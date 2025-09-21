import {  useState } from "react";
import Registration from "./registration";

export default function AppBar() {
    const [menu, setMenu] = useState(false);

    return(
        <div className="appbar">
            <div className="divapp">
                {menu && (
                <>
                <Registration/>
                </>
            )}
            <div className="menu" onClick={() => setMenu(!menu)}>
            menu
            </div>
            </div>
        </div>
    );
}


        