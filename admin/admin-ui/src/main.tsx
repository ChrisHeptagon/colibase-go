/* @refresh reload */
import { Route, BrowserRouter, Routes } from 'react-router-dom'
import LoginPage from './LoginPage'
import InitLoginPage from './InitLoginPage'
import ReactDOM from 'react-dom/client'


ReactDOM.createRoot(document.getElementById('root')!).render(
        <BrowserRouter>
        <Routes>
            <Route path="/admin-entry/login" element={<LoginPage />} />
            <Route path="/admin-entry/init" element={<InitLoginPage />} />
        </Routes>
        </BrowserRouter>
)