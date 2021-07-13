package handle

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle404(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	new(Environment).Handle404(w, r)
	res := w.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf(
			"res.StatusCode = %v, want %v",
			res.StatusCode,
			http.StatusNotFound,
		)
	}
}
