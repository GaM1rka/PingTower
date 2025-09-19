import { useState } from "react";
import { createChecker } from "../api/checkers";
import "../index.css"


function checkURL(URL: string) {
    const urlRegex = /^(https?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w-./?%&=]*)?$/i;
    return urlRegex.test(URL);
}

function create(URL:string,time:string) {
    if (checkURL(URL)) {
        createChecker(URL,time);
    } else {
        console.error("invalid URL")
    }
}

export default function Createchecker() {
    const [URL, setEmail] = useState("");
    const [time, setPassword] = useState("");
    const [open, setOpen] = useState(false);
    
    return (
    <>
    <div style={{ display: "flex" }}>
    <div className="createChecker" onClick={() => setOpen(!open)}></div>
    {open &&(<div className="createChecker expanded">
      <input value={URL} onChange={(e) => setEmail(e.target.value)} />
      <input value={time} type="time" onChange={(e) => setPassword(e.target.value)} />
      <button onClick={() => create(URL, time)}>create checker</button>
    </div>)}
    </div>
    </>
  );
}