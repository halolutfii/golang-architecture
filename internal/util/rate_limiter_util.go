package util

import (
	"context"
	"fmt"
	"golang-clean-architecture/internal/model"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RateLimiterUtil struct {
	Redis      *redis.Client
	MaxRequest int64
	Duration   time.Duration
}

func NewRateLimiterUtil(redis *redis.Client) *RateLimiterUtil {
	return &RateLimiterUtil{
		Redis:      redis,
		MaxRequest: 1,
		Duration:   time.Second * 1,
	}
}

// IsAllowed menerapkan rate limiting dengan algoritma SLIDING WINDOW.
//
// Berbeda dengan fixed-window (yang me-reset counter di titik waktu tetap dan
// rawan celah/tumpukan di batas window), sliding window selalu menghitung jumlah
// request dalam rentang "Duration terakhir dari waktu sekarang". Jendelanya
// bergerak mengikuti waktu request, sehingga hasilnya presisi berapa pun timing-nya.
//
// Implementasi memakai Redis sorted set:
//   - anggota (member) = id unik per request, score = timestamp (nanodetik)
//   - ZREMRANGEBYSCORE membuang entri yang sudah keluar dari jendela
//   - ZCARD menghitung entri yang masih di dalam jendela
//
// Semua operasi dijalankan dalam satu TxPipeline agar atomic terhadap request
// yang datang bersamaan (mis. http.batch), sehingga tidak ada race condition.
func (u RateLimiterUtil) IsAllowed(ctx context.Context, auth *model.Auth) bool {
	key := "rate_limit:" + auth.ID

	now := time.Now().UnixNano()
	windowStart := now - int64(u.Duration)

	pipe := u.Redis.TxPipeline()

	// 1. Buang request yang sudah di luar jendela (lebih tua dari windowStart).
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// 2. Catat request saat ini: score = now, member = id unik.
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: uuid.NewString()})

	// 3. Hitung jumlah request dalam jendela (termasuk yang barusan ditambahkan).
	countCmd := pipe.ZCard(ctx, key)

	// 4. Set TTL agar key otomatis bersih jika token idle.
	pipe.Expire(ctx, key, u.Duration)

	if _, err := pipe.Exec(ctx); err != nil {
		fmt.Println("Error executing rate limiter pipeline: ", err)
		return false
	}

	// Diizinkan bila jumlah request dalam jendela tidak melebihi batas.
	return countCmd.Val() <= u.MaxRequest
}

// --- Implementasi lama: FIXED WINDOW (disimpan untuk perbandingan) ---
//
// Kelemahan: window di-reset dari request pertama (via EXPIRE saat increment==1),
// sehingga tidak selaras dengan ritme request. Pada trafik tepat di batas, ini
// menyebabkan sebagian window "belum reset" -> request ditolak walau seharusnya
// lolos (hasil test batch: ~32% lolos, bukan 50/50).
//
// func (u RateLimiterUtil) IsAllowed(ctx context.Context, auth *model.Auth) bool {
// 	key := auth.ID
//
// 	increment, err := u.Redis.Incr(ctx, key).Result()
// 	if err != nil {
// 		fmt.Println("Error incrementing: ", err)
// 		return false
// 	}
//
// 	if increment == 1 {
// 		err := u.Redis.Expire(ctx, key, u.Duration).Err()
// 		if err != nil {
// 			fmt.Println("Error setting expiration: ", err)
// 			return false
// 		}
// 	}
//
// 	return increment <= u.MaxRequest
// }
