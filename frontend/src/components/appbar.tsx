import { useState } from "react";
import Registration from "./registration";

export default function AppBar() {

    const [menu, setMenu] = useState(false);
    return(
        <div className="appbar">
            {menu && (
                <Registration/>
            )}
            <div className="menu" onClick={() => setMenu(!menu)}>
            menu
            </div>
        </div>
    );
}


        