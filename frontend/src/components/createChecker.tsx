import { useState } from "react";
import { createChecker } from "../hooks/checkers";
import "../index.css"
import { useToast } from "../hooks/toast";

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
    const { show } = useToast();
    const MMtime: number = (time as unknown) as number ;
    if (checkURL(site)) {
        createChecker({site, MMtime});
    } else {
        show("invalid URL", "error")
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
      <input placeholder="URL (https://example.com)" className="input" value={URL} onChange={(e) => setEmail(e.target.value)} />
      <input placeholder="Period (in minutes)" min={1} className="input" value={time} type="number" onChange={(e) => setPassword(e.target.value)} />
      <button className="input" onClick={() => create(URL, time)}>create checker</button>
    </div>
    </div>
    </>
  );
}