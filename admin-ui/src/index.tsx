/* @refresh reload */
import { render } from 'solid-js/web'
import { Route, Router, Routes } from '@solidjs/router'
import LoginPage from './LoginPage'
import InitLoginPage from './InitLoginPage'

const root = document.getElementById('root')

render(() => (
    <Router>
        <Routes>
            <Route path="/admin-ui/login" element={<LoginPage />} />
            <Route path="/admin-ui/init" element={<InitLoginPage />} />
        </Routes>
    </Router>
),root!)
