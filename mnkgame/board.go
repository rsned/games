package mnkgame

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Board represents the cells and their states comprising an m-n-k game.
// The board has a variety of optionally set attributes to make parts
// of the game clearer such as custom labels.
//
// TODO(rsned): Separate rendering of board from its state management.
type Board struct {
	// the dimension of the board.
	rows int
	cols int

	// targetSize is the number in a row that must match to win.
	// This value is capped at the shortest dimension of the board.
	targetSize int

	// cells is the actual board layout rows x cols in size.
	cells [][]Marker

	hasLabels bool

	// If there are custom or game specific labels for the boards dimensions
	// commonly used in notation for the game. e.g., 1-8, a-h, etc.
	rowLabels []string
	colLabels []string

	// If there are labels, we only need to keep track of how long the longest
	// row label is since the column label size is the rest of the string of the
	// move string. (Note that the various N-in-a-row type games only set one
	// cell compared to something like chess or checkers where this is a start
	// and destination to the movement of a piece.
	//
	// This field assumes a fixed size for the row labels. e.g. if there are more
	// than 9 rows, the row labels will all be padded ahead of time to take the
	// same space. We do not pad or align the labels.
	rowLabelSize int

	// These maps are the reverse move string lookups from human notation
	// back to coords.
	rowLabelMap map[string]int
	colLabelMap map[string]int

	// winTests is the set of all N-in-a-row that fit within the current boards
	// dimensions. It is precomputed once at start time so the per-move checking
	// can just iterate over it.
	winTests CoordsList
}

// newBoard creates a new instance of a board of the given dimensions and n-in-a-row
// target and fills the board with the empty markers.
func newBoard(rows, cols, targetSize int) *Board {
	b := &Board{
		rows:       rows,
		cols:       cols,
		targetSize: targetSize,
	}
	if b.targetSize > b.rows || b.targetSize > b.cols {
		// TODO(rsned): find the min of rows and cols.
		b.targetSize = b.rows
	}

	b.winTests = b.generateAllWinningCoordinateSets()

	// Initialize the board to the required dimensions and pre-fill it with
	// the empty marker.
	b.cells = make([][]Marker, rows, rows)
	for i := range b.cells {
		b.cells[i] = make([]Marker, cols, cols)
		for k := 0; k < cols; k++ {
			b.cells[i][k] = MarkerEmpty
		}
	}

	return b
}

// SetLabels sets the given set of labels for the rows and columns in the
// board and updates the corresponding state elements of the board.
func (b *Board) SetLabels(rowLabels, colLabels []string) {
	if len(rowLabels) == 0 || len(colLabels) == 0 {
		return
	}
	b.hasLabels = true

	// Basic sanity checks on the number of labels.
	if len(rowLabels) < b.rows {
		rowLabels = append(rowLabels, make([]string, b.rows-len(rowLabels))...)
	}

	// Basic sanity checks on the number of labels.
	if len(colLabels) < b.cols {
		colLabels = append(colLabels, make([]string, b.cols-len(colLabels))...)
	}

	b.rowLabels = rowLabels[:b.rows]
	b.colLabels = colLabels[:b.cols]
	b.rowLabelMap = map[string]int{}
	b.colLabelMap = map[string]int{}

	b.rowLabelSize = 0
	for i, v := range b.rowLabels {
		b.rowLabelMap[v] = i
		if len(v) > b.rowLabelSize {
			b.rowLabelSize = len(v)
		}
	}

	for i, v := range b.colLabels {
		b.colLabelMap[v] = i
	}

}

// decodeMove applies the reverse transformation of the strings generated in
// OpenPositions to map back to a board coordinate.
func (b *Board) decodeMove(move string) (Coord, bool) {
	if b.hasLabels {
		row := move[0:b.rowLabelSize]
		col := move[b.rowLabelSize:]
		var c Coord
		var ok1, ok2 bool
		c.Row, ok1 = b.rowLabelMap[row]
		c.Col, ok2 = b.colLabelMap[col]
		return c, ok1 && ok2
	}

	var coord Coord

	// TODO(rsned): If these assumptions change, update this.
	//
	// User strings are 1-based, board cell coordinates are 0-based.
	// Separator character is ","
	parts := strings.Split(move, ",")
	if len(parts) < 2 {
		return coord, false
	}

	if i, err := strconv.Atoi(parts[0]); (err == nil) && (i > 0 && i <= b.rows) {
		// i was a real integer and within the board bounds.
		coord.Row = i - 1
	} else {
		return coord, false
	}

	if i, err := strconv.Atoi(parts[1]); (err == nil) && (i > 0 && i <= b.cols) {
		// i was a real integer and within the board bounds.
		coord.Col = i - 1
	} else {
		return coord, false
	}
	return coord, true
}

