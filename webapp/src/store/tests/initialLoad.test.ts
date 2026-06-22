jest.mock('../../octoClient', () => ({
    __esModule: true,
    default: {
        getMe: jest.fn(),
        getMyConfig: jest.fn(),
        getTeam: jest.fn(),
        getTeams: jest.fn(),
        getBoards: jest.fn(),
        getMyBoardMemberships: jest.fn(),
        getTeamTemplates: jest.fn(),
        getBoardsCloudLimits: jest.fn(),
        getBoard: jest.fn(),
        getAllBlocks: jest.fn(),
    },
}))

import client from '../../octoClient'

import {
    initialLoad,
    initialReadOnlyLoad,
    loadBoardData,
    loadBoards,
    loadMyBoardsMemberships,
    getUserBlockSubscriptions,
    getUserBlockSubscriptionList,
} from '../initialLoad'

import {ErrorId} from '../../errors'

const mockClient = client as jest.Mocked<typeof client>

beforeEach(() => {
    jest.clearAllMocks()
})

describe('initialLoad thunk', () => {
    test('returns initial data when user and team exist', async () => {
        mockClient.getMe.mockResolvedValue({id: 'user1'} as any)
        mockClient.getMyConfig.mockResolvedValue([{key: 'language', value: 'en'}] as any)
        mockClient.getTeam.mockResolvedValue({id: 'team1'} as any)
        mockClient.getTeams.mockResolvedValue([{id: 'team1'}] as any)
        mockClient.getBoards.mockResolvedValue([{id: 'board1'}] as any)
        mockClient.getMyBoardMemberships.mockResolvedValue([{boardId: 'board1'}] as any)
        mockClient.getTeamTemplates.mockResolvedValue([{id: 'template1'}] as any)
        mockClient.getBoardsCloudLimits.mockResolvedValue({cards: 10} as any)

        const result = await initialLoad()(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('initialLoad/fulfilled')
        expect(result.payload).toEqual({
            team: {id: 'team1'},
            teams: [{id: 'team1'}],
            boards: [{id: 'board1'}],
            boardsMemberships: [{boardId: 'board1'}],
            boardTemplates: [{id: 'template1'}],
            limits: {cards: 10},
            myConfig: [{key: 'language', value: 'en'}],
        })
    })

    test('rejects when user is not logged in', async () => {
        mockClient.getMe.mockResolvedValue(null as any)
        mockClient.getMyConfig.mockResolvedValue([] as any)
        mockClient.getTeam.mockResolvedValue({id: 'team1'} as any)
        mockClient.getTeams.mockResolvedValue([] as any)
        mockClient.getBoards.mockResolvedValue([] as any)
        mockClient.getMyBoardMemberships.mockResolvedValue([] as any)
        mockClient.getTeamTemplates.mockResolvedValue([] as any)
        mockClient.getBoardsCloudLimits.mockResolvedValue(null as any)

        const result = await initialLoad()(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('initialLoad/rejected')
        expect((result as any).error.message).toBe(ErrorId.NotLoggedIn)
    })

    test('rejects when team is undefined', async () => {
        mockClient.getMe.mockResolvedValue({id: 'user1'} as any)
        mockClient.getMyConfig.mockResolvedValue([] as any)
        mockClient.getTeam.mockResolvedValue(null as any)
        mockClient.getTeams.mockResolvedValue([] as any)
        mockClient.getBoards.mockResolvedValue([] as any)
        mockClient.getMyBoardMemberships.mockResolvedValue([] as any)
        mockClient.getTeamTemplates.mockResolvedValue([] as any)
        mockClient.getBoardsCloudLimits.mockResolvedValue(null as any)

        const result = await initialLoad()(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('initialLoad/rejected')
        expect((result as any).error.message).toBe(ErrorId.TeamUndefined)
    })
})

describe('initialReadOnlyLoad thunk', () => {
    test('returns board and blocks when board exists', async () => {
        mockClient.getBoard.mockResolvedValue({id: 'board1'} as any)
        mockClient.getAllBlocks.mockResolvedValue([{id: 'block1'}] as any)

        const result = await initialReadOnlyLoad('board1')(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('initialReadOnlyLoad/fulfilled')
        expect(result.payload).toEqual({
            board: {id: 'board1'},
            blocks: [{id: 'block1'}],
        })
    })

    test('rejects when read only board is invalid', async () => {
        mockClient.getBoard.mockResolvedValue(null as any)
        mockClient.getAllBlocks.mockResolvedValue([] as any)

        const result = await initialReadOnlyLoad('board1')(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('initialReadOnlyLoad/rejected')
        expect((result as any).error.message).toBe(ErrorId.InvalidReadOnlyBoard)
    })
})

describe('loadBoardData thunk', () => {
    test('returns blocks for board', async () => {
        mockClient.getAllBlocks.mockResolvedValue([{id: 'block1'}, {id: 'block2'}] as any)

        const result = await loadBoardData('board1')(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('loadBoardData/fulfilled')
        expect(result.payload).toEqual({
            blocks: [{id: 'block1'}, {id: 'block2'}],
        })
    })
})

describe('loadBoards thunk', () => {
    test('returns boards', async () => {
        mockClient.getBoards.mockResolvedValue([{id: 'board1'}, {id: 'board2'}] as any)

        const result = await loadBoards()(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('loadBoards/fulfilled')
        expect(result.payload).toEqual({
            boards: [{id: 'board1'}, {id: 'board2'}],
        })
    })
})

describe('loadMyBoardsMemberships thunk', () => {
    test('returns board memberships', async () => {
        mockClient.getMyBoardMemberships.mockResolvedValue([{boardId: 'board1'}] as any)

        const result = await loadMyBoardsMemberships()(jest.fn(), jest.fn(), undefined)

        expect(result.type).toBe('loadMyBoardsMemberships/fulfilled')
        expect(result.payload).toEqual({
            boardsMemberships: [{boardId: 'board1'}],
        })
    })
})

describe('initialLoad selectors', () => {
    test('getUserBlockSubscriptions returns subscriptions', () => {
        const state: any = {
            users: {
                blockSubscriptions: [
                    {blockId: 'block1', subscriptionId: 'sub1'},
                    {blockId: 'block2', subscriptionId: 'sub2'},
                ],
            },
        }

        expect(getUserBlockSubscriptions(state)).toEqual([
            {blockId: 'block1', subscriptionId: 'sub1'},
            {blockId: 'block2', subscriptionId: 'sub2'},
        ])
    })

    test('getUserBlockSubscriptionList returns subscriptions through selector', () => {
        const state: any = {
            users: {
                blockSubscriptions: [
                    {blockId: 'block1', subscriptionId: 'sub1'},
                ],
            },
        }

        expect(getUserBlockSubscriptionList(state)).toEqual([
            {blockId: 'block1', subscriptionId: 'sub1'},
        ])
    })
})
