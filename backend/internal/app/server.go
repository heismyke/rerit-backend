package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rerit/backend/internal/store"
)

type Server struct {
	store     *store.Store
	staticDir string
	logger    *slog.Logger
	mux       *http.ServeMux
	origins   []string
}

func New(st *store.Store, staticDir string, logger *slog.Logger) *Server {
	s := &Server{
		store:     st,
		staticDir: staticDir,
		logger:    logger,
		mux:       http.NewServeMux(),
		origins:   allowedOrigins(),
	}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.security(s.cors(s.logging(s.mux)))
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /api/bootstrap", s.bootstrap)
	s.mux.HandleFunc("POST /api/auth/login", s.login)

	s.mux.HandleFunc("GET /api/properties", s.listProperties)
	s.mux.HandleFunc("POST /api/properties", s.createProperty)
	s.mux.HandleFunc("PUT /api/properties/{id}", s.updateProperty)
	s.mux.HandleFunc("DELETE /api/properties/{id}", s.deleteProperty)

	s.mux.HandleFunc("GET /api/audit-cases", s.listAuditCases)
	s.mux.HandleFunc("POST /api/audit-cases", s.createAuditCase)
	s.mux.HandleFunc("PUT /api/audit-cases/{id}", s.updateAuditCase)
	s.mux.HandleFunc("DELETE /api/audit-cases/{id}", s.deleteAuditCase)
	s.mux.HandleFunc("POST /api/audit-cases/{id}/result", s.sendAuditResult)

	s.mux.HandleFunc("GET /api/flagged-cases", s.listFlaggedCases)
	s.mux.HandleFunc("PUT /api/flagged-cases/{id}", s.updateFlaggedCase)
	s.mux.HandleFunc("POST /api/flagged-cases/{id}/result", s.sendFlaggedResult)

	s.mux.HandleFunc("GET /api/successful-filings", s.listSuccessfulFilings)
	s.mux.HandleFunc("PUT /api/successful-filings/{id}", s.updateSuccessfulFiling)

	s.mux.HandleFunc("GET /api/notices", s.listNotices)
	s.mux.HandleFunc("POST /api/notices/{id}/respond", s.respondNotice)

	s.mux.HandleFunc("GET /api/payments", s.listPayments)
	s.mux.HandleFunc("POST /api/payments", s.createPayment)

	s.mux.HandleFunc("/", s.static)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) bootstrap(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.Snapshot())
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	user, err := s.store.Authenticate(req.Email, req.Role)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": randomToken(),
	})
}

func (s *Server) listProperties(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.ListProperties())
}

func (s *Server) createProperty(w http.ResponseWriter, r *http.Request) {
	var p store.Property
	if !decodeJSON(w, r, &p) {
		return
	}
	created, err := s.store.CreateProperty(p)
	writeResult(w, created, err)
}

func (s *Server) updateProperty(w http.ResponseWriter, r *http.Request) {
	var p store.Property
	if !decodeJSON(w, r, &p) {
		return
	}
	updated, err := s.store.UpdateProperty(r.PathValue("id"), p)
	writeResult(w, updated, err)
}

func (s *Server) deleteProperty(w http.ResponseWriter, r *http.Request) {
	writeDelete(w, s.store.DeleteProperty(r.PathValue("id")))
}

func (s *Server) listAuditCases(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.ListAuditCases())
}

func (s *Server) createAuditCase(w http.ResponseWriter, r *http.Request) {
	var c store.AuditCase
	if !decodeJSON(w, r, &c) {
		return
	}
	created, err := s.store.CreateAuditCase(c)
	writeResult(w, created, err)
}

func (s *Server) updateAuditCase(w http.ResponseWriter, r *http.Request) {
	var c store.AuditCase
	if !decodeJSON(w, r, &c) {
		return
	}
	updated, err := s.store.UpdateAuditCase(r.PathValue("id"), c)
	writeResult(w, updated, err)
}

func (s *Server) deleteAuditCase(w http.ResponseWriter, r *http.Request) {
	writeDelete(w, s.store.DeleteAuditCase(r.PathValue("id")))
}

func (s *Server) sendAuditResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	db, err := s.store.SendAuditResult(r.PathValue("id"), req.Status, req.Notes)
	writeResult(w, db, err)
}

func (s *Server) listFlaggedCases(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.ListFlaggedCases())
}

func (s *Server) updateFlaggedCase(w http.ResponseWriter, r *http.Request) {
	var c store.FlaggedCase
	if !decodeJSON(w, r, &c) {
		return
	}
	updated, err := s.store.UpdateFlaggedCase(r.PathValue("id"), c)
	writeResult(w, updated, err)
}

func (s *Server) sendFlaggedResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	updated, err := s.store.SendFlaggedResult(r.PathValue("id"), req.Status, req.Notes)
	writeResult(w, updated, err)
}

func (s *Server) listSuccessfulFilings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.ListSuccessfulFilings())
}

func (s *Server) updateSuccessfulFiling(w http.ResponseWriter, r *http.Request) {
	var f store.SuccessfulFiling
	if !decodeJSON(w, r, &f) {
		return
	}
	updated, err := s.store.UpdateSuccessfulFiling(r.PathValue("id"), f)
	writeResult(w, updated, err)
}

func (s *Server) listNotices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.ListNotices())
}

func (s *Server) respondNotice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Response string `json:"response"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	updated, err := s.store.RespondNotice(r.PathValue("id"), req.Response)
	writeResult(w, updated, err)
}

func (s *Server) listPayments(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.ListPayments())
}

func (s *Server) createPayment(w http.ResponseWriter, r *http.Request) {
	var p store.Payment
	if !decodeJSON(w, r, &p) {
		return
	}
	created, err := s.store.CreatePayment(p)
	writeResult(w, created, err)
}

func (s *Server) static(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	path := filepath.Join(s.staticDir, filepath.Clean(r.URL.Path))
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		setContentType(w, path)
		http.ServeFile(w, r, path)
		return
	}
	index := filepath.Join(s.staticDir, "index.html")
	http.ServeFile(w, r, index)
}

func (s *Server) security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && s.originAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "600")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) originAllowed(origin string) bool {
	for _, allowed := range s.origins {
		if allowed == "*" || allowed == origin {
			return true
		}
		if strings.HasPrefix(allowed, "https://*.") {
			suffix := strings.TrimPrefix(allowed, "https://*")
			if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, suffix) {
				return true
			}
		}
	}
	return false
}

func (s *Server) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &logResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lrw, r)
		s.logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", lrw.status, "duration_ms", time.Since(start).Milliseconds())
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer r.Body.Close()
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func writeResult(w http.ResponseWriter, value any, err error) {
	if errors.Is(err, os.ErrNotExist) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func writeDelete(w http.ResponseWriter, err error) {
	if errors.Is(err, os.ErrNotExist) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

func setContentType(w http.ResponseWriter, path string) {
	if ct := mime.TypeByExtension(filepath.Ext(path)); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
}

func allowedOrigins() []string {
	value := os.Getenv("RERIT_ALLOWED_ORIGINS")
	if value == "" {
		value = "http://localhost:5173,http://localhost:5174,https://*.vercel.app"
	}
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			origins = append(origins, part)
		}
	}
	return origins
}
