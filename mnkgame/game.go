package mnkgame

// Outcome is an enumeration of the various possible states of a game.
type Outcome int

// Define the enumeration of outcomes.
const (
	OutcomeIncomplete Outcome = iota
	OutcomeWin
	OutcomeDraw
	OutcomeLoss
)

func (r Outcome) String() string {
	switch r {
	case OutcomeWin:
		return "Win"
	case OutcomeDraw:
		return "Draw"
	case OutcomeLoss:
		return "Loss"
	default:
		return "Game Unfinished"
	}
}

/*
An MNKGame is an abstract board game in which two players take turns in placing a stone
of their color on an m-by-n board, the winner being the player who first gets k stones of
their own color in a row, horizontally, vertically, or diagonally.
*/
type MNKGame struct {
	name string

	rows int
	cols int
	size int

	player1 *Player
	player2 *Player

	board *Board
}

// RenderBoard returns a string representation of the current board state.
func (t *MNKGame) RenderBoard() string {
	return t.board.String()
}

// OpenPositions returns a list of all the open positions on the board.
func (t *MNKGame) OpenPositions() []string {
	return t.board.OpenPositions()
}

// PotentialMoves returns a list of potential moves available.
//
// TODO(rsned): Augment this to support games like Nine Mens Morris and others
// that allow markers to move after they have been played.
func (t *MNKGame) PotentialMoves() []string {
	return t.board.OpenPositions()
}

// ApplyMove attempts to apply the users choice of move. If any errors occur,
// such as an illegal move, the error will be non-nil.
func (t *MNKGame) ApplyMove(player *Player, move string) error {
	return t.board.ApplyMove(player, move)
}

// Outcome reports the current status of the game for each player.
//
// TODO(rsned): Convert this to take a player and return their outcome to
// make it easier to simplify the game loop.
func (t *MNKGame) Outcome() (Outcome, Outcome) {
	return t.board.Outcome()
}

// TicTacToe returns a new instance of an m-n-k game as defined by the common Tic Tac Toe rules.
func TicTacToe(p1, p2 *Player) *MNKGame {
	g := &MNKGame{
		name: "Tic-Tac-Toe",
		rows: 3,
		cols: 3,
		size: 3,

		player1: p1,
		player2: p2,
	}

	g.player1.marker = MarkerX
	g.player2.marker = MarkerWhiteStone

	g.board = newBoard(g.rows, g.cols, g.size)

	// For tic-tac-toe we use these common labels.
	// TL -    Top Left, TC -    Top Center, TR -    Top Right,
	// CL - Center Left, CC - Center Center, CR - Center Right,
	// BL - Bottom Left, BC - Bottom Center, BR - Bottom Right,
	g.board.SetLabels([]string{"T", "C", "B"}, []string{"L", "C", "R"})

	return g
}

// Connect4 returns a new instance using the parameters in a connect 4 game.
func Connect4(p1, p2 *Player) *MNKGame {
	g := &MNKGame{
		name: "Connect 4",
		rows: 6,
		cols: 7,
		size: 4,

		player1: p1,
		player2: p2,
	}

	g.board = newBoard(g.rows, g.cols, g.size)
	g.board.SetLabels([]string{"", "", "", "", "", ""},
		[]string{"1", "2", "3", "4", "5", "6", "7"})

	return g
}

/*
TODO(rsned): Other common game options include:

Gomoku
15x15 x 5

Order and Chaos is a variant of the game tic-tac-toe on a 6×6 gameboard with 5 in a row

Something like Three Mens Morris or Nine Mens Morris would require a little more logic
in the OpenPositions and ApplyMove.


TODO(rsned): Update game to allow custom rule handling for moves.  e.g. Connect4
only takes moves using columns but not rows, and then 'gravity' moves the marker
to the next available row slot in the column.
*/
