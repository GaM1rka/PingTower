import { toast } from "react-toastify";

type ToastType = "success" | "error" | "info";

export const useToast = () => {
  const show = (message: string, type: ToastType = "info") => {
    toast(message, { type });
  };

  return { show };
};