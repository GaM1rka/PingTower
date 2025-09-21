import { useState } from "react";
import { handleRegistration } from "../api/users"
import "../index.css"

export default function Registration() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    
    return (
    <div className="registration">
      <input placeholder="email" className="input in2" value={email} onChange={(e) => setEmail(e.target.value)} />
      <input placeholder="password" className="input in2" value={password} type="password" onChange={(e) => setPassword(e.target.value)} />
      <button className="input in2" onClick={() => handleRegistration(email,password)}>registration</button>
    </div>
  );
}