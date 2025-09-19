import { useState, useEffect } from "react";
import API from "./api";

export const useCheckers = () => {
  const [checkers, setCheckers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const getCheckers = async () => {
    try {
      const res = await API.get("/checker");
      setCheckers(res.data);
    } catch (e: any) {
      setError(e.message || "Error fetching checkers");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getCheckers();
  }, []);

  return { checkers, loading, error, refresh: getCheckers };
};

export const useChecker = (id: number) => {
  const [checker, setChecker] = useState<any | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const getChecker = async () => {
    try {
      const res = await API.get(`/checker/${id}`);
      setChecker(res.data);
    } catch (e: any) {
      setError(e.message || "Error fetching checker");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getChecker();
  }, [id]);

  return { checker, loading, error, refresh: getChecker };
};

export const createChecker = async (URL: string, time: string) => {
  try {
    await API.post(`/checker`, { URL, time });
    console.log("added");
    window.location.reload();
  } catch (e: any) {
    console.error(`ERROR: ${e}`);
  }
};
