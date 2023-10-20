/* @refresh reload */
import { render } from 'solid-js/web'
import { Route, Router, Routes } from '@solidjs/router'
import LoginPage from './LoginPage'

const root = document.getElementById('root')

render(() => (
    <Router>
        <Routes>
            <Route path="/admin-ui/login" element={<LoginPage />} />
        </Routes>
    </Router>
),root!)
