package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"math/rand"
)

func main() {
	http.HandleFunc("/", home)
	//http.HandleFunc("/start", start)
	// Mode de jeux page select 1vs1 ou 1vsrobot
	http.HandleFunc("/mode", mode)
	// players
	//http.HandleFunc("/players", players)
	// optionHandler
	http.HandleFunc("/options", optionHandler)
	//http.HandleFunc("/selectoptions", selectoptions)
	//match
	http.HandleFunc("/match", matchHandler)
	// Router de xu lu nuoc di cua nguoi choi
	http.HandleFunc("/play",playMove)
	// Reset game
	http.HandleFunc("/reset",resetGame)
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
		"Mode": gameConflig.Mode,//mode,//gameConflig.Mode,
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
	Winner 		  int
	Message		  string
}
func matchHandler(w http.ResponseWriter, r *http.Request) {
	// Khởi tạo bàn cờ theo độ khó đã chọn
	board := make([][]int, gameConflig.Rows)
	for i := range board {
		board[i] = make([]int, gameConflig.Cols)
	}
// vua sua thanh currentGame = Game { ( truoc do là :=)
	currentGame = Game{
		Board:         board,
		CurrentPlayer: 1,
		GameActive:    true,
		Winner:			0,
		Message:		"",
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
		GameActive	  bool
		Winner		  int
		Message		  string
		Mode		  string
	}{ // sua game.Board thanh currentGame.Board
		Board:         currentGame.Board,
		CurrentPlayer: currentGame.CurrentPlayer,
		ColumnRange:   columnRange,
		GameActive:    currentGame.GameActive,
		Winner:        currentGame.Winner,
		Message:       currentGame.Message,
		Mode:          gameConflig.Mode,
	}

	tmpl := template.Must(template.ParseFiles("template/match.html"))
	tmpl.Execute(w, data)
}
// Trien khai nuoc di cua Robot
// Helper functions (Nên đặt trong một file riêng như game_logic.go)
/*func checkWin(b [][]int, piece int) bool { 
	rows := len(board)
	cols := len(board[0])

	// Kiểm tra ngang
	count := 0
	for c := 0; c < cols; c++ {
		if board[row][c] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
	}

	// Kiểm tra dọc
	count = 0
	for r := 0; r < rows; r++ {
		if board[r][col] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
	}

	// Kiểm tra chéo chính (\)
	count = 0
	startRow := row - min(row, col)
	startCol := col - min(row, col)
	for startRow < rows && startCol < cols {
		if board[startRow][startCol] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
		startRow++
		startCol++
	}

	// Kiểm tra chéo phụ (/)
	count = 0
	startRow = row + min(rows-1-row, col)
	startCol = col - min(rows-1-row, col)
	for startRow >= 0 && startCol < cols {
		if board[startRow][startCol] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
		startRow--
		startCol++
	}
	return false 
	}
*/
/*
func getNextOpenRow(b [][]int, col int) int {  
	return -1 
}
*/
// ...
/*
func playMove(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" || !currentGame.GameActive {
        http.Redirect(w, r, "/match", http.StatusSeeOther)
        return
    }

    // 1. Xử lý nước đi của Người chơi (Client)
    r.ParseForm()
    colStr := r.FormValue("col") // Lấy cột từ form
    col, err := strconv.Atoi(colStr)

    if err != nil || col < 0 || col > gameConflig {
        http.Redirect(w,r,"/match",http.StatusSeeOther)
		return
    }
    
    row := getNextOpenRow(currentGame.Board, col)

    if row != -1 {
        // Thả quân người chơi
       // currentGame.Board[row][col] = currentGame.CurrentPlayer
	    http.Redirect(w,r,"/match", http.StatusSeeOther)
		return
        
        // 2. Kiểm tra thắng
        if checkWin(currentGame.Board,row, col, currentGame.CurrentPlayer) {
            currentGame.GameActive = false
			currentGame.Winner = currentGame.CurrentPlayer
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
} */



// chu y den nuoc co hoà , neu ca 2  nguoi choi deu hao nthanh nuoc di cua ho nhung khong co ai thang thi ca 2 deu hoa 

func playMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || !currentGame.GameActive {
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	colStr := r.FormValue("col")
	col, err := strconv.Atoi(colStr)
	
	if err != nil || col < 0 || col >= gameConflig.Cols {
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}

	// 1. Xử lý nước đi của người chơi
	row := getNextOpenRow(currentGame.Board, col)
	
	if row == -1 {
		// Cột đầy
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}

	// Thả quân người chơi
	currentGame.Board[row][col] = currentGame.CurrentPlayer

	// 2. Kiểm tra thắng
	if checkWin(currentGame.Board, row, col, currentGame.CurrentPlayer) {
		currentGame.GameActive = false
		currentGame.Winner = currentGame.CurrentPlayer
		if currentGame.CurrentPlayer == 1 {
			currentGame.Message = "Người chơi 1 (Đỏ) thắng!"
		} else {
			currentGame.Message = "Người chơi 2 (Vàng) thắng!"
		}
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}

	// Kiểm tra hòa
	if isBoardFull(currentGame.Board) {
		currentGame.GameActive = false
		currentGame.Message = "Hòa!"
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}

	// 3. Chuyển lượt
	currentGame.CurrentPlayer = 3 - currentGame.CurrentPlayer

	// 4. Nếu là chế độ robot và đến lượt robot
	if currentGame.GameActive && currentGame.CurrentPlayer == 2 && gameConflig.Mode == "robot" {
		// AI chọn cột
		bestCol := getBestMoveAI(currentGame.Board, gameConflig.Difficulty)
		
		aiRow := getNextOpenRow(currentGame.Board, bestCol)
		if aiRow != -1 {
			currentGame.Board[aiRow][bestCol] = 2

			// 5. Kiểm tra thắng của robot
			if checkWin(currentGame.Board, aiRow, bestCol, 2) {
				currentGame.GameActive = false
				currentGame.Winner = 2
				currentGame.Message = "Robot thắng!"
				http.Redirect(w, r, "/match", http.StatusSeeOther)
				return
			}

			// Kiểm tra hòa
			if isBoardFull(currentGame.Board) {
				currentGame.GameActive = false
				currentGame.Message = "Hòa!"
				http.Redirect(w, r, "/match", http.StatusSeeOther)
				return
			}

			// 6. Chuyển lượt lại cho người chơi
			currentGame.CurrentPlayer = 1
		}
	}

	http.Redirect(w, r, "/match", http.StatusSeeOther)
}
func resetGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/match", http.StatusSeeOther)
}
// Tìm hàng trống tiếp theo trong cột
func getNextOpenRow(board [][]int, col int) int {
	for row := len(board) - 1; row >= 0; row-- {
		if board[row][col] == 0 {
			return row
		}
	}
	return -1
}

// Kiểm tra thắng
func checkWin(board [][]int, row, col, player int) bool {
	rows := len(board)
	cols := len(board[0])

	// Kiểm tra ngang
	count := 0
	for c := 0; c < cols; c++ {
		if board[row][c] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
	}

	// Kiểm tra dọc
	count = 0
	for r := 0; r < rows; r++ {
		if board[r][col] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
	}

	// Kiểm tra chéo chính (\)
	count = 0
	startRow := row - min(row, col)
	startCol := col - min(row, col)
	for startRow < rows && startCol < cols {
		if board[startRow][startCol] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
		startRow++
		startCol++
	}

	// Kiểm tra chéo phụ (/)
	count = 0
	startRow = row + min(rows-1-row, col)
	startCol = col - min(rows-1-row, col)
	for startRow >= 0 && startCol < cols {
		if board[startRow][startCol] == player {
			count++
			if count >= 4 {
				return true
			}
		} else {
			count = 0
		}
		startRow--
		startCol++
	}

	return false
}

