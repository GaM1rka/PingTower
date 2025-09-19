import { useState } from "react";
import { handleRegistration } from "../api/users"
import "../index.css"

export default function Registration() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    
    return (
    <div className="registration">
      <input value={email} onChange={(e) => setEmail(e.target.value)} />
      <input value={password} type="password" onChange={(e) => setPassword(e.target.value)} />
      <button onClick={() => handleRegistration(email,password)}>registration</button>
    </div>
  );
}