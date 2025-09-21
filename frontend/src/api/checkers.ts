import { useState, useEffect, useCallback } from "react";
import API from "./api";

type CheckerLog = {
  id: number;         // site id
  req_time: string;   // ISO string
  resp_time: number;  // ms, can be -1
  status: "ok" | "bad" | "initial";
  site: string;       // url
};

type CreateCheckerPayload = {
  site: string;
  time: string; // seconds
};

export const useCheckers = () => {
  const [checkers, setCheckers] = useState<CheckerLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const getCheckers = useCallback(async (signal?: AbortSignal) => {
    try {
      setLoading(true);
      setError(null);
      const res = await API.get<CheckerLog[]>("/checkers", { signal });
      setCheckers(res.data);
    } catch (e: any) {
      if (e.name === "CanceledError") return;
      setError(e?.response?.data?.error ?? e.message ?? "Error fetching checkers");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    getCheckers(controller.signal);
    return () => controller.abort();
  }, [getCheckers]);

  return { checkers, loading, error, refresh: () => getCheckers() };
};

export const useChecker = (id?: number) => {
  const [checker, setChecker] = useState<CheckerLog[] | null>(null);
  const [loading, setLoading] = useState(!!id);
  const [error, setError] = useState<string | null>(null);

  const getChecker = useCallback(async (signal?: AbortSignal) => {
    if (!id) return;
    try {
      setLoading(true);
      setError(null);
      const res = await API.get<CheckerLog[]>(`/checker/${id}`, { signal });
      setChecker(res.data);
    } catch (e: any) {
      if (e.name === "CanceledError") return;
      setError(e?.response?.data?.error ?? e.message ?? "Error fetching checker");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    if (!id) return;
    const controller = new AbortController();
    getChecker(controller.signal);
    return () => controller.abort();
  }, [id, getChecker]);

  return { checker, loading, error, refresh: () => getChecker() };
};

export const createChecker = async ({ site, time }: CreateCheckerPayload) => {
  // Validate a bit on the client:
  if (!site) throw new Error("Site is required");
  if (typeof time !== "number" || time <= 0) throw new Error("Time must be a positive number (seconds)");

  const res = await API.post("/checkers", { site, time });
  return res.data; // let caller update UI or re-fetch; don't reload the page
};
