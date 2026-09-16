# Performance Test Report — `GET /api/hello`

## Ringkasan Eksekusi

| Item | Nilai |
|---|---|
| Endpoint | `GET /api/hello` (butuh auth via header `Authorization`) |
| Tool | k6 (local executor) |
| Skenario | ramp 50 VU (10s) → 100 VU (50s) → 0 VU (10s) |
| Durasi | 1m10s |

## Perbandingan: Token DB Lookup vs JWT

| Metrik | DB Lookup (Run 1) | DB Lookup (Run 2) | **JWT** |
|---|---|---|---|
| p95 | 118.20ms ❌ | 86.99ms ✓ | **13.93ms ✓** |
| avg | 67.86ms | 55.57ms | **8.44ms** |
| median | 58.70ms | 50.17ms | **9.05ms** |
| min | 39.66ms | 34.07ms | **0s** |
| max | 233.86ms | 445.06ms | **44.3ms** |
| throughput | ~942 req/s | ~1.148 req/s | **~7.499 req/s** |
| total iterations | 65.953 | 80.394 | **524.929** |
| checks sukses | 94.52% | 98.44% | **100.00%** |
| error rate | 0.00% | 0.00% | **0.00%** |
| p95 < 100ms | ❌ | ✓ (tipis) | **✓ (jauh)** |
| error rate 0% | ✓ | ✓ | **✓** |

## Perbaikan JWT relatif ke DB Lookup

| | vs Run 1 (118ms) | vs Run 2 (87ms) |
|---|---|---|
| p95 | ~8.5x lebih cepat | ~6.2x lebih cepat |
| throughput | ~8x lebih banyak | ~6.5x lebih banyak |

## Analisis

- **JWT lolos kedua target dengan margin besar.** p95 = 13.93ms (jauh di bawah 100ms), error rate 0%, dan **100% checks sukses** dari 1.049.858 checks — tidak ada satu pun request yang melebihi 100ms. Ini kualitas hasil yang jauh lebih baik daripada versi DB lookup yang selalu punya ekor lambat.

- **Kenapa JWT jauh lebih cepat?** Perbedaan arsitekturnya fundamental:
  - **DB lookup (sebelumnya):** setiap request memanggil `FindByToken` → query ke MySQL untuk memvalidasi token. Ini I/O-bound; latency-nya ditentukan round-trip DB dan sensitif terhadap beban DB (terlihat dari variasi Run 1 vs Run 2: 118ms vs 87ms).
  - **JWT (sekarang):** token divalidasi dengan memverifikasi **signature secara kriptografis di memori**, tanpa menyentuh database. Ini CPU-bound dan sangat cepat, sehingga latency konsisten (median 9ms, p95 14ms, max hanya 44ms).

- **Konsistensi jauh lebih baik.** Bandingkan sebaran:
  - DB lookup: median 50ms → p95 87ms → max 445ms (ekor panjang, tidak stabil).
  - JWT: median 9ms → p95 14ms → max 44ms (sebaran rapat, stabil).
  Selisih min-max JWT sangat kecil, tanda tidak ada bottleneck I/O yang bikin outlier.

- **Throughput naik ~6–8x** (dari ~1.000 ke ~7.500 req/s) dengan VU yang sama, karena tiap request tidak lagi menahan koneksi DB.

## Kesimpulan

Mengganti verifikasi token dari **DB lookup** ke **JWT** menyelesaikan bottleneck utama endpoint `/api/hello`:

- p95: **~87–118ms → 14ms** (lolos target 100ms dengan margin besar dan stabil)
- throughput: **~1.000 → ~7.500 req/s**
- konsistensi: ekor latency hilang (max 445ms → 44ms)

Ini pola yang sama seperti perbaikan endpoint categories (dari DB langsung ke cache): **memindahkan verifikasi dari database ke proses in-memory** adalah pengungkit performa terbesar. Dengan JWT, tidak diperlukan lagi query DB per-request maupun cache token di Redis untuk kebutuhan autentikasi ini.

## Catatan

- Run JWT dijalankan tanpa `-e TOKEN` (script memakai token default). Pastikan token JWT yang dipakai valid dan belum kedaluwarsa saat test agar error rate tetap 0%.
- Trade-off JWT yang perlu diingat (di luar performa): token JWT tidak bisa langsung "dicabut" di sisi server sebelum expiry, berbeda dengan token DB yang bisa dihapus. Pertimbangkan strategi expiry/refresh atau blocklist bila diperlukan.
- Throughput bersifat spesifik terhadap mesin test (k6 dan server berjalan di mesin yang sama). Bandingkan tren/rasio, bukan angka absolut lintas mesin.
