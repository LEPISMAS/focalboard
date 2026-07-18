import { test, expect } from '@playwright/test';
import * as fs from 'fs';

// Helper function to generate focalboard-compatible random IDs
function generateGuid() {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < 26; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

test('Setup Performance Board with 2000 Cards', async ({ request }) => {
  console.log('--- Setting up Performance Test Environment ---');
  
  const username = 'concurrency_test_user_1';
  const email = 'concurrency_test_user_1@test.com';
  const password = 'Password123!';

  // Attempt registration (will succeed on fresh db, safe to fail/ignore on existing db)
  const registerResp = await request.post('/api/v2/register', {
    data: {
      username,
      email,
      password,
      token: '',
    },
    headers: {
      'X-Requested-With': 'XMLHttpRequest',
    }
  });
  console.log(`Registration attempt returned status: ${registerResp.status()}`);

  // Login
  console.log(`Logging in as ${username}...`);
  const loginResp = await request.post('/api/v2/login', {
    data: {
      type: 'normal',
      username,
      password,
    },
    headers: {
      'X-Requested-With': 'XMLHttpRequest',
    }
  });
  expect(loginResp.status()).toBe(200);
  const loginData = await loginResp.json();
  const userToken = loginData.token;
  console.log('Login successful.');

  // Create Board and View
  const tempBoardId = generateGuid();
  const tempViewId = generateGuid();
  const tempStatusPropId = generateGuid();

  const boardPayload = {
    id: tempBoardId,
    teamId: '0',
    channelId: '',
    title: 'Performance Testing Board',
    type: 'P',
    minimumRole: '',
    description: '',
    icon: '',
    showDescription: false,
    isTemplate: false,
    templateVersion: 0,
    properties: {},
    cardProperties: [
      {
        id: tempStatusPropId,
        name: 'Status',
        type: 'select',
        options: [
          { id: 'opt-todo', value: 'To Do', color: 'propColorDefault' },
          { id: 'opt-inprogress', value: 'In Progress', color: 'propColorDefault' },
          { id: 'opt-done', value: 'Done', color: 'propColorDefault' }
        ]
      }
    ],
    createAt: Date.now(),
    updateAt: Date.now(),
    deleteAt: 0
  };

  const blocksPayload = [
    {
      id: tempViewId,
      parentId: tempBoardId,
      boardId: tempBoardId,
      type: 'view',
      title: 'Board view',
      fields: {
        viewType: 'board',
        groupById: tempStatusPropId,
        visiblePropertyIds: [tempStatusPropId],
        visibleOptionIds: ['opt-todo', 'opt-inprogress', 'opt-done'],
        hiddenOptionIds: []
      },
      createAt: Date.now(),
      updateAt: Date.now(),
      deleteAt: 0
    }
  ];

  const createResponse = await request.post('/api/v2/boards-and-blocks', {
    data: {
      boards: [boardPayload],
      blocks: blocksPayload
    },
    headers: {
      'Authorization': `Bearer ${userToken}`,
      'X-Requested-With': 'XMLHttpRequest'
    }
  });
  expect(createResponse.status()).toBe(200);
  const createData = await createResponse.json();
  const boardId = createData.boards[0].id;
  console.log(`Performance Board created. Board ID: ${boardId}`);

  // Generate 2000 cards
  console.log('Generating 2000 cards...');
  const cards: any[] = [];
  for (let i = 0; i < 2000; i++) {
    cards.push({
      id: generateGuid(),
      parentId: boardId,
      boardId: boardId,
      type: 'card',
      title: `Performance Card ${i}`,
      fields: {
        properties: {
          [tempStatusPropId]: 'opt-todo'
        }
      },
      createAt: Date.now(),
      updateAt: Date.now(),
      deleteAt: 0
    });
  }

  // Insert cards in batches of 400
  console.log('Inserting 2000 cards via API in batches...');
  for (let i = 0; i < 5; i++) {
    const batch = cards.slice(i * 400, (i + 1) * 400);
    const insertResponse = await request.post(`/api/v2/boards/${boardId}/blocks`, {
      data: batch,
      headers: {
        'Authorization': `Bearer ${userToken}`,
        'X-Requested-With': 'XMLHttpRequest',
      }
    });
    expect(insertResponse.status()).toBe(200);
    console.log(`Inserted batch ${i + 1}/5`);
  }

  // Save perf context to file
  const perfContext = {
    boardId,
    userToken,
    statusPropId: tempStatusPropId
  };
  fs.writeFileSync('./perf_context.json', JSON.stringify(perfContext, null, 2));
  console.log('Performance test setup complete. perf_context.json written successfully.');
});
