package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tt := []struct {
		name           string
		statusCode     int
		expectedStatus string
	}{
		{name: "returns 200, with status 'healthy'", statusCode: http.StatusOK, expectedStatus: "healthy"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				t.Fatalf("could not create request; %v", err)
			}

			rec := httptest.NewRecorder()

			HealthHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			// check status code
			if got, expected := res.StatusCode, tc.statusCode; got != expected {
				t.Errorf("status code expected: %v; got: %v", expected, got)
			}

			var decoded healthResponse
			err = json.NewDecoder(res.Body).Decode(&decoded)
			if err != nil {
				t.Fatalf("could not unmarshal json; %v", err)
			}

			if got, expected := decoded.Status, tc.expectedStatus; got != expected {
				t.Errorf("status expected: %v; got: %v", expected, got)
			}
		})
	}
}
