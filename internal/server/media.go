package server

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) MediaHandler(w http.ResponseWriter, r *http.Request) {
	objectName := strings.TrimPrefix(r.URL.Path, "/media/")
	if objectName == "" {
		http.NotFound(w, r)
		return
	}

	reader, size, contentType, err := s.storage.GetImage(r.Context(), objectName)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", contentType)
	if size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
	_, _ = io.Copy(w, reader)
}