// ApplyMove applies the given move for the given player to the board.
// If there are errors preventing the move, they are returned.
func (b *Board) ApplyMove(player *Player, move string) error {
	m, ok := b.decodeMove(move)
	if !ok {
		return fmt.Errorf("Unable to decipher the requested move: %q", move)
	}

	if b.cells[m.Row][m.Col] != MarkerEmpty {
		return fmt.Errorf("Move not available")
	}

	b.cells[m.Row][m.Col] = player.marker
	return nil
}

// OpenPositions returns the set all possible cells that have not yet been filled.
// If there are notation labels, those values are returned. Otherwise, a list of
// cell coordinates is returned.
func (b *Board) OpenPositions() []string {
	var open []string
	for i, row := range b.cells {
		for j, col := range row {
			if col == MarkerEmpty {
				var l string
				if b.hasLabels {
					// TODO(rsned): implement a GenerateNotation to make these.
					l = fmt.Sprintf("%s%s", b.rowLabels[i], b.colLabels[j])
				} else {
					l = fmt.Sprintf("%d,%d", i+1, j+1)
				}
				open = append(open, l)
			}
		}
	}
	return open
}

const (
	cellBorder = lineHorizontal + lineHorizontal + lineHorizontal
)

// Board border and separator tokens.
// See https://en.wikipedia.org/wiki/Box_Drawing  for more symbols.
const (
	cornerTopLeft          = "┌"
	cornerTopLeftThick     = "┏"
	cornerTopRight         = "┐"
	cornerTopRightThick    = "┓"
	cornerBottomLeft       = "└"
	cornerBottomLeftThick  = "┗"
	cornerBottomRight      = "┘"
	cornerBottomRightThick = "┛"
	lineHorizontal         = "─"
	lineHorizontalThick    = "━"
	lineVertical           = "│"
	lineVerticalThick      = "┃"
	cross                  = "┼"
	crossThick             = "╋"
	teeLeft                = "├"
	teeLeftThick           = "┣"
	teeRight               = "┤"
	teeRightThick          = "┫"
	teeUp                  = "┴"
	teeUpThick             = "┻"
	teeDown                = "┬"
	teeDownThick           = "┳"
)

var (
	boardCache = map[string]string{}
)

// String returns a fixed width layout text version of the current boards state.
//
// TODO(rsned): Consider renaming this method and leaving String() as a simpler
// state dump of the instance.
func (b *Board) String() string {

	return b.renderBoard()
}

