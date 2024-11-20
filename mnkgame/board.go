package mnkgame

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
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
		if utf8.RuneCountInString(v) > b.rowLabelSize {
			b.rowLabelSize = utf8.RuneCountInString(v)
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

// BoardOptions packages up the various settings used when rendering the game board.
type BoardOptions struct {
	HasOuterBorder bool // Should there be a line around the labels outside the main board.
	HasInnerBorder bool // Should there be a line around the main board area.
	HasInnerGrid   bool // Should we render the lines separating each row and column.
	HasLabels      bool // Do we have labels to show.
	LabelWidth     int  // Width of longest label to be displayed.
	MarkerWidth    int  // Width of the widest player marker symbol.
	Padding        int  // Amount of whitespace on either side of labels and markers.
}

// Various board border and separator tokens.
// See https://en.wikipedia.org/wiki/Box_Drawing for more symbols.
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

	// whiteSpace is a block of text we can substring from.
	whiteSpace = "                                            "
)

var (
	cellBorder = []rune(strings.Repeat(lineHorizontal, 20))
)

var (
	boardCache = map[string]string{}
)

// String returns a fixed width layout text version of the current boards state.
//
// TODO(rsned): Consider renaming this method and leaving String() as a simpler
// state dump of the instance.
func (b *Board) String() string {

	return b.renderBoard(nil)
}

// generateStaticElements computes the dimensions of the board and renders
// the parts of the board that don't change every iteration for the rendering
// throughout the remainder of the run.
func (b *Board) generateStaticElements(bo *BoardOptions) {
	// Figure out the overall width of the output starting with the number of
	// columns plus padding on either side.
	boardWidth := b.cols * (bo.MarkerWidth + 2*bo.Padding)
	if bo.HasOuterBorder {
		boardWidth += utf8.RuneCountInString(cornerTopLeftThick) +
			utf8.RuneCountInString(cornerTopRightThick)
	}
	if bo.HasInnerBorder {
		boardWidth += utf8.RuneCountInString(cornerTopLeft) +
			utf8.RuneCountInString(cornerTopRight)
	}
	if bo.HasInnerGrid {
		boardWidth += (b.cols - 1) * utf8.RuneCountInString(lineVertical)
	}
	if bo.HasLabels {
		// Add spacing to left and right sides of the board.
		boardWidth += 2 * (bo.LabelWidth + 2*bo.Padding)
	}

	// -----------------------------------------------
	// Top and bottom outer border
	// -----------------------------------------------

	tob := cornerTopLeftThick +
		strings.Repeat(lineHorizontalThick, boardWidth-2) +
		cornerTopRightThick +
		"\n"
	boardCache["topOuterBorder"] = tob

	bob := cornerBottomLeftThick +
		strings.Repeat(lineHorizontalThick, boardWidth-2) +
		cornerBottomRightThick +
		"\n"
	boardCache["botOuterBorder"] = bob

	// -----------------------------------------------
	// Top and bottom labels
	// -----------------------------------------------

	var rowBuf bytes.Buffer
	// If has outer border
	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}
	// Row label padding
	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}
	// Inner border padding
	if bo.HasInnerBorder {
		rowBuf.WriteString(" ")
	}
	for i, col := range b.colLabels {
		rowBuf.WriteString(fmt.Sprintf("%s%s%s", whiteSpace[0:bo.Padding],
			col, whiteSpace[0:bo.Padding]))
		if bo.HasInnerGrid && i != b.cols-1 {
			rowBuf.WriteString(" ")
		}
	}
	// Inner border padding
	if bo.HasInnerBorder {
		rowBuf.WriteString(" ")
	}
	// Row label padding
	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}
	// If has outer border
	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}
	rowBuf.WriteString("\n")
	boardCache["colLabels"] = rowBuf.String()

	// -----------------------------------------------
	// Top and bottom inner borders
	// -----------------------------------------------

	rowBuf.Reset()
	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}
	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}

	// If has inner border
	if bo.HasInnerBorder {
		rowBuf.WriteString(cornerTopLeft)
	}

	for i := range b.cols {
		if bo.HasInnerBorder {
			rowBuf.WriteString(string(cellBorder[0:(bo.MarkerWidth + 2*bo.Padding)]))
		} else {
			rowBuf.WriteString(whiteSpace[0 : bo.MarkerWidth+2*bo.Padding])
		}
		if i != b.cols-1 {
			if bo.HasInnerBorder {
				if bo.HasInnerGrid {
					rowBuf.WriteString(teeDown)
				} else {
					rowBuf.WriteString(lineHorizontal)
				}
			}
		}
	}

	if bo.HasInnerBorder {
		rowBuf.WriteString(cornerTopRight)
	}

	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}

	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}

	rowBuf.WriteString("\n")

	boardCache["topInnerBorder"] = rowBuf.String()

	// Bottom inner border
	rowBuf.Reset()
	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}
	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}

	// If has inner border
	if bo.HasInnerBorder {
		rowBuf.WriteString(cornerBottomLeft)
	}

	for i := range b.cols {
		if bo.HasInnerBorder {
			rowBuf.WriteString(string(cellBorder[0:(bo.MarkerWidth + 2*bo.Padding)]))
		} else {
			rowBuf.WriteString(whiteSpace[0 : bo.MarkerWidth+2*bo.Padding])
		}
		if i != b.cols-1 {
			if bo.HasInnerBorder {
				if bo.HasInnerGrid {
					rowBuf.WriteString(teeUp)
				} else {
					rowBuf.WriteString(lineHorizontal)
				}
			} else {
				rowBuf.WriteString(" ")
			}
		}
	}

	if bo.HasInnerBorder {
		rowBuf.WriteString(cornerBottomRight)
	}

	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}

	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}

	rowBuf.WriteString("\n")

	boardCache["bottomInnerBorder"] = rowBuf.String()

	// -----------------------------------------------
	// Inner grid separator lines
	// -----------------------------------------------

	// TODO(rsned): To allow for vertical padding, generate the same row without
	// the marker values added in.

	rowBuf.Reset()
	// If has outer border
	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}
	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}

	for i := range b.cols {
		switch i {
		case 0: // Opening cell of the line.

			if bo.HasInnerBorder {
				if bo.HasInnerGrid {
					rowBuf.WriteString(teeLeft)
				} else {
					rowBuf.WriteString(lineVertical)
				}
			}
			rowBuf.WriteString(string(cellBorder[0:(bo.MarkerWidth + 2*bo.Padding)]))
			rowBuf.WriteString(cross)
		case b.cols - 1: // Closing cell of the line.

			rowBuf.WriteString(string(cellBorder[0:(bo.MarkerWidth + 2*bo.Padding)]))

			if bo.HasInnerBorder {
				if bo.HasInnerGrid {
					rowBuf.WriteString(teeRight)
				} else {
					rowBuf.WriteString(lineVertical)
				}
			}
		default:
			// Inner cells along the separator row are basically "---+"
			if bo.HasInnerGrid {
				rowBuf.WriteString(string(cellBorder[0:(bo.MarkerWidth + 2*bo.Padding)]))
				rowBuf.WriteString(cross)
			}
		}
	}

	if bo.HasLabels {
		rowBuf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
	}

	if bo.HasOuterBorder {
		rowBuf.WriteString(lineVerticalThick)
	}
	rowBuf.WriteString("\n")
	boardCache["innerSeparator"] = rowBuf.String()

}

