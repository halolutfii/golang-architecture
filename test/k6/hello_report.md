# Performance Test Report — `GET /api/hello`

## Ringkasan Eksekusi

| Item | Nilai |
|---|---|
| Endpoint | `GET /api/hello` (butuh auth via header `Authorization`) |
| Tool | k6 (local executor) |
| Skenario | ramp 50 VU (10s) → 100 VU (50s) → 0 VU (10s) |
| Durasi | 1m10s |

## Target vs Hasil (per run)

| Target | Threshold | Run 1 | Run 2 | Status akhir |
|---|---|---|---|---|
| Response time p95 < 100ms | `p(95)<100` | 118.2ms ❌ | **86.99ms ✓** | ✓ LOLOS |
| Error rate 0% | `rate==0` | 0.00% ✓ | 0.00% ✓ | ✓ LOLOS |

## Perbandingan Metrik `http_req_duration`

| Metrik | Run 1 | Run 2 |
|---|---|---|
| min | 39.66ms | 34.07ms |
| avg | 67.86ms | **55.57ms** |
| median | 58.70ms | **50.17ms** |
| p90 | 102.09ms | **71.11ms** |
| p95 | 118.20ms | **86.99ms** |
| max | 233.86ms | 445.06ms |
| throughput | ~942 req/s | **~1.148 req/s** |
| total iterations | 65.953 | 80.394 |
| error rate | 0.00% | 0.00% |

## Detail Checks (Run 2)

| Check | Hasil |
|---|---|
| status is 200 | ✓ 100% (semua request sukses) |
| duration < 100ms | ✗ 96% — 77.887 lolos / 2.507 gagal |
| `http_req_failed` | 0.00% (0 dari 80.394) |

## Analisis

- **Kedua target tercapai di Run 2.** p95 turun dari 118.2ms → 86.99ms (di bawah target 100ms), dan error rate tetap 0%. Semua percentile membaik: median 58→50ms, p90 102→71ms.

- **Kenapa Run 2 lebih baik dari Run 1?** Beban dua run identik, jadi perbedaan ini berasal dari kondisi mesin/lingkungan, bukan perubahan kode. Faktor umum: MySQL sudah "warm" (buffer pool terisi, koneksi ter-reuse), cache OS lebih siap, dan beban latar mesin lebih ringan saat Run 2. Throughput naik ~22% (942 → 1.148 req/s) mendukung ini.

- **p95 masih dekat batas (87ms vs target 100ms).** Marginnya tipis (~13ms). Karena `/api/hello` melakukan query DB verifikasi token (`FindByToken`) di setiap request, hasilnya sensitif terhadap kondisi MySQL. Run berikutnya bisa saja kembali menembus 100ms jika mesin lebih sibuk — lihat catatan konsistensi di bawah.

- **Perhatikan `max = 445ms` di Run 2** (naik dari 234ms di Run 1). Meski p95 lebih baik, ada ekor (outlier) yang lebih lambat — kemungkinan lonjakan sesaat saat koneksi DB rebutan. Karena hanya menyentuh <4% request (yang gagal cek <100ms), tidak memengaruhi p95, tapi menunjukkan latency belum sepenuhnya stabil di ekor.

## Rekomendasi

Karena p95 masih menempel di batas dan bergantung kondisi DB, untuk margin yang aman dan stabil:

1. **Cache token di Redis** (token → user, dengan TTL). Menghilangkan query DB di mayoritas request. Pola ini terbukti menurunkan p95 endpoint categories dari ratusan ms ke belasan ms, dan akan memberi margin jauh di bawah 100ms secara konsisten.
2. **Pastikan kolom `token` di tabel `users` ter-index**, agar `FindByToken` tidak melakukan full-table scan.

## Catatan Konsistensi

- Dua run dengan beban identik memberi p95 berbeda (118ms vs 87ms). Ini menegaskan bahwa untuk endpoint yang menyentuh DB, **satu run tidak cukup** untuk menyimpulkan lolos/tidak. Disarankan menjalankan 3–5 kali dan mengambil median, atau menetapkan margin target yang cukup (mis. optimasi hingga p95 << 100ms) sebelum menyatakan lolos.
- Token diberikan via `-e TOKEN=<token>` karena berganti setiap login.
- Throughput bersifat spesifik terhadap mesin test (k6 dan server berjalan di mesin yang sama). Bandingkan tren/rasio, bukan angka absolut lintas mesin.
