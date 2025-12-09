import http from 'k6/http';
import { check, sleep } from 'k6';

// 1. Configuration (Options)
export const options = {
  // Use a different executor to control requests per second (RPS) directly
  scenarios: {
    constant_stress: {
      executor: 'constant-arrival-rate', // Inject requests at a constant rate
      rate: 1000,                        // Target: 100 iterations (requests) per second
      timeUnit: '1s',                   // The rate is per 1 second
      duration: '1m',                   // Run the test for 1 minute
      preAllocatedVUs: 200,             // Start with 200 Virtual Users pre-allocated
      maxVUs: 1000,                     // Allow up to 1000 VUs if needed to maintain the rate
    },
  },

  // Set much tighter thresholds to fail immediately if performance degrades
  thresholds: {
    // 95% of requests MUST complete within 100ms
    'http_req_duration': ['p(95) < 100'],
    // 99.9% of requests must succeed
    'checks': ['rate>0.999'],
    // Total throughput must not drop below a minimum expected rate
    'iterations': ['rate>90'],
  },
};

// 2. Main Logic (Default Function)
export default function () {
  // Target a simple, unauthenticated endpoint through your API Gateway
  // The sleep(1) is NOW largely ignored because the executor controls the arrival rate.
  const res = http.get('http://localhost:8080/api/v1/test');

  // Perform checks (assertions) on the response
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
}