func (b *Board) renderBoard() string {
	var buf bytes.Buffer

	const cellWidth = 3
	boardWidth := b.cols*cellWidth /* cell contents */ +
		(b.cols - 1) /* innerGrid */ +
		2 /* inner border */ +
		2 /* outer border */ +
		(2 * 3) /* label width */

	// If has outer border.
	if _, ok := boardCache["topOuterBorder"]; !ok {
		tob := cornerTopLeftThick +
			strings.Repeat(lineHorizontalThick, boardWidth-2) +
			cornerTopRightThick +
			"\n"
		boardCache["topOuterBorder"] = tob
	}
	buf.WriteString(boardCache["topOuterBorder"])

	if b.hasLabels {
		if _, ok := boardCache["colLabels"]; !ok {
			var rowBuf bytes.Buffer
			// If has outer border
			rowBuf.WriteString(lineVerticalThick)
			// Row label padding
			if b.hasLabels {
				rowBuf.WriteString("   ")
			}
			// Inner border padding
			rowBuf.WriteString(" ")
			for i, col := range b.colLabels {
				rowBuf.WriteString(fmt.Sprintf(" %s ", col))
				// If has inner grid
				if i != b.cols-1 {
					rowBuf.WriteString(" ")
				}
			}
			// Inner border padding
			rowBuf.WriteString(" ")
			// Row label padding
			if b.hasLabels {
				rowBuf.WriteString("   ")
			}
			// If has outer border
			rowBuf.WriteString(lineVerticalThick)
			rowBuf.WriteString("\n")
			boardCache["colLabels"] = rowBuf.String()
		}
		buf.WriteString(boardCache["colLabels"])
	}

	// If has inner border
	// Top inner border row
	if _, ok := boardCache["topInnerBorder"]; !ok {
		var rowBuf bytes.Buffer
		// If has outer border
		rowBuf.WriteString(lineVerticalThick)
		if b.hasLabels {
			rowBuf.WriteString("   ")
		}

		// If has inner border
		rowBuf.WriteString(cornerTopLeft)

		for i := range b.cols {
			rowBuf.WriteString(cellBorder)
			// If has inner grid
			if i != b.cols-1 {
				rowBuf.WriteString(teeDown)
				// } else {
				// rowBuf.WriteString(lineHorizontal)
				// }
			}
		}
		rowBuf.WriteString(cornerTopRight)

		if b.hasLabels {
			rowBuf.WriteString("   ")
		}
		// If has outer border
		rowBuf.WriteString(lineVerticalThick)
		rowBuf.WriteString("\n")
		boardCache["topInnerBorder"] = rowBuf.String()
	}

	// If has inner border
	buf.WriteString(boardCache["topInnerBorder"])

	// Main board elements

	for i, row := range b.cells {
		// If has outer border
		buf.WriteString(lineVerticalThick)
		if b.hasLabels {
			buf.WriteString(fmt.Sprintf(" %s ", b.rowLabels[i]))
		}

		// If has inner border
		buf.WriteString(lineVertical)

		// For each active cell in this row of the board
		for i, col := range row {
			buf.WriteString(fmt.Sprintf(" %s ", col))
			if i != b.cols-1 {
				buf.WriteString(lineVertical)
			}
		}

		// If has inner border.
		buf.WriteString(lineVertical)

		// If has labels.
		if b.hasLabels {
			buf.WriteString("   ")
		}
		// If has outer border
		buf.WriteString(lineVerticalThick)
		buf.WriteString("\n")

		if i == b.cols-1 {
			// Bottom inner border row
			if _, ok := boardCache["bottomInnerBorder"]; !ok {
				var rowBuf bytes.Buffer
				// If has outer border
				rowBuf.WriteString(lineVerticalThick)
				if b.hasLabels {
					rowBuf.WriteString("   ")
				}
				// If has inner border
				rowBuf.WriteString(cornerBottomLeft)

				for i := range b.cols {
					rowBuf.WriteString(cellBorder)
					// If has inner grid
					if i != b.cols-1 {
						rowBuf.WriteString(teeUp)
						// } else {
						// rowBuf.WriteString(lineHorizontal)
						// }
					}
				}
				rowBuf.WriteString(cornerBottomRight)

				if b.hasLabels {
					rowBuf.WriteString("   ")
				}
				// If has outer border
				rowBuf.WriteString(lineVerticalThick)
				rowBuf.WriteString("\n")
				boardCache["bottomInnerBorder"] = rowBuf.String()
			}
			buf.WriteString(boardCache["bottomInnerBorder"])
			break
		}

		// If has inner grid lines
		if _, ok := boardCache["innerSeparator"]; !ok {
			var rowBuf bytes.Buffer
			// If has outer border
			rowBuf.WriteString(lineVerticalThick)
			if b.hasLabels {
				rowBuf.WriteString("   ")
			}

			for i := range row {
				switch i {
				case 0:
					// If has inner border AND NOT is last row
					rowBuf.WriteString(teeLeft)
					rowBuf.WriteString(cellBorder)
					rowBuf.WriteString(cross)
				case b.cols - 1:
					// If has inner border AND NOT is last row
					rowBuf.WriteString(cellBorder)
					rowBuf.WriteString(teeRight)
				default:
					rowBuf.WriteString(cellBorder)
					rowBuf.WriteString(cross)
				}
			}

			if b.hasLabels {
				rowBuf.WriteString("   ")
			}
			// If has outer border
			rowBuf.WriteString(lineVerticalThick)
			rowBuf.WriteString("\n")
			boardCache["innerSeparator"] = rowBuf.String()
		}
		// If has inner grid lines
		buf.WriteString(boardCache["innerSeparator"])
	}

	if b.hasLabels {
		buf.WriteString(boardCache["colLabels"])
	}

	// If has outer border.
	if _, ok := boardCache["botOuterBorder"]; !ok {
		tob := cornerBottomLeftThick +
			strings.Repeat(lineHorizontalThick, boardWidth-2) +
			cornerBottomRightThick +
			"\n"
		boardCache["botOuterBorder"] = tob
	}
	// If has outer border
	buf.WriteString(boardCache["botOuterBorder"])

	return buf.String()
}

