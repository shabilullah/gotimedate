# Stress Testing Guide

This guide explains how to stress test the Go TimeDate API using [k6](https://k6.io/).

## Prerequisites

- **k6**: [Installation Instructions](https://k6.io/docs/getting-started/installation/)
- **hey** (Optional for quick tests): `go install github.com/rakyll/hey@latest`

## Running Tests

### 1. REST API Load Test

This script tests the `/api/v1/time` and `/api/v1/timezones` endpoints.

```bash
k6 run scripts/stress/load_test.js
```

To override the base URL:
```bash
k6 run -e BASE_URL=http://your-server:8080 scripts/stress/load_test.js
```

### 2. WebSocket Load Test

This script tests the `/ws/time` streaming endpoint, simulating multiple active subscribers.

```bash
k6 run scripts/stress/ws_test.js
```

To override the WebSocket URL:
```bash
k6 run -e WS_URL=ws://your-server:8080/ws/time scripts/stress/ws_test.js
```

### 3. Quick Benchmarking with `hey`

For a simple burst of requests to a specific endpoint:

```bash
hey -n 1000 -c 50 http://localhost:8080/api/v1/time
```

## Metrics to Watch

- **http_req_duration**: Aim for `p(95) < 200ms`.
- **http_req_failed**: Should be 0%.
- **ws_connecting**: Latency for establishing WebSocket connections.
- **Service CPU/Memory**: Monitor your system resources during the test to see where it caps out.

## Interpreting Results & Finding Capacity

To find the maximum load your system can handle, you should look at these specific metrics in the k6 summary output:

### 1. Throughput (Requests Per Second)
Look for the `http_reqs` metric. It shows the total count and the rate per second:
```text
http_reqs..................: 12000  133.333333/s
```
In this example, the system handled **133 requests per second (RPS)**.

### 2. Identifying the Limit
To find the absolute maximum capacity:
1.  **Increase VUs**: Edit `load_test.js` or `ws_test.js` and increase the `target` in `stages` (e.g., from 20 to 100, 500, or 1000).
2.  **Watch for Failures**: Look at `http_req_failed`. As soon as this is > 0%, you have exceeded the system's capacity.
3.  **Watch for Latency Spikes**: Look at `http_req_duration`. If the `p(95)` (95th percentile) jumps from a few milliseconds to several seconds, the system is bottlenecked (likely CPU or database/IO).

### 3. WebSocket Capacity
For WebSockets, the key metric is **concurrent connections**. The `vus` (Virtual Users) metric tells you how many active connections were maintained simultaneously. If connections start dropping or timing out, you've hit the limit of the Go server's memory or file descriptor limits.

## Troubleshooting Bottlenecks

If you see timeouts or connection errors while CPU usage is still low (e.g., < 50%), you are likely hitting one of these common bottlenecks:

### 1. Logging Overhead
The default Logger middleware writes every request to the console/file. This is a synchronous I/O operation that is much slower than the Go logic itself.
- **Fix**: Set `LOG_LEVEL=error` in your `config.env` to disable per-request logging.

### 2. OS File Descriptor Limits
Every connection is a "file". If the limit is 1,024, you can't have 10,000 users.
- **Linux Fix**: Run `ulimit -n 50000` before starting the server.
- **Windows Fix**: This is handled via the registry, but generally Windows allows more by default for handles, though socket limits still apply.

### 3. Ephemeral Port Exhaustion
If `k6` and the server are on the same machine, they might run out of ports to talk to each other.
- **Fix**: Run `k6` from a different machine or increase the OS ephemeral port range.

### 4. Middleware Timeouts
The system has a default timeout of 60 seconds (`middleware.Timeout(60*time.Second)`). If the server is stuck behind a logging bottleneck, requests might sit in the queue until they hit this 60s limit and get cancelled.
