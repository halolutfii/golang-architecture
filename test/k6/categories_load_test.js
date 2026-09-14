import http from 'k6/http';
import { check } from 'k6';

// how to run
// k6 run test/k6/categories_load_test.js or k6 run -e BASE_URL=http://localhost:3000 test/k6/categories_load_test.js or & "C:\Program Files\k6\k6.exe" run test/k6/categories_load_test.js


// Base URL bisa di-override: k6 run -e BASE_URL=http://localhost:3000 categories_load_test.js
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export const options = {
  stages: [
    { duration: '10s', target: 50 },  // ramp-up ke 50 VU dalam 10 detik
    { duration: '50s', target: 100 }, // naik & tahan di 100 VU selama 50 detik
    { duration: '10s', target: 0 },   // ramp-down ke 0 VU dalam 10 detik
  ],
  thresholds: {
    // p(95) waktu response harus <= 200ms
    http_req_duration: ['p(95)<=200'],
  },
};

export default function () {
  const res = http.get(`${BASE_URL}/api/categories`);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'p95 duration <= 200ms': (r) => r.timings.duration <= 200,
  });
}
