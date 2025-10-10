package main

import (
	"html/template"
	"log"
	"net/http"
)

type GameConflig struct {
	Mode       string
	Difficulty string
	rows       int
	cols       int
}

var gameConflig GameConflig

type Game struct {
	Board         [][]int //
	CurrentPlayer int
	GameActive    bool
}

func main() {
	//home routeur
	http.HandleFunc("/", home)
	// play
	//http.HandleFunc("/start", start)
	// Mode de jeux page select 1vs1 ou 1vsrobot
	http.HandleFunc("/mode", mode)
	// players
	//http.HandleFunc("/players", players)
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

func mode(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("template/mode.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}

/*
	func players(w http.ResponseWriter, r *http.Request) {
		template, erreur := template.ParseFiles("tempates/players.html")
		if erreur != nil {
			log.Fatal(erreur)
		}
		template.Execute(w, nil)
	}
*/
/*func options(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("tempates/options.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	template.Execute(w, nil)
}*/
/*
func options(w http.ResponseWriter, r *http.Request) {
	if r.Method != "post" {
		http.Redirect(w, r, "/options", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	data := OptionsData{
		GameMode:    r.FormValue("mode"),
		Player1name: r.FormValue("player1"),
		Player2name: r.FormValue("player2"),
	}
	if data.Player1name == "" && data.GameMode == "1v1" {
		data.Player1name = "Player 1"
		data.Player2name = "Player 2"
	}
	template, erreur := template.ParseFiles("templates/options.html")
	if erreur != nil {
		log.Println("", erreur)
		http.Error(w, " Erreur server", http.StatusInternalServerError)
		return
	}
	template.Execute(w, data)
}

type OptionsData struct {
	GameMode    string
	Player1name string
	Player2name string
}
*/

func difficultyHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = gameConflig.Mode
	}
	template, erreur := template.ParseFiles("template/difficulty.html")
	if erreur != nil {
		log.Fatal(erreur)
	}
	data := [string]string{
		"Mode": gameConflig.Mode,
	}
	template.Execute(w, data)
}
func selectDifficultyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/difficulty", http.StatusSeeOther)
		return
	}

	difficulty := r.FormValue("difficulty")

	// Configurer la grille selon la difficulté
	gameConflig.Difficulty = difficulty
	switch difficulty {
	case "easy":
		gameConflig.Rows = 6
		gameConflig.Cols = 7
	case "normal":
		gameConflig.Rows = 6
		gameConflig.Cols = 9
	case "hard":
		gameConflig.Rows = 7
		gameConflig.Cols = 8
	default:
		gameConflig.Rows = 6
		gameConflig.Cols = 7
	}
	log.Printf("Difficulté choisie:", difficulty, gameConflig.Rows, gameConflig.Cols)
	http.Redirect(w, r, "/players", http.StatusSeeOther)
}
