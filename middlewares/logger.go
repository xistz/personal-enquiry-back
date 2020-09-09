package middlewares

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"text/template"
	"time"
)

type logContent struct {
	StartTime     string
	Method        string
	Path          string
	Status        int
	ContentLength int
	Duration      time.Duration
}

// Logger middleware to log requests
func Logger(next http.Handler) http.Handler {
	logFormat := "{{.StartTime}} | {{.Method}} | {{.Path}} | {{.Status}} | {{.ContentLength}} | {{.Duration}}"
	t := template.Must(template.New("logger").Parse(logFormat))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// :method :url :status :res[content-length] - :response-time ms
		start := time.Now()

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())

		duration := time.Since(start)

		buf := new(bytes.Buffer)
		t.Execute(
			buf,
			logContent{
				StartTime:     start.Format(time.RFC3339),
				Method:        r.Method,
				Path:          r.URL.Path,
				Status:        rec.Code,
				ContentLength: rec.Body.Len(),
				Duration:      duration,
			},
		)
		log.SetFlags(0)
		log.Println(buf.String())
	})
}
