import http from 'k6/http';
import { check } from 'k6';

const context = JSON.parse(open('./perf_context.json'));

function generateGuid() {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < 26; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 500,
      timeUnit: '1s',
      duration: '10s',
      preAllocatedVUs: 50,
      maxVUs: 100,
    },
  },
  thresholds: {
    // We expect requests to be successful
    http_req_failed: ['rate<0.01'],
  },
};

export default function () {
  const url = `http://localhost:8000/api/v2/boards/${context.boardId}/blocks`;
  const cardId = generateGuid();
  const payload = JSON.stringify([{
    id: cardId,
    parentId: context.boardId,
    boardId: context.boardId,
    type: 'card',
    title: `Perf Insertion Card ${cardId}`,
    fields: {
      properties: {
        [context.statusPropId]: 'opt-todo'
      }
    },
    createAt: Date.now(),
    updateAt: Date.now(),
    deleteAt: 0
  }]);

  const params = {
    headers: {
      'Authorization': `Bearer ${context.userToken}`,
      'X-Requested-With': 'XMLHttpRequest',
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}