// checkOutcome tests a given set of coords for the given player to see if
// there is a full match.
func (b *Board) checkOutcome(coords Coords, p *Player) bool {
	if len(coords) != b.targetSize {
		return false
	}

	win := true
	for _, c := range coords {
		win = win && (b.cells[c.Row][c.Col] == p.marker)
	}
	return win
}

// Outcome reports the game outcome state for both players.
func (b *Board) Outcome() (player1, player2 Outcome) {
	var win bool
	// Check P1
	for _, coords := range b.winTests {
		win = win || b.checkOutcome(coords, Player1)
	}
	if win {
		return OutcomeWin, OutcomeLoss
	}

	// Check P2
	win = false
	for _, coords := range b.winTests {
		win = win || b.checkOutcome(coords, Player2)
	}
	if win {
		return OutcomeLoss, OutcomeWin
	}

	// Check if board is full.
	if len(b.OpenPositions()) == 0 {
		return OutcomeDraw, OutcomeDraw
	}

	return OutcomeIncomplete, OutcomeIncomplete
}

// generateAllWinningCoordinateSets is used to figure out based on the board
// parameters all sets of coordinates that represent winning sequences.
//
// TODO(rsned): Consider moving this to a standalone method instead of relying
// on the method to get values from the Board.
func (b *Board) generateAllWinningCoordinateSets() CoordsList {
	potentialWins := CoordsList{}

	// Start at the origin corner and walk all the cells in order, top
	// left to botton right. For each cell attempt to generate the horizontal,
	// vertical, and both diagonals going rightward and downward.
	// If there are Board.targetSize cells available in the given direction,
	// grab that many coordinates and save it. Then move on to the next cell
	// until complete.
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			// Horizontal
			if col <= (b.cols - b.targetSize) {
				vals := Coords{}
				for k := 0; k < b.targetSize; k++ {
					vals.Add(Coord{
						Row: row,
						Col: col + k,
					})
				}
				slices.SortFunc(vals, coordCompare)
				potentialWins.Add(vals)
			}

			// Vertical
			if row <= (b.rows - b.targetSize) {
				vals := Coords{}
				for k := 0; k < b.targetSize; k++ {
					vals.Add(Coord{
						Row: row + k,
						Col: col,
					})
				}
				slices.SortFunc(vals, coordCompare)
				potentialWins.Add(vals)
			}

			// Diagonal TL->BR
			// Because we increase by 1 step in both directions,
			// The last row and col we can use is the difference between
			// the number of rows minus the target size.  e.g. With a
			// target of 2 and 2 rows, the only cell is row 0. Same
			// on the col side.  If the rows is 5 and target is 2, then
			// row 3 would be the final choice. etc. etc.
			if row <= (b.rows-b.targetSize) && col <= (b.cols-b.targetSize) {
				vals := Coords{}
				for k := 0; k < b.targetSize; k++ {
					vals.Add(Coord{
						Row: row + k,
						Col: col + k,
					})
				}
				slices.SortFunc(vals, coordCompare)
				potentialWins.Add(vals)
			}

			// Diagonal BL->TR
			// This direction, the rows will grow 'upwards' aka, row will
			// decrease while col will increase. The row just has to
			// be at least targetSize so it won't exceed the bounds when
			// subtracting targetSize steps. Columns just uses the same as prev.
			if (row+1) >= b.targetSize && col <= b.cols-b.targetSize {
				vals := Coords{}
				for k := 0; k < b.targetSize; k++ {
					vals.Add(Coord{
						Row: row - k,
						Col: col + k,
					})
				}
				slices.SortFunc(vals, coordCompare)
				potentialWins.Add(vals)
			}
		}
	}

	slices.SortFunc(potentialWins, coordSliceCompare)
	return potentialWins
}
