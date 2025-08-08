package middleware

import "net/http"

type WrapperWrite struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrapperWrite) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}
