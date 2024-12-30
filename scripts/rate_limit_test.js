// We want to ensure the rate limiter works as expected.
//
// To run the script:
//  1. Install k6 tool from https://grafana.com/docs/k6/latest/set-up/install-k6/
//  2. Run the app: air
//  3. Run the script: k6 run ./scripts/rate_limit_test.js

import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { duration: "5s", target: 5 },
    { duration: "5s", target: 10 },
    { duration: "15s", target: 100 },
  ],
};

export default function () {
  const res = http.post(
    "http://localhost:3000/api/v1/search",
    JSON.stringify({ query: "something" }),
    { headers: { "Content-Type": "application/json" } }
  );

  check(res, {
    "Status 200": (r) => r.status === 200,
    "Rate limited": (r) => r.status === 429,
  });

  sleep(1);
}
