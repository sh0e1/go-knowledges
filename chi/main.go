package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer func() {
		_ = logger.Sync()
	}()

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, Logger(logger))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	r.Get("/delay", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(w, "Hello Delay")
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func Logger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				l.Info(fmt.Sprintf("finished http request with code %d", ww.Status()),
					zap.Int64("ts", t1.UnixNano()),
					zap.Time("http.start_time", t1),
					zap.String("http.proto", r.Proto),
					zap.String("http.method", r.Method),
					zap.String("http.path", r.URL.Path),
					zap.String("access.client_ip", r.RemoteAddr),
					zap.String("access.user_agent", r.UserAgent()),
					zap.Int("http.code", ww.Status()),
					zap.Int("http.size", ww.BytesWritten()),
					zap.Float32("http.time_ms", float32(time.Since(t1).Nanoseconds()/1000)/1000),
					zap.String("req.id", middleware.GetReqID(r.Context())),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
