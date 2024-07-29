package handlers

import "net/http"

func getSessionId(r *http.Request) string {
	return r.Header.Get("Session-ID")
}