// Kiểm tra board đầy
func isBoardFull(board [][]int) bool {
	for _, cell := range board[0] {
		if cell == 0 {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
// Choi cung voi IA
func getBestMoveAI(board [][]int, difficulty string) int {
	cols := len(board[0])
	validCols := []int{}

	// Lấy danh sách cột còn trống
	for col := 0; col < cols; col++ {
		if getNextOpenRow(board, col) != -1 {
			validCols = append(validCols, col)
		}
	}

	if len(validCols) == 0 {
		return 0
	}

	switch difficulty {
	case "easy":
		// Random move
		return validCols[rand.Intn(len(validCols))]
		
	case "normal":
		// Kiểm tra nước thắng hoặc chặn đối thủ
		// 1. Thử thắng
		for _, col := range validCols {
			row := getNextOpenRow(board, col)
			board[row][col] = 2
			if checkWin(board, row, col, 2) {
				board[row][col] = 0
				return col
			}
			board[row][col] = 0
		}

		// 2. Chặn đối thủ
		for _, col := range validCols {
			row := getNextOpenRow(board, col)
			board[row][col] = 1
			if checkWin(board, row, col, 1) {
				board[row][col] = 0
				return col
			}
			board[row][col] = 0
		}

		// 3. Random
		return validCols[rand.Intn(len(validCols))]
		
	case "hard":
		// Minimax với độ sâu giới hạn
		bestScore := -999999
		bestCol := validCols[0]

		for _, col := range validCols {
			row := getNextOpenRow(board, col)
			board[row][col] = 2
			
			score := minimax(board, 3, false, -999999, 999999)
			
			board[row][col] = 0

			if score > bestScore {
				bestScore = score
				bestCol = col
			}
		}
		return bestCol
		
	default:
		return validCols[rand.Intn(len(validCols))]
	}
}
func minimax(board [][]int, depth int, isMaximizing bool, alpha, beta int) int {
	cols := len(board[0])

	// Kiểm tra terminal states
	for col := 0; col < cols; col++ {
		row := getNextOpenRow(board, col)
		if row != -1 {
			// Kiểm tra AI thắng
			board[row][col] = 2
			if checkWin(board, row, col, 2) {
				board[row][col] = 0
				return 100000
			}
			board[row][col] = 0

			// Kiểm tra người chơi thắng
			board[row][col] = 1
			if checkWin(board, row, col, 1) {
				board[row][col] = 0
				return -100000
			}
			board[row][col] = 0
		}
	}

	// Kiểm tra độ sâu hoặc board đầy
	if depth == 0 || isBoardFull(board) {
		return evaluateBoard(board)
	}

	validCols := []int{}
	for col := 0; col < cols; col++ {
		if getNextOpenRow(board, col) != -1 {
			validCols = append(validCols, col)
		}
	}

	if isMaximizing {
		maxScore := -999999
		for _, col := range validCols {
			row := getNextOpenRow(board, col)
			board[row][col] = 2
			
			score := minimax(board, depth-1, false, alpha, beta)
			
			board[row][col] = 0

			maxScore = max(maxScore, score)
			alpha = max(alpha, score)
			if beta <= alpha {
				break
			}
		}
		return maxScore
	} else {
		minScore := 999999
		for _, col := range validCols {
			row := getNextOpenRow(board, col)
			board[row][col] = 1
			
			score := minimax(board, depth-1, true, alpha, beta)
			
			board[row][col] = 0

			minScore = min(minScore, score)
			beta = min(beta, score)
			if beta <= alpha {
				break
			}
		}
		return minScore
	}
}
// Đánh giá board
func evaluateBoard(board [][]int) int {
	score := 0
	rows := len(board)
	cols := len(board[0])

	// Ưu tiên cột giữa
	centerCol := cols / 2
	for row := 0; row < rows; row++ {
		if board[row][centerCol] == 2 {
			score += 3
		} else if board[row][centerCol] == 1 {
			score -= 3
		}
	}

	// Đánh giá các chuỗi
	score += evaluateSequences(board, 2) - evaluateSequences(board, 1)

	return score
}
func evaluateSequences(board [][]int, player int) int {
	score := 0
	rows := len(board)
	cols := len(board[0])

	// Kiểm tra ngang
	for r := 0; r < rows; r++ {
		for c := 0; c <= cols-4; c++ {
			window := []int{board[r][c], board[r][c+1], board[r][c+2], board[r][c+3]}
			score += scoreWindow(window, player)
		}
	}

	// Kiểm tra dọc
	for c := 0; c < cols; c++ {
		for r := 0; r <= rows-4; r++ {
			window := []int{board[r][c], board[r+1][c], board[r+2][c], board[r+3][c]}
			score += scoreWindow(window, player)
		}
	}

	// Kiểm tra chéo chính
	for r := 0; r <= rows-4; r++ {
		for c := 0; c <= cols-4; c++ {
			window := []int{board[r][c], board[r+1][c+1], board[r+2][c+2], board[r+3][c+3]}
			score += scoreWindow(window, player)
		}
	}

	// Kiểm tra chéo phụ
	for r := 3; r < rows; r++ {
		for c := 0; c <= cols-4; c++ {
			window := []int{board[r][c], board[r-1][c+1], board[r-2][c+2], board[r-3][c+3]}
			score += scoreWindow(window, player)
		}
	}

	return score
}

func scoreWindow(window []int, player int) int {
	score := 0
	opponent := 3 - player

	playerCount := 0
	opponentCount := 0
	emptyCount := 0

	for _, cell := range window {
		if cell == player {
			playerCount++
		} else if cell == opponent {
			opponentCount++
		} else {
			emptyCount++
		}
	}

	if playerCount == 4 {
		score += 100
	} else if playerCount == 3 && emptyCount == 1 {
		score += 5
	} else if playerCount == 2 && emptyCount == 2 {
		score += 2
	}

	if opponentCount == 3 && emptyCount == 1 {
		score -= 4
	}

	return score
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}



//func play(w, http.ResponseWriter, r *Request)











































































/*
func optionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		gameConflig.Difficulty = r.FormValue("difficulty")
		switch gameConflig.Difficulty {
		case "easy":
			gameConflig.Rows, gameConflig.Cols = 6, 7
		case "normal":
			gameConflig.Rows, gameConflig.Cols = 6, 9
		case "hard":
			gameConflig.Rows, gameConflig.Cols = 7, 8
		default:
			gameConflig.Rows, gameConflig.Cols = 6, 7
		}
		// create board
		board := make([][]int, gameConflig.Rows)
		for i := range board {
			board[i] = make([]int, gameConflig.Cols)
		}
		currentGame = Game{
			Board: board, 
			CurrentPlayer: 1,
			GameActive: true
		}
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/options.html"))
	tmpl.Execute(w, gameConflig)
}

// MATCH
func matchHandler(w http.ResponseWriter, r *http.Request) {
	columnRange := make([]int, gameConflig.Cols)
	for i := range columnRange {
		columnRange[i] = i
	}
	data := struct {
		Board         [][]int
		CurrentPlayer int
		ColumnRange   []int
	}{
		currentGame.Board,
		currentGame.CurrentPlayer,
		columnRange,
	}
	tmpl := template.Must(template.ParseFiles("templates/match.html"))
	tmpl.Execute(w, data)
}

// PLAY
func playMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !currentGame.GameActive {
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}

	colStr := r.FormValue("col")
	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Redirect(w, r, "/match", http.StatusSeeOther)
		return
	}
	row := getNextOpenRow(currentGame.Board, col)
	if row != -1 {
		currentGame.Board[row][col] = currentGame.CurrentPlayer
		if checkWin(currentGame.Board, currentGame.CurrentPlayer) {
			currentGame.GameActive = false
		}
		currentGame.CurrentPlayer = 3 - currentGame.CurrentPlayer
	}
	http.Redirect(w, r, "/match", http.StatusSeeOther)
}

func getNextOpenRow(b [][]int, col int) int {
	for r := len(b) - 1; r >= 0; r-- {
		if b[r][col] == 0 {
			return r
		}
	}
	return -1
}

func checkWin(b [][]int, p int) bool {
	rows := len(b)
	cols := len(b[0])
	for r := 0; r < rows; r++ {
		for c := 0; c < cols-3; c++ {
			if b[r][c] == p && b[r][c+1] == p && b[r][c+2] == p && b[r][c+3] == p {
				return true
			}
		}
	}
	for c := 0; c < cols; c++ {
		for r := 0; r < rows-3; r++ {
			if b[r][c] == p && b[r+1][c] == p && b[r+2][c] == p && b[r+3][c] == p {
				return true
			}
		}
	}
	for r := 0; r < rows-3; r++ {
		for c := 0; c < cols-3; c++ {
			if b[r][c] == p && b[r+1][c+1] == p && b[r+2][c+2] == p && b[r+3][c+3] == p {
				return true
			}
		}
	}
	for r := 3; r < rows; r++ {
		for c := 0; c < cols-3; c++ {
			if b[r][c] == p && b[r-1][c+1] == p && b[r-2][c+2] == p && b[r-3][c+3] == p {
				return true
			}
		}
	}
	return false
}
	*/