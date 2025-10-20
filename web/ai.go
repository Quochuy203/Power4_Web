package web

import (
	"math/rand"
	"time"
)

// ChooseBotMove choisit une colonne pour le robot selon la difficulté
func (g *Game) ChooseBotMove() int {
	rand.Seed(time.Now().UnixNano())
	
	validMoves := g.getValidMoves()
	if len(validMoves) == 0 {
		return -1
	}

	switch g.BotLevel {
	case "easy":
		// Facile: joue aléatoirement
		return validMoves[rand.Intn(len(validMoves))]
		
	case "normal":
		// Normal: cherche coup gagnant, puis blocage, puis stratégique
		return g.chooseNormalMove(validMoves)
		
	case "hard":
		// Difficile: utilise minimax avec profondeur 5
		return g.chooseMinimax(validMoves, 5)
		
	default:
		return validMoves[rand.Intn(len(validMoves))]
	}
}

// getValidMoves retourne toutes les colonnes jouables
func (g *Game) getValidMoves() []int {
	validMoves := []int{}
	for c := 0; c < g.Cols; c++ {
		// Vérifie si la colonne n'est pas pleine (selon la gravité)
		if g.GravityDown {
			if g.Grid[0][c] == 0 {
				validMoves = append(validMoves, c)
			}
		} else {
			if g.Grid[g.Rows-1][c] == 0 {
				validMoves = append(validMoves, c)
			}
		}
	}
	return validMoves
}

// chooseNormalMove : stratégie pour difficulté normale
func (g *Game) chooseNormalMove(validMoves []int) int {
	// 1. Cherche un coup gagnant pour le robot (joueur 2)
	for _, c := range validMoves {
		tmp := g.clone()
		tmp.PlaceTokenForPlayer(c, 2)
		if tmp.CheckWin() == 2 {
			return c // GAGNER !
		}
	}

	// 2. Bloque l'adversaire s'il peut gagner au prochain coup
	for _, c := range validMoves {
		tmp := g.clone()
		tmp.PlaceTokenForPlayer(c, 1)
		if tmp.CheckWin() == 1 {
			return c // BLOQUER !
		}
	}

	// 3. Cherche à créer des opportunités (alignement de 3)
	bestScore := -1000
	bestCol := validMoves[0]
	
	for _, c := range validMoves {
		tmp := g.clone()
		tmp.PlaceTokenForPlayer(c, 2)
		score := tmp.evaluateBoardAdvanced(2)
		
		if score > bestScore {
			bestScore = score
			bestCol = c
		}
	}

	return bestCol
}

// chooseMinimax : utilise l'algorithme minimax avec élagage alpha-beta
func (g *Game) chooseMinimax(validMoves []int, depth int) int {
	bestScore := -999999
	bestCol := validMoves[0]

	for _, c := range validMoves {
		tmp := g.clone()
		tmp.PlaceTokenForPlayer(c, 2)
		
		// Vérifie victoire immédiate
		if tmp.CheckWin() == 2 {
			return c
		}
		
		score := tmp.minimax(depth-1, -999999, 999999, false)
		
		if score > bestScore {
			bestScore = score
			bestCol = c
		}
	}

	return bestCol
}

// minimax : algorithme minimax avec élagage alpha-beta
func (g *Game) minimax(depth int, alpha, beta int, isMaximizing bool) int {
	// Vérifie les conditions terminales
	winner := g.CheckWin()
	if winner == 2 {
		return 1000 + depth // Préfère gagner rapidement
	}
	if winner == 1 {
		return -1000 - depth // L'adversaire gagne
	}
	if g.CheckDraw() || depth == 0 {
		return g.evaluateBoardAdvanced(2)
	}

	validMoves := g.getValidMoves()
	if len(validMoves) == 0 {
		return 0
	}

	if isMaximizing {
		// Tour du robot (joueur 2)
		maxEval := -999999
		for _, c := range validMoves {
			tmp := g.clone()
			tmp.PlaceTokenForPlayer(c, 2)
			eval := tmp.minimax(depth-1, alpha, beta, false)
			if eval > maxEval {
				maxEval = eval
			}
			if eval > alpha {
				alpha = eval
			}
			if beta <= alpha {
				break // Élagage beta
			}
		}
		return maxEval
	} else {
		// Tour de l'adversaire (joueur 1)
		minEval := 999999
		for _, c := range validMoves {
			tmp := g.clone()
			tmp.PlaceTokenForPlayer(c, 1)
			eval := tmp.minimax(depth-1, alpha, beta, true)
			if eval < minEval {
				minEval = eval
			}
			if eval < beta {
				beta = eval
			}
			if beta <= alpha {
				break // Élagage alpha
			}
		}
		return minEval
	}
}

