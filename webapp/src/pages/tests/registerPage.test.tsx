// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {render, screen, fireEvent, waitFor} from '@testing-library/react'
import {Router, Route} from 'react-router-dom'
import {createMemoryHistory} from 'history'
import {mocked} from 'jest-mock'

import {wrapIntl} from '../../testUtils'
import client from '../../octoClient'
import RegisterPage from '../registerPage'
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

describe('pages/registerPage', () => {
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
                    <Route path='/register'>
                        <RegisterPage/>
                    </Route>
                </Router>,
            ),
        )

        history.push('/register')
        await waitFor(() => {
            expect(screen.getByText('Root Page')).toBeDefined()
        })
    })

    test('renders registration form when logged out', () => {
        render(
            wrapIntl(
                <Router history={history}>
                    <RegisterPage/>
                </Router>,
            ),
        )

        expect(screen.getByPlaceholderText('Enter email')).toBeDefined()
        expect(screen.getByPlaceholderText('Enter username')).toBeDefined()
        expect(screen.getByPlaceholderText('Enter password')).toBeDefined()
        expect(screen.getByRole('button', {name: 'Register'})).toBeDefined()
        expect(screen.getByText('or log in if you already have an account')).toBeDefined()
    })

    test('handles input changes', () => {
        render(
            wrapIntl(
                <Router history={history}>
                    <RegisterPage/>
                </Router>,
            ),
        )

        const emailInput = screen.getByPlaceholderText('Enter email') as HTMLInputElement
        const usernameInput = screen.getByPlaceholderText('Enter username') as HTMLInputElement
        const passwordInput = screen.getByPlaceholderText('Enter password') as HTMLInputElement

        fireEvent.change(emailInput, {target: {value: 'test@email.com'}})
        fireEvent.change(usernameInput, {target: {value: 'testuser'}})
        fireEvent.change(passwordInput, {target: {value: 'testpass'}})

        expect(emailInput.value).toBe('test@email.com')
        expect(usernameInput.value).toBe('testuser')
        expect(passwordInput.value).toBe('testpass')
    })

    test('handles successful registration and login', async () => {
        mockedOctoClient.register.mockResolvedValue({
            code: 200,
            json: {},
        })
        mockedOctoClient.login.mockResolvedValue(true)
        history.push = jest.fn()

        render(
            wrapIntl(
                <Router history={history}>
                    <RegisterPage/>
                </Router>,
            ),
        )

        fireEvent.change(screen.getByPlaceholderText('Enter email'), {target: {value: 'test@email.com'}})
        fireEvent.change(screen.getByPlaceholderText('Enter username'), {target: {value: 'user'}})
        fireEvent.change(screen.getByPlaceholderText('Enter password'), {target: {value: 'pass'}})
        fireEvent.click(screen.getByRole('button', {name: 'Register'}))

        await waitFor(() => {
            expect(mockedOctoClient.register).toHaveBeenCalledWith('test@email.com', 'user', 'pass', '')
            expect(mockedOctoClient.login).toHaveBeenCalledWith('user', 'pass')
            expect(fetchMe).toHaveBeenCalled()
            expect(history.push).toHaveBeenCalledWith('/')
        })
    })

    test('handles signupToken query parameter', async () => {
        const originalLocation = window.location
        delete (window as any).location;
        (window as any).location = {search: '?t=token123'}

        mockedOctoClient.register.mockResolvedValue({
            code: 200,
            json: {},
        })
        mockedOctoClient.login.mockResolvedValue(true)

        render(
            wrapIntl(
                <Router history={history}>
                    <RegisterPage/>
                </Router>,
            ),
        )

        fireEvent.change(screen.getByPlaceholderText('Enter email'), {target: {value: 'test@email.com'}})
        fireEvent.change(screen.getByPlaceholderText('Enter username'), {target: {value: 'user'}})
        fireEvent.change(screen.getByPlaceholderText('Enter password'), {target: {value: 'pass'}})
        fireEvent.click(screen.getByRole('button', {name: 'Register'}))

        await waitFor(() => {
            expect(mockedOctoClient.register).toHaveBeenCalledWith('test@email.com', 'user', 'pass', 'token123')
        })

        window.location = originalLocation
    })

    test('handles 401 registration error link invalid', async () => {
        mockedOctoClient.register.mockResolvedValue({
            code: 401,
            json: {},
        })

        render(
            wrapIntl(
                <Router history={history}>
                    <RegisterPage/>
                </Router>,
            ),
        )

        fireEvent.click(screen.getByRole('button', {name: 'Register'}))

        await waitFor(() => {
            expect(screen.getByText('Invalid registration link, please contact your administrator')).toBeDefined()
        })
    })

    test('handles other registration error messages', async () => {
        mockedOctoClient.register.mockResolvedValue({
            code: 400,
            json: {error: 'Username already taken'},
        })

        render(
            wrapIntl(
                <Router history={history}>
                    <RegisterPage/>
                </Router>,
            ),
        )

        fireEvent.click(screen.getByRole('button', {name: 'Register'}))

        await waitFor(() => {
            expect(screen.getByText('Username already taken')).toBeDefined()
        })
    })
})
