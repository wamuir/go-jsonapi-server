package handle

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCopyHeader(t *testing.T) {

	h := make(http.Header)
	h.Add("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	h.Add("Content-Type", "text/html; charset=utf-8")
	h.Add("Expires", "Thu, 01 Jan 1970 00:00:00 UTC")
	h.Add("Pragma", "no-cache")
	h.Add("X-Accel-Expires", "0")
	h.Add("X-Content-Type-Options", "nosniff")
	h.Add("X-Ratelimit-Limit", "100")
	h.Add("X-Ratelimit-Remaining", "100")
	h.Add("X-Ratelimit-Reset", "1626202320")
	h.Add("Date", "Tue, 13 Jul 2021 18:51:45 GMT")
	h.Add("Content-Length", "1057")

	d := make(http.Header)
	copyHeader(d, h)

	if ok := reflect.DeepEqual(d, h); !ok {
		t.Fatal("Headers key-value pairs not equal")
	}
}

func TestValidateMIME(t *testing.T) {

	if err := ValidateMIME(""); err == nil {
		t.Fatal("Unexpected nil error")
	} else if err.Status != fmt.Sprintf("%d", http.StatusUnsupportedMediaType) {
		t.Fatalf("Unexpected %v", err.Status)
	}

	if err := ValidateMIME("application/json"); err == nil {
		t.Fatal("Unexpected nil error")
	} else if err.Status != fmt.Sprintf("%d", http.StatusUnsupportedMediaType) {
		t.Fatalf("Unexpected %v", err.Status)
	}

	if err := ValidateMIME("application/vnd.api+json"); err != nil {
		t.Fatalf("Unexpected %v", err)
	}

}
