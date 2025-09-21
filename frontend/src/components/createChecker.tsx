import { useState } from "react";
import { createChecker } from "../api/checkers";
import "../index.css"


function checkURL(URL: string) {
    const urlRegex = /^(https?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w-./?%&=]*)?$/i;
    return urlRegex.test(URL);
}

function create(site:string,time:string) {
    if (checkURL(site)) {
        createChecker({site,time});
    } else {
        console.error("invalid URL")
    }
}

export default function Createchecker() {
    const [URL, setEmail] = useState("");
    const [time, setPassword] = useState("");
    
    return (
    <>
    <div style={{ display: "flex" }}>
    <div className="createChecker">
      <input placeholder="URL" className="input" value={URL} onChange={(e) => setEmail(e.target.value)} />
      <input placeholder="Period" className="input" value={time} type="time" onChange={(e) => setPassword(e.target.value)} />
      <button className="input" onClick={() => create(URL, time)}>create checker</button>
    </div>
    </div>
    </>
  );
}