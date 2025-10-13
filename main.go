package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	// Router de xu lu nuoc di cua nguoi choi
	http.HandleFunc("/play",playMove)
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

var gameConflig GameConflig

type GameConflig struct {
	Mode       string
	Difficulty string
	Rows       int
	Cols       int
}
var currentGame Game


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

	currentGame := Game{
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
	}{ // sua game.Board thanh currentGame.Board
		Board:         currentGame.Board,
		CurrentPlayer: currentGame.CurrentPlayer,
		ColumnRange:   columnRange,
	}

	tmpl := template.Must(template.ParseFiles("template/match.html"))
	tmpl.Execute(w, data)
}
// Trien khai nuoc di cua Robot
// Helper functions (Nên đặt trong một file riêng như game_logic.go)
func checkWin(b [][]int, piece int) bool { /* ... logic kiểm tra thắng ... */
	 return false }
func getNextOpenRow(b [][]int, col int) int { /* ... */ 
	return -1 }
// ...

func playMove(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" || !currentGame.GameActive {
        http.Redirect(w, r, "/match", http.StatusSeeOther)
        return
    }

    // 1. Xử lý nước đi của Người chơi (Client)
    r.ParseForm()
    colStr := r.FormValue("col") // Lấy cột từ form
    col, err := strconv.Atoi(colStr)
    if err != nil {
        http.Redirect(w,r,"/match",http.StatusSeeOther)
		return
    }
    
    row := getNextOpenRow(currentGame.Board, col)

    if row != -1 {
        // Thả quân người chơi
        currentGame.Board[row][col] = currentGame.CurrentPlayer
        
        // 2. Kiểm tra thắng
        if checkWin(currentGame.Board, currentGame.CurrentPlayer) {
            currentGame.GameActive = false
            // Có thể thêm logic lưu kết quả và hiển thị thông báo
        }
        
        // 3. Chuyển lượt
        currentGame.CurrentPlayer = 3 - currentGame.CurrentPlayer // Chuyển từ 1 sang 2 hoặc ngược lại
    }

    // 4. LƯỢT CỦA ROBOT (nếu GameActive và là chế độ 1vsRobot)
    if currentGame.GameActive && currentGame.CurrentPlayer == 2 && gameConflig.Mode == "robot" {
        // Đây là nơi bạn gọi AI
        
        // bestCol := getBestMoveAI(currentGame.Board, gameConflig.Difficulty)
        bestCol := 3 // Tạm thời chọn cột 3
        
        aiRow := getNextOpenRow(currentGame.Board, bestCol)
        if aiRow != -1 {
            currentGame.Board[aiRow][bestCol] = 2 // Quân của Robot là 2
            
            // 5. Kiểm tra thắng của Robot
            if checkWin(currentGame.Board, 2) {
                currentGame.GameActive = false
                // ...
            }
            
            // 6. Chuyển lượt lại cho Người chơi
            currentGame.CurrentPlayer = 1
        }
    }
    
    // Chuyển hướng về trang /match để hiển thị trạng thái mới
    http.Redirect(w, r, "/match", http.StatusSeeOther)
}
