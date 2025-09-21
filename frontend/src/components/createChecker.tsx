import { useState } from "react";
import { createChecker } from "../api/checkers";
import "../index.css"

// function convertTime(time: string): number {
//   const [hours, minutes] = time.split(':').map(Number);
//   if (isNaN(hours) || isNaN(minutes)) throw new Error('Invalid time format');
//   return hours * 60 + minutes;
// }

function checkURL(URL: string) {
    const urlRegex = /^(https?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w-./?%&=]*)?$/i;
    return urlRegex.test(URL);
}

function create(site:string,time:string) {
    const MMtime: number = (time as unknown) as number ;
    if (checkURL(site)) {
        createChecker({site, MMtime});
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
      <input placeholder="Period" className="input" value={time} type="number" onChange={(e) => setPassword(e.target.value)} />
      <button className="input" onClick={() => create(URL, time)}>create checker</button>
    </div>
    </div>
    </>
  );
}