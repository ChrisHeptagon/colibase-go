import { useEffect } from "react";
import "./LoginPage.scss";
import LoginForm from "./components/loginComponent";

const LoginPage = () => {
  useEffect(() => {
     fetch("/api/user-initialization-status").then((res) => {
      if (res.status === 500) {
        return JSON.stringify({ status: "User not initialized" });
      }
      if (res.status === 200) {
        window.location.href = "/admin-entry/login";
      }
    })
  }, []);

 
  return (
    <>
      <title>Colibase - Initial Login Page</title>
      <div className="background">
        <LoginForm
        headerText="Initial Login"
        />
      </div>
    </>
  );
};

export default LoginPage;
