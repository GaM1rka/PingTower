import { useState, useEffect, useCallback } from "react";
import API from "./api";
import { useToast } from "./toast";

type CheckerLog = {
  id: number;         // site id
  req_time: string;   // ISO string
  resp_time: number;  // ms, can be -1
  status: "ok" | "bad" | "initial";
  site: string;       // url
};

type CreateCheckerPayload = {
  site: string;
  MMtime: number; // seconds
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
  const [checker, setChecker] = useState<CheckerLog[] | null>([]);
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

export const createChecker = async ({ site, MMtime }: CreateCheckerPayload) => {
  const { show } = useToast();
  // Validate a bit on the client:
  if (!site) throw new Error("Site is required");
  const res = await API.post("/checkers", { site, MMtime });
  show("Checker created!", "success")
  console.log('success!')
  return res.data; // let caller update UI or re-fetch; don't reload the page
};
