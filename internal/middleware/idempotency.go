package middleware

import (
	"bytes"
	"database/sql"
	"net/http"
)

func Idempotency(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get Key from Header
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Check if key exists
			var storedBody string
			var storedStatus int
			err := db.QueryRow("SELECT response_body, status_code FROM idempotency_keys WHERE key = $1", key).Scan(&storedBody, &storedStatus)

			if err == nil {
				// Key exists! Return cached response immediately
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Idempotency-Hit", "true")
				w.WriteHeader(storedStatus)
				w.Write([]byte(storedBody))
				return
			}

			// 3. New Key - intercept response
			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
			}

			next.ServeHTTP(recorder, r)

			// 4. Save response to DB
			_, _ = db.Exec(
				"INSERT INTO idempotency_keys (key, response_body, status_code) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
				key,
				recorder.body.String(),
				recorder.statusCode,
			)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
