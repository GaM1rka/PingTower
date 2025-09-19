import { useState } from "react";
import {useChecker} from "../api/checkers"

type CheckerProps = {
  id: number;
  url: string;
  status: boolean;
};

export default function Checker({ id, url, status }: CheckerProps) {
    const [expanded, setExpanded] = useState(false);
    return (
        <div className={expanded ? "checker expanded" : "checker"} 
            onClick={() => {
            setExpanded(!expanded);
            const info = useChecker(id)
            }}>
            {expanded && (
                <div>
                    <text>{url}</text>
                </div>
            )}
        </div>
    );

}