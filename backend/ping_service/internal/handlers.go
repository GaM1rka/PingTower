package internal

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"ping_service/models"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
	},
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req models.PingRequest
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}
	if req.Site == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "site is required"})
		return
	}

	target := req.Site
	if _, err := url.ParseRequestURI(target); err != nil || (!hasScheme(target)) {
		target = "https://" + req.Site
	}
	u, err := url.Parse(target)
	if err != nil || u.Host == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid site URL"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	pingTime := time.Now().UTC()

	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "cannot create request: " + err.Error()})
		return
	}

	start := time.Now()
	resp, err := httpClient.Do(reqHTTP)
	elapsed := time.Since(start)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, models.PingResponse{
			PingTime:     pingTime.Format(time.RFC3339Nano),
			ResponseTime: -1,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		writeJSON(w, http.StatusBadGateway, models.PingResponse{
			PingTime:     pingTime.Format(time.RFC3339Nano),
			ResponseTime: elapsed.Milliseconds(),
		})
		return
	}

	writeJSON(w, http.StatusOK, models.PingResponse{
		PingTime:     pingTime.Format(time.RFC3339Nano),
		ResponseTime: elapsed.Milliseconds(),
	})
}

func hasScheme(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != ""
}
