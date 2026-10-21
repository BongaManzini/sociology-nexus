package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"sociology-nexus/internal/views"
)

type StudyTopic struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Theorists   []string `json:"theorists"`
	Paradigm    string   `json:"paradigm"`
}

var sampleTopics = []StudyTopic{
	{
		ID:          "structural-functionalism",
		Title:       "Structural Functionalism",
		Description: "Views society as a complex system whose parts work together to promote solidarity and stability.",
		Theorists:   []string{"Émile Durkheim", "Talcott Parsons", "Robert K. Merton"},
		Paradigm:    "Macro-Level",
	},
	{
		ID:          "conflict-theory",
		Title:       "Conflict Theory",
		Description: "Focuses on the competition between different social groups for scarce resources like power and wealth.",
		Theorists:   []string{"Karl Marx", "Max Weber", "C. Wright Mills"},
		Paradigm:    "Macro-Level",
	},
	{
		ID:          "symbolic-interactionism",
		Title:       "Symbolic Interactionism",
		Description: "Examines how individuals interact through shared symbols, language, and everyday social meaning.",
		Theorists:   []string{"George Herbert Mead", "Herbert Blumer", "Erving Goffman"},
		Paradigm:    "Micro-Level",
	},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static assets if present
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Main root route renders Templ component
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Student"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = views.Hello(name).Render(r.Context(), w)
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/topics", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(sampleTopics)
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "ok", "service": "sociology-nexus"}`))
		})
	})

	fmt.Printf("🚀 Sociology Nexus server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