func (b *Board) renderBoard(bo *BoardOptions) string {
	if bo == nil {
		bo = &BoardOptions{
			HasOuterBorder: true,
			HasInnerBorder: true,
			HasInnerGrid:   true,
			HasLabels:      b.hasLabels,
			LabelWidth:     b.rowLabelSize,
			MarkerWidth:    1,
			Padding:        4,
		}
	}

	var once sync.Once
	once.Do(func() { b.generateStaticElements(bo) })
	var buf bytes.Buffer

	if bo.HasOuterBorder {
		buf.WriteString(boardCache["topOuterBorder"])
	}

	if bo.HasLabels {
		buf.WriteString(boardCache["colLabels"])
	}

	// If has inner border
	if bo.HasInnerBorder {
		buf.WriteString(boardCache["topInnerBorder"])
	}

	// Main board elements
	for i, row := range b.cells {
		if bo.HasOuterBorder {
			buf.WriteString(lineVerticalThick)
		}
		if bo.HasLabels {
			buf.WriteString(fmt.Sprintf("%s%s%s", whiteSpace[0:bo.Padding],
				b.rowLabels[i], whiteSpace[0:bo.Padding]))
		}
		if bo.HasInnerBorder {
			buf.WriteString(lineVertical)
		}

		// For each active cell in this row of the board
		for i, col := range row {
			buf.WriteString(fmt.Sprintf("%s%s%s", whiteSpace[0:bo.Padding],
				col, whiteSpace[0:bo.Padding]))
			if i != b.cols-1 {
				if bo.HasInnerGrid {
					buf.WriteString(lineVertical)
				}
			}
		}

		if bo.HasInnerBorder {
			buf.WriteString(lineVertical)
		}
		if bo.HasLabels {
			buf.WriteString(whiteSpace[0 : bo.LabelWidth+2*bo.Padding])
		}
		if bo.HasOuterBorder {
			buf.WriteString(lineVerticalThick)
		}
		buf.WriteString("\n")

		if bo.HasInnerGrid && i != b.cols-1 {
			buf.WriteString(boardCache["innerSeparator"])
		}
		if i == b.cols-1 && bo.HasInnerBorder {
			buf.WriteString(boardCache["bottomInnerBorder"])
		}
	}

	if bo.HasLabels {
		buf.WriteString(boardCache["colLabels"])
	}

	if bo.HasOuterBorder {
		buf.WriteString(boardCache["botOuterBorder"])
	}

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
