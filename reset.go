package main

import "net/http"

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	}

	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(r.Context())

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset"))
}
