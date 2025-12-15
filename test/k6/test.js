import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    overload: {
      executor: "constant-arrival-rate",
      rate: 50,        // 50 req/sec
      timeUnit: "1s",
      duration: "10s",
      preAllocatedVUs: 20,
      maxVUs: 50,
    },
  },
};

export default function () {
  const res = http.get("http://localhost:8080/healthz");

  check(res, {
    "status is 200": (r) => r.status === 200,
  });

  sleep(1);
}
