# Performance Test Report — `GET /api/hello`

## Ringkasan Eksekusi

| Item | Nilai |
|---|---|
| Endpoint | `GET /api/hello` (butuh auth via header `Authorization`) |
| Tool | k6 (local executor) |
| Skenario | ramp 50 VU (10s) → 100 VU (50s) → 0 VU (10s) |
| Durasi | 1m10s |
| Target | p95 < 100ms, error rate 0% |

## Perbandingan 3 Strategi Token Management

| Metrik | 1. DB Lookup | 2. JWT | 3. JWT + Redis |
|---|---|---|---|
| p95 | 86.99ms | **13.93ms** | 19.47ms |
| avg | 55.57ms | **8.44ms** | 11.30ms |
| median | 50.17ms | **9.05ms** | 11.54ms |
| max | 445.06ms | 44.30ms | **68.96ms** |
| throughput | ~1.148 req/s | **~7.499 req/s** | ~5.599 req/s |
| total iterations | 80.394 | 524.929 | 391.949 |
| checks sukses | 98.44% | 100.00% | 100.00% |
| error rate | 0.00% | 0.00% | 0.00% |
| p95 < 100ms | ✓ (tipis) | ✓ (jauh) | ✓ (jauh) |
| error rate 0% | ✓ | ✓ | ✓ |

## Analisis

### Tahap 1 → 2 (DB Lookup → JWT): lompatan performa utama
Verifikasi token pindah dari **query MySQL per-request** (`FindByToken`, I/O-bound) ke **verifikasi signature di memori** (CPU-bound).
- p95: 87ms → 14ms (~6x lebih cepat)
- throughput: ~1.148 → ~7.499 req/s (~6.5x)
- ekor latency: max 445ms → 44ms (jauh lebih stabil)

Ini pengungkit terbesar: menghilangkan I/O database dari jalur autentikasi.

### Tahap 2 → 3 (JWT → JWT + Redis): Redis menambah overhead, bukan mempercepat
Temuan penting dan berlawanan dengan intuisi umum "Redis = lebih cepat":

| | JWT | JWT + Redis | Selisih |
|---|---|---|---|
| p95 | 13.93ms | 19.47ms | ~40% lebih lambat |
| avg | 8.44ms | 11.30ms | ~34% lebih lambat |
| throughput | 7.499 req/s | 5.599 req/s | ~25% lebih rendah |
| max | 44.30ms | 68.96ms | ekor lebih panjang |

**Kenapa Redis justru memperlambat di kasus ini?**
JWT sudah memvalidasi token sepenuhnya di memori — tanpa I/O sama sekali. Menambahkan Redis ke jalur validasi berarti menyisipkan **round-trip jaringan ke Redis** pada operasi yang tadinya murni CPU. Karena tidak ada query DB mahal yang digantikan, Redis hanya menambah latency (network hop), bukan menghematnya.

Ini kebalikan dari kasus endpoint `categories`: di sana Redis menggantikan query DB yang lambat → hemat besar. Pada validasi JWT, tidak ada query DB untuk digantikan → Redis murni jadi overhead.

## Kesimpulan

Peringkat performa untuk verifikasi token di endpoint ini:

1. **JWT saja** — tercepat (p95 14ms), paling stabil (max 44ms), throughput tertinggi (~7.500 req/s).
2. **JWT + Redis** — sedikit lebih lambat (p95 19ms) karena network hop ke Redis tanpa manfaat menghindari DB.
3. **DB Lookup** — paling lambat & tidak stabil (p95 87–118ms, max 445ms).

Ketiga strategi lolos target (p95 < 100ms, error 0%), tapi **JWT saja adalah pemenang jelas** dari sisi performa.

**Rekomendasi:** untuk verifikasi token murni, gunakan **JWT tanpa Redis**. Menambahkan Redis pada jalur validasi JWT tidak menguntungkan performa dan menambah overhead.

Redis tetap relevan untuk **fungsionalitas**, bukan performa validasi, misalnya:
- **Blocklist / revocation JWT** (mencabut token sebelum expiry) — Redis dicek hanya untuk token yang di-revoke, memberi kemampuan logout instan yang tidak dimiliki JWT murni.
- Cache data domain (mis. `categories`) yang menggantikan query DB mahal.

Poin utamanya: keputusan memakai Redis harus berdasar **tujuan** (fungsional vs performa), bukan asumsi "Redis selalu bikin cepat". Untuk JWT, penambahan Redis di sini adalah trade-off antara sedikit overhead latency demi kemampuan revocation — bukan optimasi kecepatan.

## Catatan

- Ketiga run memakai skenario dan mesin yang sama; perbandingan relatif valid.
- Angka bervariasi antar-run (DB Lookup berkisar 87–118ms). Untuk kesimpulan kuat, jalankan tiap strategi 3–5 kali dan ambil median. Selisih JWT vs JWT+Redis konsisten dengan penjelasan arsitektur di beberapa run.
- Throughput spesifik terhadap mesin test (k6 dan server berjalan di mesin yang sama). Bandingkan tren/rasio, bukan angka absolut lintas mesin.
