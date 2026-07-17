// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render, screen, fireEvent, waitFor} from '@testing-library/react'
import {Router, Route} from 'react-router-dom'
import {createMemoryHistory} from 'history'
import {mocked} from 'jest-mock'

import {wrapIntl} from '../../testUtils'
import client from '../../octoClient'
import LoginPage from '../loginPage'
import {fetchMe} from '../../store/users'
import {useAppSelector} from '../../store/hooks'

jest.mock('../../octoClient')
const mockedOctoClient = mocked(client, true)

const mockDispatch: jest.Mock = jest.fn((action) => {
    if (typeof action === 'function') {
        return action(mockDispatch, () => ({}), undefined)
    }
    return action
})

jest.mock('../../store/hooks', () => ({
    useAppDispatch: () => mockDispatch,
    useAppSelector: jest.fn(),
}))

jest.mock('../../store/users', () => ({
    fetchMe: jest.fn(() => () => Promise.resolve()),
    getLoggedIn: jest.fn(),
}))

describe('pages/loginPage', () => {
    let history = createMemoryHistory()

    beforeEach(() => {
        jest.resetAllMocks()
        history = createMemoryHistory()
        mocked(useAppSelector).mockReturnValue(false)
    })

    test('redirects to root if already logged in', async () => {
        mocked(useAppSelector).mockReturnValue(true)

        render(
            wrapIntl(
                <Router history={history}>
                    <Route
                        exact={true}
                        path='/'
                    >
                        <div>{'Root Page'}</div>
                    </Route>
                    <Route path='/login'>
                        <LoginPage/>
                    </Route>
                </Router>,
            ),
        )

        history.push('/login')
        await waitFor(() => {
            expect(screen.getByText('Root Page')).toBeDefined()
        })
    })

    test('renders login form when logged out', () => {
        render(
            wrapIntl(
                <Router history={history}>
                    <LoginPage/>
                </Router>,
            ),
        )

        expect(screen.getByPlaceholderText('Enter username')).toBeDefined()
        expect(screen.getByPlaceholderText('Enter password')).toBeDefined()
        expect(screen.getByRole('button', {name: 'Log in'})).toBeDefined()
        expect(screen.getByText('or create an account if you don\'t have one')).toBeDefined()
    })

    test('handles input changes', () => {
        render(
            wrapIntl(
                <Router history={history}>
                    <LoginPage/>
                </Router>,
            ),
        )

        const usernameInput = screen.getByPlaceholderText('Enter username') as HTMLInputElement
        const passwordInput = screen.getByPlaceholderText('Enter password') as HTMLInputElement

        fireEvent.change(usernameInput, {target: {value: 'testuser'}})
        fireEvent.change(passwordInput, {target: {value: 'testpass'}})

        expect(usernameInput.value).toBe('testuser')
        expect(passwordInput.value).toBe('testpass')
    })

    test('handles successful login and redirects to root by default', async () => {
        mockedOctoClient.login.mockResolvedValue(true)
        history.push = jest.fn()

        render(
            wrapIntl(
                <Router history={history}>
                    <LoginPage/>
                </Router>,
            ),
        )

        fireEvent.change(screen.getByPlaceholderText('Enter username'), {target: {value: 'user'}})
        fireEvent.change(screen.getByPlaceholderText('Enter password'), {target: {value: 'pass'}})
        fireEvent.click(screen.getByRole('button', {name: 'Log in'}))

        await waitFor(() => {
            expect(mockedOctoClient.login).toHaveBeenCalledWith('user', 'pass')
            expect(fetchMe).toHaveBeenCalled()
            expect(history.push).toHaveBeenCalledWith('/')
        })
    })

    test('handles successful login and redirects to r query parameter if present', async () => {
        mockedOctoClient.login.mockResolvedValue(true)
        history.push = jest.fn()
        history.location.search = '?r=/redirect-path'

        render(
            wrapIntl(
                <Router history={history}>
                    <LoginPage/>
                </Router>,
            ),
        )

        fireEvent.change(screen.getByPlaceholderText('Enter username'), {target: {value: 'user'}})
        fireEvent.change(screen.getByPlaceholderText('Enter password'), {target: {value: 'pass'}})
        fireEvent.click(screen.getByRole('button', {name: 'Log in'}))

        await waitFor(() => {
            expect(mockedOctoClient.login).toHaveBeenCalledWith('user', 'pass')
            expect(history.push).toHaveBeenCalledWith('/redirect-path')
        })
    })

    test('handles failed login and displays error message', async () => {
        mockedOctoClient.login.mockResolvedValue(false)

        render(
            wrapIntl(
                <Router history={history}>
                    <LoginPage/>
                </Router>,
            ),
        )

        fireEvent.change(screen.getByPlaceholderText('Enter username'), {target: {value: 'user'}})
        fireEvent.change(screen.getByPlaceholderText('Enter password'), {target: {value: 'pass'}})
        fireEvent.click(screen.getByRole('button', {name: 'Log in'}))

        await waitFor(() => {
            expect(screen.getByText('Login failed')).toBeDefined()
        })

        // typing clears error
        fireEvent.change(screen.getByPlaceholderText('Enter username'), {target: {value: 'newuser'}})
        expect(screen.queryByText('Login failed')).toBeNull()
    })
})
