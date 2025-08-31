package server

import (
	"bankstats/internal/domain"
	"bankstats/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"
)

type HTTPServer struct {
	Svc *service.Service
}

func (s *HTTPServer) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		end := time.Now()
		if v := q.Get("end"); v != "" {
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				http.Error(w, "bad end (expect YYYY-MM-DD)", http.StatusBadRequest)
				return
			}
			end = t
		}
		start := end.AddDate(0, 0, -89)
		codesParam := q.Get("codes")
		if codesParam == "" {
			codesParam = q.Get("code")
		}
		codeSet := make(map[string]bool)
		if codesParam != "" {
			for _, p := range strings.Split(codesParam, ",") {
				if c := strings.ToUpper(strings.TrimSpace(p)); c != "" {
					codeSet[c] = true
				}
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
		defer cancel()

		points, err := s.Svc.CollectRange(ctx, start, end, codeSet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		stats := s.Svc.CalcStats(points, 90)

		type response struct {
			Start string       `json:"start"`
			End   string       `json:"end"`
			Codes []string     `json:"codes,omitempty"`
			Stats domain.Stats `json:"stats"`
		}
		var codes []string
		for c := range codeSet {
			codes = append(codes, c)
		}
		sort.Strings(codes)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(response{
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
			Codes: codes,
			Stats: stats,
		})
	})

	return mux
}

func (s *HTTPServer) Run(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 180 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv.ListenAndServe()
}
