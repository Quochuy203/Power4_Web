package main

import (
	"power4_web/web"
	"net/http"
)

func main() {
	http.HandleFunc("/", web.IndexHandler)
	http.HandleFunc("/welcome", web.WelcomeHandler)
	http.HandleFunc("/game", web.GameHandler)
	http.HandleFunc("/play", web.PlayHandler)
	http.HandleFunc("/rematch", web.RematchHandler)
	http.HandleFunc("/mockup-victoire", web.MockupVictoireHandler)
	http.HandleFunc("/mockup-nul", web.MockupNulHandler)
	http.HandleFunc("/mockup-gravite-inverse", web.MockupGraviteInverseHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
