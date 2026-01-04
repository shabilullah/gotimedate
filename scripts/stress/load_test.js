import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 50000 }, // ramp up to 50000 users
        { duration: '1m', target: 50000 },  // stay at 50000 users
        { duration: '30s', target: 0 },  // ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
        http_req_failed: ['rate<0.01'],    // less than 1% errors
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
    const res = http.get(`${BASE_URL}/api/v1/time`);
    check(res, {
        'status is 200': (r) => r.status === 200,
        'has timestamp': (r) => r.json().timestamp !== undefined,
    });

    const tzRes = http.get(`${BASE_URL}/api/v1/timezones`);
    check(tzRes, {
        'timezones status is 200': (r) => r.status === 200,
    });

    sleep(1);
}
