package middleware

import (
	"bytes"
	"net/http"
)

type CustomResponseRecorder struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) *CustomResponseRecorder {
	return &CustomResponseRecorder{
		ResponseWriter: w,
		Body:           new(bytes.Buffer),
	}
}

func (r *CustomResponseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *CustomResponseRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}
