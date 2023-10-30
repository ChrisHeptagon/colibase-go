import Sidebar from "./components/sidebar";
import "./homePage.scss";

export default function HomePage() {
    return (
        <div className="home-page">
            <title>Colibase - Dashboard</title>
            <Sidebar />
            <div className="main-content">
                <h1
                >
                    Dashboard
                    </h1>
                </div>
        </div>
    );
    }