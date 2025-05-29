package middlewares

import "net/http"

type wrappedResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{ResponseWriter: w}
}

func (w wrappedResponseWriter) Status() int {
	return w.status
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	w.status = code
	w.ResponseWriter.WriteHeader(code)
	w.wroteHeader = true
	return
}
