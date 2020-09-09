package middlewares

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	testHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	logger := Logger(testHandler())
	logger.ServeHTTP(rec, req)

	logged := strings.Split(buf.String(), "|")

	if got, expected := len(logged), 6; got != expected {
		t.Errorf("number of fields expected: %v; got: %v", expected, got)
	}

	// check timestamp
	validTimestamp := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}$`)
	if !validTimestamp.MatchString(strings.TrimSpace(logged[0])) {
		t.Errorf("timestamp has incorrect format")
	}

	// check method
	if got, expected := strings.TrimSpace(logged[1]), "GET"; got != expected {
		t.Errorf("method expected: %v; got: %v", expected, got)
	}

	// check path
	if got, expected := strings.TrimSpace(logged[2]), "/test"; got != expected {
		t.Errorf("path expected: %v; got: %v", expected, got)
	}

	// check status
	if got, expected := strings.TrimSpace(logged[3]), "204"; got != expected {
		t.Errorf("status expected: %v; got: %v", expected, got)
	}

	// check content length
	validContentLength := regexp.MustCompile(`^\d+$`)
	if !validContentLength.MatchString(strings.TrimSpace(logged[4])) {
		t.Errorf("content length has incorrect format")
	}

	// check duration
	validDuration := regexp.MustCompile(`^\d+\.?\d+(µs|ms)$`)
	if !validDuration.MatchString(strings.TrimSpace(logged[5])) {
		t.Errorf("duration has incorrect format")
	}
}
