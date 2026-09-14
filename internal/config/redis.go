package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper) *redis.Client {
	host := viper.GetString("redis.host")
	database := viper.GetInt("redis.database")

	var client = redis.NewClient(&redis.Options{
		Addr: host,
		DB:   database,
	})

	return client
}

// How to inspect Redis locally (container name: "redis").
// Open an interactive redis-cli session:
//   docker exec -it redis redis-cli

// PING                 # cek hidup, balas PONG
// KEYS *               # lihat semua key
// GET categories       # lihat isi cache categories (JSON)
// TTL categories       # sisa waktu hidup key (detik)
// DBSIZE               # jumlah key
// DEL categories       # hapus key categories (buat test cache miss)
// FLUSHDB              # hapus semua key di db aktif
// exit                 # keluar
