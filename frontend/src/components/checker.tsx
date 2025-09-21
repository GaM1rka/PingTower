import { useMemo, useState } from "react";
import { useChecker } from "../hooks/checkers";

type CheckerProps = {
  id: number;
  url: string;
  status?: "ok" | "bad" | "initial";
};

function formatTime(iso: string) {
  try {
    const d = new Date(iso);
    return d.toLocaleString();
  } catch {
    return iso;
  }
}

export default function Checker({ id, url, status: initialStatus }: CheckerProps) {
  const [expanded, setExpanded] = useState(false);
  const { checker: logs, loading, error, } = useChecker(id);
  
  const latestStatus = useMemo<"ok" | "bad" | "initial" | undefined>(() => {
    if (logs && logs.length > 0) return logs[0].status;
    return initialStatus;
  }, [logs, initialStatus]);

  return (
    <div className={`checker ${expanded ? "expanded" : ""} ${latestStatus}`}>
        <button
          type="button"
          className="input in3"
          onClick={() => setExpanded((v) => !v)}
          aria-expanded={expanded}
          aria-controls={`checker-body-${id}`}
        >
          {expanded ? "close" : "open"}
        </button>

      {/* поменять потом или накидать стилей */}
      <div className="expandedInfo">
        <a
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          className="link"
          onClick={(e) => e.stopPropagation()}
        >
          {url}
        </a>

        <span
          className={`checker__badge checker__badge--${
            latestStatus ?? "initial"
          }`}
          title={`status: ${latestStatus ?? "initial"}`}
        >
          status: {latestStatus ?? "initial"}
        </span>
      {/* </div> */}

      {expanded && (
        <div id={`checker-body-${id}`} className="checker__body">
          {loading && <p>loading...</p>}
          {error && <p>error: {error}</p>}
          {!loading && !error && logs && logs.length > 0 && (
            <table className="checker__table">
              <thead>
                <tr>
                  <th>ping time</th>
                  <th>response time (ms)</th>
                  <th>status</th>
                </tr>
              </thead>
              <tbody>
                {logs.map((row, idx) => (
                  <tr key={`${row.req_time}-${idx}`}>
                    <td>{formatTime(row.req_time)}</td>
                    <td>{row.resp_time}</td>
                    <td>{row.status}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {!loading && !error && logs && logs.length === 0 && (
            <p>No logs yet.</p>
          )}
        </div>
      )}
      </div>
    </div>
  );
}
