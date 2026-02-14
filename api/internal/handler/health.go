package handler

import "net/http"

// Health returns a simple health check handler.
func Health(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
