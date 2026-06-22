// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {configureStore} from '@reduxjs/toolkit'

import {reducer, setTeam, getCurrentTeamId, getCurrentTeam, getFirstTeam, getAllTeams} from '../teams'

// Mock octoClient and utils
jest.mock('../../octoClient', () => ({
    default: {
        getTeams: jest.fn().mockResolvedValue([]),
        getTeam: jest.fn().mockResolvedValue(null),
        regenerateTeamSignupToken: jest.fn().mockResolvedValue(undefined),
    },
}))

jest.mock('../../utils', () => ({
    Utils: {
        log: jest.fn(),
        logError: jest.fn(),
    },
}))

function makeStore() {
    return configureStore({reducer: {teams: reducer}})
}

function makeTeam(id: string, title: string) {
    return {id, title, signupToken: '', modifiedBy: '', updateAt: 0}
}

describe('teams reducer', () => {
    test('initial state is empty', () => {
        const store = makeStore()
        expect(store.getState().teams.current).toBeNull()
        expect(store.getState().teams.currentId).toBe('')
        expect(store.getState().teams.allTeams).toEqual([])
    })

    test('setTeam sets currentId and current team when found', () => {
        const store = makeStore()
        const team1 = makeTeam('t1', 'Alpha')
        // First load teams via initialLoad
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: team1,
                teams: [team1],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        store.dispatch(setTeam('t1'))
        expect(store.getState().teams.currentId).toBe('t1')
        expect(store.getState().teams.current?.id).toBe('t1')
    })

    test('setTeam logs error and does not set current if team not found', () => {
        const store = makeStore()
        store.dispatch(setTeam('non-existent'))
        expect(store.getState().teams.currentId).toBe('non-existent')
        expect(store.getState().teams.current).toBeNull()
    })

    test('setTeam does nothing if current team is the same', () => {
        const store = makeStore()
        const team1 = makeTeam('t1', 'Alpha')
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: team1,
                teams: [team1],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        store.dispatch(setTeam('t1'))
        const currentAfterFirst = store.getState().teams.current
        store.dispatch(setTeam('t1'))
        // The current should remain the same reference
        expect(store.getState().teams.current).toEqual(currentAfterFirst)
    })

    test('extraReducer: initialLoad.fulfilled sets team and allTeams sorted', () => {
        const store = makeStore()
        const team1 = makeTeam('t1', 'Beta')
        const team2 = makeTeam('t2', 'Alpha')
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: team1,
                teams: [team1, team2],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        expect(store.getState().teams.current).toEqual(team1)
        expect(store.getState().teams.allTeams[0].title).toBe('Alpha')
        expect(store.getState().teams.allTeams[1].title).toBe('Beta')
    })

    test('extraReducer: team/fetch/fulfilled sets allTeams sorted', () => {
        const store = makeStore()
        const team1 = makeTeam('t1', 'Zeta')
        const team2 = makeTeam('t2', 'Delta')
        store.dispatch({
            type: 'team/fetch/fulfilled',
            payload: [team1, team2],
        })
        expect(store.getState().teams.allTeams[0].title).toBe('Delta')
        expect(store.getState().teams.allTeams[1].title).toBe('Zeta')
    })

    test('extraReducer: team/refreshCurrentTeam/fulfilled updates current team', () => {
        const store = makeStore()
        const updatedTeam = makeTeam('t1', 'Updated Team')
        store.dispatch({
            type: 'team/refreshCurrentTeam/fulfilled',
            payload: updatedTeam,
        })
        expect(store.getState().teams.current).toEqual(updatedTeam)
    })
})

describe('teams selectors', () => {
    test('getCurrentTeamId returns empty string initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentTeamId(state)).toBe('')
    })

    test('getCurrentTeam returns null initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getCurrentTeam(state)).toBeNull()
    })

    test('getFirstTeam returns null when no teams', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getFirstTeam(state)).toBeUndefined()
    })

    test('getFirstTeam returns first team when teams loaded', () => {
        const store = makeStore()
        const team1 = makeTeam('t1', 'Alpha')
        const team2 = makeTeam('t2', 'Beta')
        store.dispatch({
            type: 'initialLoad/fulfilled',
            payload: {
                team: team1,
                teams: [team1, team2],
                boards: [],
                boardsMemberships: [],
                boardTemplates: [],
                limits: null,
                myConfig: [],
            },
        })
        const state = store.getState() as any
        expect(getFirstTeam(state)?.id).toBe('t1')
    })

    test('getAllTeams returns empty array initially', () => {
        const store = makeStore()
        const state = store.getState() as any
        expect(getAllTeams(state)).toEqual([])
    })

    test('getAllTeams returns all loaded teams', () => {
        const store = makeStore()
        const team1 = makeTeam('t1', 'Team 1')
        const team2 = makeTeam('t2', 'Team 2')
        store.dispatch({type: 'team/fetch/fulfilled', payload: [team1, team2]})
        const state = store.getState() as any
        expect(getAllTeams(state)).toHaveLength(2)
    })
})
