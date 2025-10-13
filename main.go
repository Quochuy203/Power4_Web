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
	//http.HandleFunc("/players", players)
	// optionHandler
	http.HandleFunc("/options", optionHandler)
	// 
	//http.HandleFunc("/selectoptions", selectoptions)
	//match
	http.HandleFunc("/match", matchHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	template, erreur := template.ParseFiles("/home.html")
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

var gameConflig GameConflig

type GameConflig struct {
	Mode       string
	Difficulty string
	Rows       int
	Cols       int
}



/*
func options(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		difficulty := r.FormValue("difficulty")
		
		switch difficulty {
		case "easy":
			gameConflig.Difficulty = "easy"
			gameConflig.rows = 6
			gameConflig.cols = 7
		case "normal":
			gameConflig.Difficulty = "normal"
			gameConflig.rows = 6
			gameConflig.cols = 9
		case "hard":
			gameConflig.Difficulty = "hard"
			gameConflig.rows = 7
			gameConflig.cols = 8
		default:
			http.Error(w, "Invalid difficulty", http.StatusBadRequest)
			return
			}
			
			// Sau khi chọn độ khó, chuyển sang trang chơi hoặc khởi tạo game
			http.Redirect(w, r, "/start", http.StatusSeeOther)
			return
			}
			
			// Nếu là GET thì hiển thị trang chọn độ khó
			tmpl, err := template.ParseFiles("template/options.html")
			if err != nil {
				log.Fatal(err)
	}
	tmpl.Execute(w, nil)
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
						
func optionHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
		if mode == "" {
			mode = gameConflig.Mode
		}
		template, erreur := template.ParseFiles("template/options.html")
		if erreur != nil {
		log.Fatal(erreur)
	}
	data := map[string]string{
		"Mode": mode,//gameConflig.Mode,
	}
	template.Execute(w, data)
}
func selectoptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/options", http.StatusSeeOther)
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
		gameConflig.Cols = 7
		gameConflig.Rows = 6
	}
	//log.Printf("Difficulté choisie:", difficulty, gameConflig.Rows, gameConflig.Cols)
	http.Redirect(w, r, "/match", http.StatusSeeOther)
}

// MATCH
type Game struct {
	Board         [][]int //
	CurrentPlayer int
	GameActive    bool
}
func matchHandler(w http.ResponseWriter, r *http.Request) {
	// Khởi tạo bàn cờ theo độ khó đã chọn
	board := make([][]int, gameConflig.Rows)
	for i := range board {
		board[i] = make([]int, gameConflig.Cols)
	}

	game := Game{
		Board:         board,
		CurrentPlayer: 1,
		GameActive:    true,
	}

	// Tạo mảng cột để render dropdown
	columnRange := make([]int, gameConflig.Cols)
	for i := range columnRange {
		columnRange[i] = i
	}

	// Truyền dữ liệu vào template
	data := struct {
		Board         [][]int
		CurrentPlayer int
		ColumnRange   []int
	}{
		Board:         game.Board,
		CurrentPlayer: game.CurrentPlayer,
		ColumnRange:   columnRange,
	}

	tmpl := template.Must(template.ParseFiles("templates/match.html"))
	tmpl.Execute(w, data)
}

