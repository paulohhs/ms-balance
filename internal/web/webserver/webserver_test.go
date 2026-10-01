package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestWebServer_AddHandler(t *testing.T) {
	server := NewWebServer(":0")
	server.AddHandler(http.MethodGet, "/ping/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(chi.URLParam(r, "id")))
	})

	rec := httptest.NewRecorder()
	server.Router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping/42", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "42", rec.Body.String())

	rec = httptest.NewRecorder()
	server.Router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ping/42", nil))
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
