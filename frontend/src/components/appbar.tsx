import { useEffect, useState } from "react";
import Registration from "./registration";

export default function AppBar() {
    const [showHello, setShowHello] = useState(false);

  useEffect(() => {
    const value = localStorage.getItem("email");
    if (value) setShowHello(true);
  }, []);

  
    return (
    <div className="appbar">
      {showHello ? (
        <h1>Привет, {localStorage.getItem("email")}</h1>
      ) : (
        <Registration />
      )}
    </div>
    );
}


        