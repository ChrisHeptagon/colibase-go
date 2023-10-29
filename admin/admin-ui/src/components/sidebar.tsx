import './sidebar.scss';
import {
    TbDoorExit
} from 'react-icons/tb';
import {
    HiMiniCog8Tooth
} from 'react-icons/hi2';

import {
    MdDashboard
} from 'react-icons/md';
import ColiBaseLogo from '../assets/colibase_logo.webp';


export default function Sidebar() {
    return (
        <aside className="sidebar">
            <span className="sidebar-logo"
            aria-label="logo"
            role="img"
            >
                <a href="/admin-ui/dashboard">
                    < img src={ColiBaseLogo} alt="C" />
                </a>
            </span>
            <div className="sidebar-menu">
                <ul>
                    <li><a href="/admin-ui/dashboard">
                            <MdDashboard />
                        Dashboard</a></li>
                    <li><a href="/admin-ui/settings">
                            <HiMiniCog8Tooth />
                        Settings</a></li>
                    <li><a href="/admin-entry/logout">
                            <TbDoorExit />
                        Logout
                    </a></li>
                </ul>
            </div>
        </aside>
    );
}