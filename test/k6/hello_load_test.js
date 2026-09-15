import http from 'k6/http';
import { check } from 'k6';

// Base URL & token bisa di-override lewat env var:
//   k6 run -e BASE_URL=http://localhost:3000 -e TOKEN=<token-mu> hello_load_test.js
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';
const TOKEN = __ENV.TOKEN || '6d4b50ee-ad0b-413e-91ae-031220627ae7';

export const options = {
  stages: [
    { duration: '10s', target: 50 },  // ramp-up ke 50 VU dalam 10 detik
    { duration: '50s', target: 100 }, // naik & tahan di 100 VU selama 50 detik
    { duration: '10s', target: 0 },   // ramp-down ke 0 VU dalam 10 detik
  ],
  thresholds: {
    // response time p95 harus di bawah 100ms
    http_req_duration: ['p(95)<100'],
    // error rate harus 0% (tidak boleh ada request yang gagal)
    http_req_failed: ['rate==0'],
  },
};

const params = {
  headers: {
    Authorization: TOKEN,
  },
};

export default function () {
  const res = http.get(`${BASE_URL}/api/hello`, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'duration < 100ms': (r) => r.timings.duration < 100,
  });
}
