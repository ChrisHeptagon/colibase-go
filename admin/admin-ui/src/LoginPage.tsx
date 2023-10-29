import { useEffect } from "react";
import "./LoginPage.scss";
import LoginForm from "./components/loginComponent";



const LoginPage = () => {
  useEffect(() => {
    fetch("/api/user-initialization-status").then((res) => {
      if (res.status === 500) {
        window.location.href = "/admin-entry/init";
        return JSON.stringify({ status: "User not initialized" });
      }
      if (res.status === 200) {
        return JSON.stringify({ status: "User initialized" });
      }
    });
  }, []);
 
  return (
    <>
      <title>Colibase - Login Page</title>
      <div className="background">
        <LoginForm
        headerText="Login"
        />
      </div>
    </>
  );
};

export default LoginPage;
