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
	//http.HandleFunc("/start", start)
	// Mode de jeux page select 1vs1 ou 1vsrobot
	http.HandleFunc("/mode", mode)
	// players
	http.HandleFunc("/players", players)
	// optionHandler
	http.HandleFunc("/options", options)

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

/*func start(w http.ResponseWriter, r *http.Request) {
/*if r.Method != http.MethodPost{
	// 		http.Redirect(w, r, "/", http.StatusSeeOther)
	// 		return
// 	}*/
/*	template, erreur := template.ParseFiles("templates/start.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}*/

func mode(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("template/mode.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}
func players(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("tempates/players.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}
func options(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("tempates/options.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}

type Game struct {
	Board         [6][7]int //
	CurrentPlayer int
	GameActive    bool
}
