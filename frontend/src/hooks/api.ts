import axios from "axios";

const API = axios.create({
    //поменяю когда будет готов бэк
    //ЕСЛИ будет готов
    //чупеп
    baseURL: "http://localhost:8080"
});

API.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config
})

export default API