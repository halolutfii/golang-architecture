import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

// Base URL bisa di-override: k6 run -e BASE_URL=http://localhost:3000 rate_limiter_batch_test.js
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

// Metrik untuk memisahkan hasil 200 vs 429.
const allowed = new Counter('req_allowed_200');    // lolos rate limiter
const throttled = new Counter('req_throttled_429'); // ditolak rate limiter

export const options = {
  scenarios: {
    // 100 VU, tiap VU mengirim 2 request SEKALIGUS (batch) per iterasi.
    // Dengan rate limiter 1 req/detik/token: 1 lolos (200), 1 kena 429.
    // Iterasi dibatasi 1 per detik per VU agar window rate limiter reset.
    batch_two_at_once: {
      executor: 'constant-arrival-rate',
      rate: 100,             // 100 iterasi/detik
      timeUnit: '1s',
      duration: '60s',
      preAllocatedVUs: 100,
      maxVUs: 100,
    },
  },
  thresholds: {
    // Ekspektasi: ~50% lolos (200), ~50% ditolak (429).
    // Tidak ada threshold error==0 karena 429 memang DIHARAPKAN di sini.
    'req_allowed_200': ['count>0'],
    'req_throttled_429': ['count>0'],
  },
};

// Cache token unik per VU.
let token = null;

function login() {
  const idx = ((__VU - 1) % 100) + 1;
  const userId = `user-${String(idx).padStart(3, '0')}`;

  const res = http.post(
    `${BASE_URL}/api/users/_login`,
    JSON.stringify({ id: userId, password: userId }),
    { headers: { 'Content-Type': 'application/json' }, tags: { name: 'login' } }
  );

  check(res, { 'login 200': (r) => r.status === 200 });

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

  const authHeader = { Authorization: token };

  // Kirim 2 request /api/users/_current SEKALIGUS dengan token yang sama.
  const responses = http.batch([
    { method: 'GET', url: `${BASE_URL}/api/users/_current`, params: { headers: authHeader, tags: { name: 'current' } } },
    { method: 'GET', url: `${BASE_URL}/api/users/_current`, params: { headers: authHeader, tags: { name: 'current' } } },
  ]);

  let ok = 0;
  let limited = 0;
  for (const res of responses) {
    if (res.status === 200) {
      allowed.add(1);
      ok++;
    } else if (res.status === 429) {
      throttled.add(1);
      limited++;
    }
  }

  // Per iterasi (2 request), harapannya tepat 1 lolos dan 1 ditolak.
  check(null, {
    'tepat 1 lolos (200)': () => ok === 1,
    'tepat 1 ditolak (429)': () => limited === 1,
  });
}
