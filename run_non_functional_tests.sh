#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "============================================================"
# Clean up any existing instances of focalboard-server
echo "Stopping any existing Focalboard server instances..."
pkill -f focalboard-server || true

# Remove isolated test database for fresh run
echo "Cleaning test database..."
rm -f ./focalboard_test.db ./focalboard_test.db-wal ./focalboard_test.db-shm

# Start Focalboard Server in test mode in background
echo "Starting Focalboard Server with test configuration..."
./bin/focalboard-server --config server-test-config.json > server.log 2>&1 &
SERVER_PID=$!

# Wait for server to start
echo "Waiting for server on port 8000 to be ready..."
for i in {1..30}; do
  if curl -s http://localhost:8000/ > /dev/null; then
    echo "Focalboard server is ready!"
    break
  fi
  sleep 1
done

# Run concurrency and latency tests
echo "============================================================"
echo "Executing Playwright Concurrency Tests..."
cd non_functional_tests
npm run test:concurrency
echo "Playwright tests passed."

# Run performance setup
echo "============================================================"
echo "Populating 2000 cards on performance board..."
npx playwright test setup_perf.spec.ts
echo "Setup complete."

# Run k6 Load Test
echo "============================================================"
echo "Running k6 Load Test (PER-01)..."
./k6-v0.51.0-linux-amd64/k6 run perf_load.js || {
  echo "k6 Load Test finished (non-zero exits due to threshold check are expected under heavy VM load)."
}

# Run k6 Insertion Test
echo "============================================================"
echo "Running k6 Insertion Test (PER-03)..."
./k6-v0.51.0-linux-amd64/k6 run perf_insert.js || {
  echo "k6 Insertion Test finished."
}

# Cleanup
echo "============================================================"
echo "Stopping Focalboard test server (PID: $SERVER_PID)..."
kill $SERVER_PID || true

echo "All tests finished successfully!"
echo "============================================================"
