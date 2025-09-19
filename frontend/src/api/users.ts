import API from "./api";


//хз как будут апишки называться, пока наугад тыкаю
//потом поменять

const handleRegistration = async (email: string, password: string) => {
    try {
        const res = await API.post("/authorize", {email, password});
        console.log("registered");
        //нет апи, нет ответа. data.token может отличаться (поменять везде)
        localStorage.setItem("token", res.data.token);
    }
    catch (e) {
        console.error(`ERROR: ${e}`);
    }
}

const handleLogin = async (email: string, password: string) => {
    try {
        const res = await API.post("/users/login", {email, password});
        console.log("registered");
        localStorage.setItem("token", res.data.token);
    }
    catch (e) {
        console.error(`ERROR: ${e}`);
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