// evaluateBoardAdvanced : évaluation avancée de la position
func (g *Game) evaluateBoardAdvanced(player int) int {
	score := 0
	opponent := 3 - player // 1 si player=2, 2 si player=1

	// Évalue toutes les fenêtres de 4 cases possibles
	// Horizontal
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols-3; c++ {
			window := []int{g.Grid[r][c], g.Grid[r][c+1], g.Grid[r][c+2], g.Grid[r][c+3]}
			score += g.evaluateWindow(window, player, opponent)
		}
	}

	// Vertical
	for r := 0; r < g.Rows-3; r++ {
		for c := 0; c < g.Cols; c++ {
			window := []int{g.Grid[r][c], g.Grid[r+1][c], g.Grid[r+2][c], g.Grid[r+3][c]}
			score += g.evaluateWindow(window, player, opponent)
		}
	}

	// Diagonale ↘
	for r := 0; r < g.Rows-3; r++ {
		for c := 0; c < g.Cols-3; c++ {
			window := []int{g.Grid[r][c], g.Grid[r+1][c+1], g.Grid[r+2][c+2], g.Grid[r+3][c+3]}
			score += g.evaluateWindow(window, player, opponent)
		}
	}

	// Diagonale ↙
	for r := 0; r < g.Rows-3; r++ {
		for c := 3; c < g.Cols; c++ {
			window := []int{g.Grid[r][c], g.Grid[r+1][c-1], g.Grid[r+2][c-2], g.Grid[r+3][c-3]}
			score += g.evaluateWindow(window, player, opponent)
		}
	}

	// Bonus pour contrôle du centre
	centerCol := g.Cols / 2
	centerCount := 0
	for r := 0; r < g.Rows; r++ {
		if g.Grid[r][centerCol] == player {
			centerCount++
		}
	}
	score += centerCount * 3

	return score
}

// evaluateWindow : évalue une fenêtre de 4 cases
func (g *Game) evaluateWindow(window []int, player, opponent int) int {
	score := 0
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

	// Évaluation selon les combinaisons
	if playerCount == 4 {
		score += 100 // Victoire
	} else if playerCount == 3 && emptyCount == 1 {
		score += 10 // 3 alignés, peut gagner
	} else if playerCount == 2 && emptyCount == 2 {
		score += 5 // 2 alignés, potentiel
	}

	// Pénalise les menaces adverses
	if opponentCount == 3 && emptyCount == 1 {
		score -= 80 // URGENT : bloquer !
	} else if opponentCount == 2 && emptyCount == 2 {
		score -= 5 // Menace future
	}

	return score
}

// clone : crée une copie de la grille
func (g *Game) clone() *Game {
	newG := *g
	newG.Grid = make([][]int, g.Rows)
	for i := range g.Grid {
		newG.Grid[i] = make([]int, g.Cols)
		copy(newG.Grid[i], g.Grid[i])
	}
	return &newG
}

// PlaceTokenForPlayer place un pion pour un joueur précis (utilisé par l'IA)
func (g *Game) PlaceTokenForPlayer(col int, player int) bool {
	if g.GravityDown {
		for row := g.Rows - 1; row >= 0; row-- {
			if g.Grid[row][col] == 0 {
				g.Grid[row][col] = player
				return true
			}
		}
	} else {
		for row := 0; row < g.Rows; row++ {
			if g.Grid[row][col] == 0 {
				g.Grid[row][col] = player
				return true
			}
		}
	}
	return false
}