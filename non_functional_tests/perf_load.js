import http from 'k6/http';
import { check, sleep } from 'k6';

// Load context from setup
const context = JSON.parse(open('./perf_context.json'));

export const options = {
  scenarios: {
    board_load: {
      executor: 'constant-vus',
      vus: 50,
      duration: '15s',
    },
  },
  thresholds: {
    // Assert P95 time of response (http_req_duration) is under 300ms
    http_req_duration: ['p(95)<300'],
  },
};

export default function () {
  const url = `http://localhost:8000/api/v2/boards/${context.boardId}/blocks?all=true`;
  const params = {
    headers: {
      'Authorization': `Bearer ${context.userToken}`,
      'X-Requested-With': 'XMLHttpRequest',
    },
  };
  
  const res = http.get(url, params);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
  
  sleep(0.1);
}
