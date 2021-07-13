package handle

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle405(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	new(Environment).Handle405(w, r)
	res := w.Result()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf(
			"res.StatusCode = %v, want %v",
			res.StatusCode,
			http.StatusMethodNotAllowed,
		)
	}
}
