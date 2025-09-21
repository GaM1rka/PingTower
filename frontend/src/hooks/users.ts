import API from "./api";
import 'react-toastify/dist/ReactToastify.css';
import { useToast } from "./toast";


//хз как будут апишки называться, пока наугад тыкаю
//потом поменять

const handleRegistration = async (email: string, password: string) => {
     const { show } = useToast();
    try {
        const res = await API.post("/register", {email, password});
        console.log("registered");
        //нет апи, нет ответа. data.token может отличаться (поменять везде)
        localStorage.setItem("token", res.data.access_token);
        localStorage.setItem("email",email);
        show("success!", "success");
    }
    catch (e) {
        console.error(`ERROR: ${e}`);
        show(`ERROR: ${e}`, "error");
    }
}

const handleLogin = async (email: string, password: string) => {
    const { show } = useToast();
    try {
        const res = await API.post("/login", {email, password});
        console.log("logged in");
        localStorage.setItem("token", res.data.access_token);
        localStorage.setItem("email",email);
        show("success!", "success");
    }
    catch (e) {
        console.error(`ERROR: ${e}`);
        show(`ERROR: ${e}`, "error");
    }
}

const handleLogout = async () => {
    try {
        await API.post("/users/logout");
        console.log("logout");
        localStorage.removeItem("token");
    }
    catch (e) {
        console.error(`ERROR: ${e}`);
    }
}

export {handleRegistration, handleLogin, handleLogout};