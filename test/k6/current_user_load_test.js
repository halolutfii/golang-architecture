import http from 'k6/http';
import { check } from 'k6';

// Base URL bisa di-override: k6 run -e BASE_URL=http://localhost:3000 current_user_load_test.js
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export const options = {
  scenarios: {
    // 100 VU, masing-masing melakukan tepat 1 request per detik.
    // constant-arrival-rate menjaga laju tetap 100 req/detik total,
    // dengan 100 VU yang dialokasikan -> efektif 1 req/detik per VU.
    one_rps_per_vu: {
      executor: 'constant-arrival-rate',
      rate: 100,               // 100 iterasi per timeUnit
      timeUnit: '1s',          // -> 100 request per detik total
      duration: '60s',         // total durasi 60 detik
      preAllocatedVUs: 100,    // 100 VU disiapkan
      maxVUs: 100,             // dikunci di 100 VU (tiap VU ~1 req/detik)
    },
  },
  thresholds: {
    // Semua request harus sukses: tidak ada error, tidak ada 429.
    http_req_failed: ['rate==0'],
    checks: ['rate==1.0'],
  },
};

// Tiap VU login sekali (di init per-iterasi pertama) dan menyimpan token unik.
// Password == id, sesuai seed_users_100.sql (user-001..user-100).
let token = null;

function login() {
  const idx = ((__VU - 1) % 100) + 1;
  const userId = `user-${String(idx).padStart(3, '0')}`;

  const res = http.post(
    `${BASE_URL}/api/users/_login`,
    JSON.stringify({ id: userId, password: userId }),
    { headers: { 'Content-Type': 'application/json' }, tags: { name: 'login' } }
  );

  const ok = check(res, { 'login 200': (r) => r.status === 200 });
  if (!ok) return null;

  try {
    return res.json('data.token');
  } catch (e) {
    return null;
  }
}

export default function () {
  if (!token) {
    token = login();
  }

  const res = http.get(`${BASE_URL}/api/users/_current`, {
    headers: { Authorization: token },
    tags: { name: 'current' },
  });

  check(res, {
    'current status 200': (r) => r.status === 200,
  });
}
