package web

// Game contient les informations d'une partie
type Game struct {
	Player1       string
	Player2       string
	Difficulty    string
	Rows          int
	Cols          int
	Grid          [][]int
	CurrentPlayer int
	Winner        int
	Draw          bool
	ColRange      []int
	TurnCount     int      // compteur de tours pour la gravité
	GravityDown   bool     // true = gravité normale, false = inversée
	WinningCells  [][2]int // positions des jetons gagnants
	IsVsBot		  bool
	BotLevel	  string
}

// Cell représente une cellule enrichie pour le rendu HTML
type Cell struct {
	Value     int  // 0, 1 ou 2
	IsWinning bool // true si cette cellule fait partie des jetons gagnants
}

// NewGame crée une nouvelle partie en fonction du niveau
func NewGame(player1, player2, difficulty string) *Game {
	var rows, cols int
	switch difficulty {
	case "easy":
		rows, cols = 6, 7
	case "normal":
		rows, cols = 6, 9
	case "hard":
		rows, cols = 7, 8
	default:
		rows, cols = 6, 7
	}

	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}

	colRange := make([]int, cols)
	for i := 0; i < cols; i++ {
		colRange[i] = i
	}

	return &Game{
		Player1:       player1,
		Player2:       player2,
		Difficulty:    difficulty,
		Rows:          rows,
		Cols:          cols,
		Grid:          grid,
		CurrentPlayer: 1,
		ColRange:      colRange,
		TurnCount:     0,
		GravityDown:   true,
	}
}

// PlaceToken place un jeton dans la colonne spécifiée en tenant compte de la gravité
func (g *Game) PlaceToken(col int) bool {
	if g.GravityDown {
		for row := g.Rows - 1; row >= 0; row-- {
			if g.Grid[row][col] == 0 {
				g.Grid[row][col] = g.CurrentPlayer
				g.TurnCount++
				if g.TurnCount%10 == 0 {
					g.GravityDown = !g.GravityDown
				}
				return true
			}
		}
	} else {
		for row := 0; row < g.Rows; row++ {
			if g.Grid[row][col] == 0 {
				g.Grid[row][col] = g.CurrentPlayer
				g.TurnCount++
				if g.TurnCount%10 == 0 {
					g.GravityDown = !g.GravityDown
				}
				return true
			}
		}
	}
	return false
}

// CheckDraw vérifie si la grille est pleine
func (g *Game) CheckDraw() bool {
	for _, row := range g.Grid {
		for _, cell := range row {
			if cell == 0 {
				return false
			}
		}
	}
	return true
}

// CheckWin vérifie les alignements de 4 et enregistre les positions gagnantes
func (g *Game) CheckWin() int {
	directions := [][2]int{
		{0, 1},  // horizontal →
		{1, 0},  // vertical ↓
		{1, 1},  // diagonale ↘
		{1, -1}, // diagonale ↙
	}

	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			player := g.Grid[r][c]
			if player == 0 {
				continue
			}

			for _, d := range directions {
				cells := [][2]int{{r, c}}

				for step := 1; step < 4; step++ {
					nr := r + d[0]*step
					nc := c + d[1]*step

					if nr < 0 || nr >= g.Rows || nc < 0 || nc >= g.Cols {
						break
					}
					if g.Grid[nr][nc] == player {
						cells = append(cells, [2]int{nr, nc})
					} else {
						break
					}
				}

				if len(cells) == 4 {
					g.Winner = player
					g.WinningCells = cells
					return player
				}
			}
		}
	}
	return 0
}

// IsWinningCell vérifie si une cellule fait partie des cellules gagnantes
func (g *Game) IsWinningCell(row, col int) bool {
	for _, cell := range g.WinningCells {
		if cell[0] == row && cell[1] == col {
			return true
		}
	}
	return false
}

// RenderedGrid retourne une grille enrichie pour le template HTML
func (g *Game) RenderedGrid() [][]Cell {
	grid := make([][]Cell, g.Rows)
	for r := 0; r < g.Rows; r++ {
		row := make([]Cell, g.Cols)
		for c := 0; c < g.Cols; c++ {
			row[c] = Cell{
				Value:     g.Grid[r][c],
				IsWinning: g.IsWinningCell(r, c),
			}
		}
		grid[r] = row
	}
	return grid
}

// Reset vide la grille et réinitialise l'état de la partie
func (g *Game) Reset() {
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			g.Grid[r][c] = 0
		}
	}
	g.CurrentPlayer = 1
	g.Winner = 0
	g.Draw = false
	g.TurnCount = 0
	g.GravityDown = true
	g.WinningCells = nil
}