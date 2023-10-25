/* @refresh reload */
import { render } from 'solid-js/web'
import { Route, Router, Routes } from '@solidjs/router'
import LoginPage from './LoginPage'
import InitLoginPage from './InitLoginPage'

const root = document.getElementById('root')

render(() => (
    <Router>
        <Routes>
            <Route path="/admin-entry/login" element={<LoginPage />} />
            <Route path="/admin-entry/init" element={<InitLoginPage />} />
        </Routes>
    </Router>
),root!)
