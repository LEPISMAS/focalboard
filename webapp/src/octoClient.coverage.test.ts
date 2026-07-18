// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable @typescript-eslint/no-explicit-any */

import 'isomorphic-fetch'

import {OctoClient} from './octoClient'
import {FetchMock} from './test/fetchMock'

global.fetch = FetchMock.fn
console.log = jest.fn()
console.error = jest.fn()

const okResponse = (body: unknown = []) => new Response(JSON.stringify(body), {
    status: 200,
    headers: {'Content-Type': 'application/json'},
})

describe('OctoClient API contracts', () => {
    let client: OctoClient

    beforeEach(() => {
        FetchMock.fn.mockReset()
        FetchMock.fn.mockImplementation(async () => okResponse())
        localStorage.clear()
        client = new OctoClient('http://focalboard.test/', 'team-1', 'channel-1')
        client.token = 'session-token'
    })

    test('handles authentication, registration, configuration, and user endpoints', async () => {
        FetchMock.fn.mockImplementation(async () => okResponse({token: 'new-token'}))
        await expect(client.login('user', 'password')).resolves.toBe(true)
        expect(client.token).toBe('new-token')

        FetchMock.fn.mockImplementation(async () => okResponse({id: 'user-1'}))
        await client.getClientConfig()
        await client.register('user@example.com', 'user', 'password', 'invite')
        await client.changePassword('user/1', 'old', 'new')
        await client.getMe()
        await client.getMyBoardMemberships()
        await client.getUser('user/1')
        await client.getUsersList(['user-1'])
        await client.getMyConfig()
        await client.patchUserConfig('user/1', {category: 'boards', name: 'welcomePageViewed', value: 'true'} as any)
        await client.logout()

        const calls = FetchMock.fn.mock.calls.map(([url]) => String(url))
        expect(calls).toEqual(expect.arrayContaining([
            'http://focalboard.test/api/v2/login',
            'http://focalboard.test/api/v2/register',
            'http://focalboard.test/api/v2/users/user%2F1/changepassword',
            'http://focalboard.test/api/v2/users/me?teamID=team-1&channelID=channel-1',
            'http://focalboard.test/api/v2/logout',
        ]))
    })

    test('returns safe defaults for failed and malformed API responses', async () => {
        FetchMock.fn.mockImplementation(async () => new Response('invalid-json', {status: 500}))

        await expect(client.login('user', 'bad')).resolves.toBe(false)
        await expect(client.logout()).resolves.toBe(false)
        await expect(client.getClientConfig()).resolves.toBeNull()
        await expect(client.getMe()).resolves.toBeUndefined()
        await expect(client.getMyBoardMemberships()).resolves.toEqual([])
        await expect(client.getUser('missing')).resolves.toBeUndefined()
        await expect(client.getUsersList(['missing'])).resolves.toEqual([])
        await expect(client.getMyConfig()).resolves.toBeUndefined()
        await expect(client.patchUserConfig('missing', {} as any)).resolves.toBeUndefined()
        await expect(client.getBoards()).resolves.toEqual([])
        await expect(client.getBoard('missing')).resolves.toBeUndefined()
        await expect(client.getSharing('missing')).resolves.toBeUndefined()
        await expect(client.setSharing('missing', {} as any)).resolves.toBe(false)
    })

    test('covers block, board, membership, sharing, and archive write contracts', async () => {
        const block = {id: 'block-1', boardId: 'board-1', type: 'card', title: 'Card', fields: {}} as any
        const board = {id: 'board-1', teamId: 'team-1', type: 'board', title: 'Board', fields: {}} as any
        const member = {boardId: 'board-1', userId: 'user-1', roles: 'editor'} as any

        await client.exportBoardArchive('board-1')
        await client.exportFullArchive('team-1')
        await client.importFullArchive(new File(['archive'], 'archive.boardarchive'))
        await client.getBlocksWithParent('parent/1')
        await client.getBlocksWithParent('parent/1', 'card')
        await client.getBlocksWithType('card')
        await client.getBlocksWithBlockID('block-1', 'board-1', 'read-token')
        await client.getAllBlocks('board-1')
        await client.patchBlock('board-1', 'block-1', {title: 'Changed'} as any)
        await client.patchBlocks([block], [{title: 'Changed'}] as any)
        await client.deleteBlock('board-1', 'block/1')
        await client.undeleteBlock('board/1', 'block/1')
        await client.undeleteBoard('board-1')
        await client.followBlock('block-1', 'card', 'user-1')
        await client.unfollowBlock('block-1', 'card', 'user-1')
        await client.insertBlock('board-1', block)
        await client.insertBlocks('board-1', [block], 'source/board')
        await client.createBoardsAndBlocks({boards: [board], blocks: [block]})
        await client.deleteBoardsAndBlocks(['board-1'], ['block-1'])
        await client.createBoardMember(member)
        await client.joinBoard('board-1', true)
        await client.joinBoard('board-1', false)
        await client.updateBoardMember(member)
        await client.deleteBoardMember(member)
        await client.patchBoardsAndBlocks({boardIDs: ['board-1'], blockIDs: ['block-1'], boardPatches: [], blockPatches: []} as any)
        await client.getSharing('board-1')
        await client.setSharing('board-1', {enabled: true} as any)
        await client.regenerateTeamSignupToken()
        await client.createBoard(board)
        await client.patchBoard('board-1', {title: 'Changed'} as any)
        await client.deleteBoard('board-1')
        await client.moveBlockTo('block-1', 'after', 'block-2')

        expect(FetchMock.fn).toHaveBeenCalledTimes(32)
        expect(FetchMock.fn).toHaveBeenCalledWith(
            'http://focalboard.test/api/v2/boards/board-1/blocks?sourceBoardID=source%2Fboard',
            expect.objectContaining({method: 'POST'}),
        )
    })

    test('covers team, category, search, onboarding, limits, and insights endpoints', async () => {
        const category = {id: 'category-1', teamID: 'team-1', name: 'Work'} as any

        await client.getTeam()
        await client.getTeams()
        await client.getTeamUsers()
        await client.getTeamUsers(true)
        await client.getTeamUsersList(['user-1'], 'team-2')
        await client.searchTeamUsers('marco')
        await client.searchTeamUsers('marco', true)
        await client.getTeamTemplates()
        await client.getTeamTemplates('team-2')
        await client.getBoards()
        await client.getBoard('board-1')
        await client.duplicateBoard('board-1', false)
        await client.duplicateBoard('board-1', true, 'team/2')
        await client.duplicateBlock('board-1', 'block-1', false)
        await client.duplicateBlock('board-1', 'block-1', true)
        await client.getBlocksForBoard('team-2', 'board-1')
        await client.getBoardMembers('team-1', 'board-1')
        await client.getSidebarCategories('team-1')
        await client.createSidebarCategory(category)
        await client.updateSidebarCategory(category)
        await client.deleteSidebarCategory('team-1', 'category-1')
        await client.reorderSidebarCategories('team-1', ['category-1'])
        await client.reorderSidebarCategoryBoards('team-1', 'category-1', ['board-1'])
        await client.moveBoardToCategory('team-1', 'board-1', 'category-1', 'category-0')
        await client.hideBoard('category-1', 'board-1')
        await client.unhideBoard('category-1', 'board-1')
        await client.search('team-1', 'road map')
        await client.searchLinkableBoards('team-1', 'road map')
        await client.searchAll('road map')
        await client.getUserBlockSubscriptions('user-1')
        await client.searchUserChannels('team-1', 'town square')
        await client.getChannel('team-1', 'channel-1')
        await client.prepareOnboarding('team-1')
        await client.notifyAdminUpgrade()
        await client.getBoardsCloudLimits()
        await client.getSiteStatistics()
        await client.getMyTopBoards('week', 1, 10, 'team-1')
        await client.getTeamTopBoards('week', 1, 10, 'team-1')

        expect(FetchMock.fn).toHaveBeenCalledTimes(38)
        expect(FetchMock.fn).toHaveBeenCalledWith(
            'http://focalboard.test/api/v2/teams/team-1/boards/search?q=road%20map',
            expect.objectContaining({method: 'GET'}),
        )
    })

    test('uploads a file and handles invalid upload responses', async () => {
        const file = new File(['image'], 'image.png')
        FetchMock.fn.mockImplementationOnce(async () => okResponse({fileId: 'file-1'}))
        await expect(client.uploadFile('board-1', file)).resolves.toBe('file-1')

        FetchMock.fn.mockImplementationOnce(async () => new Response('not-json', {status: 200}))
        await expect(client.uploadFile('board-1', file)).resolves.toBeUndefined()

        FetchMock.fn.mockImplementationOnce(async () => new Response('', {status: 500}))
        await expect(client.uploadFile('board-1', file)).resolves.toBeUndefined()
    })
})
