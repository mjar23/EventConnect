package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func StaticFileRoutes(r *mux.Router) {
	// Serve static files from the "css" directory
	cssDir := http.Dir("css")
	cssHandler := http.StripPrefix("/css/", http.FileServer(cssDir))
	r.PathPrefix("/css/").Handler(cssHandler)

	// Serve JavaScript files
	serveJSFile(r, "/events.js")
	serveJSFile(r, "/event-details.js")
	serveJSFile(r, "/profile.js")
	serveJSFile(r, "/login.js")
	serveJSFile(r, "/event-comments.js")
}

func serveJSFile(r *mux.Router, path string) {
    r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "js"+path)
    })
}