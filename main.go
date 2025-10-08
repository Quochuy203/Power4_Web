package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	//home routeur
	http.HandleFunc("/", home)
	// play
	http.HandleFunc("/start", start)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("home.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}

func start(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != http.MethodPost{
		// 		http.Redirect(w, r, "/", http.StatusSeeOther)
		// 		return
	// 	}*/
	template, erreur := template.ParseFiles("templates/start.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}
