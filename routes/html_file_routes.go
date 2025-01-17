package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HTMLFileRoutes(r *mux.Router) {
	// Serve HTML files
	serveHTMLFile(r, "/", "html/events.html")
	serveHTMLFile(r, "/login.html", "html/login.html")
	serveHTMLFile(r, "/profile.html", "html/profile.html")
	serveHTMLFile(r, "/event-details.html", "html/event-details.html")
	serveHTMLFile(r, "/event-comments.html", "html/event-comments.html")
	serveHTMLFile(r, "/other-user-profile.html", "html/other-user-profile.html")
	serveHTMLFile(r, "/recommendations.html", "html/recommendations.html")
}

func serveHTMLFile(r *mux.Router, path, filename string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}