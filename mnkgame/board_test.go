package mnkgame

import (
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBoardSetLabels(t *testing.T) {
	tests := []struct {
		board         *Board
		rowLabels     []string
		colLabels     []string
		wantRowLabels []string
		wantColLabels []string
	}{
		{
			// Empty inputs should result in blank padded list.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 1,
			},
			rowLabels:     nil,
			colLabels:     nil,
			wantRowLabels: nil,
			wantColLabels: nil,
		},
		{
			// Under-length inputs should result in blank padded to right length.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
			},
			rowLabels:     []string{"a"},
			colLabels:     []string{"1"},
			wantRowLabels: []string{"a", "", ""},
			wantColLabels: []string{"1", "", ""},
		},
		{
			// Mis-matched lengths of Under-length inputs should result in blank padded to right length.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
			},
			rowLabels:     []string{"a", "b"},
			colLabels:     []string{"1"},
			wantRowLabels: []string{"a", "b", ""},
			wantColLabels: []string{"1", "", ""},
		},
		{
			// Expected length inputs should stay unchanged.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
			},
			rowLabels:     []string{"a"},
			colLabels:     []string{"1"},
			wantRowLabels: []string{"a", "", ""},
			wantColLabels: []string{"1", "", ""},
		},
		{
			// Over-length inputs should be truncated to match board dimensions.
			board: &Board{
				rows:       3,
				cols:       5,
				targetSize: 1,
			},
			rowLabels:     []string{"a", "b", "c", "d", "e", "f"},
			colLabels:     []string{"0", "1", "2", "3", "4", "5"},
			wantRowLabels: []string{"a", "b", "c"},
			wantColLabels: []string{"1", "2", "3", "4", "5"},
		},
	}

	for _, test := range tests {
		test.board.SetLabels(test.rowLabels, test.colLabels)
		if !cmp.Equal(test.board.rowLabels, test.wantRowLabels, cmpopts.EquateEmpty()) {
			t.Errorf("rowLabels(%v, %v), final rowLabels = %+v, want %+v\ndiff: %+v",
				test.rowLabels, test.colLabels, test.board.rowLabels, test.wantRowLabels,
				cmp.Diff(test.board.rowLabels, test.wantRowLabels))
		}
		if !cmp.Equal(test.board.rowLabels, test.wantRowLabels, cmpopts.EquateEmpty()) {
			t.Errorf("colLabels(%v, %v), final colLabels = %+v, want %+v\ndiff: %+v",
				test.rowLabels, test.colLabels, test.board.colLabels, test.wantColLabels,
				cmp.Diff(test.board.colLabels, test.wantColLabels))
		}
	}

}

func TestBoardApplyMove(t *testing.T) {
	tests := []struct {
		board   *Board
		player  *Player
		move    string
		wantErr bool
	}{
		{
			//
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			player:  Player1,
			move:    "A1",
			wantErr: true,
		},
	}

	for _, test := range tests {
		gotErr := test.board.ApplyMove(test.player, test.move)
		if (gotErr != nil) != test.wantErr {
			t.Errorf("ApplyMove(player, %q) error != nil = %v, want %v",
				test.move, gotErr != nil, test.wantErr)
		}

	}
}

func TestBoardOpenPositions(t *testing.T) {
	tests := []struct {
		have *Board
		want []string
	}{
		{
			// Empty board, all spots available.
			have: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
				},
			},
			want: []string{"1,1", "1,2", "1,3", "2,1", "2,2", "2,3", "3,1", "3,2", "3,3"},
		},
		{
			// Partial board, most spots open.
			have: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerX, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
				},
			},
			want: []string{"1,1", "1,2", "1,3", "2,1", "2,3", "3,1", "3,2", "3,3"},
		},
		{
			// Partial board, most spots open.
			have: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerWhiteStone, MarkerEmpty},
					[]Marker{MarkerX, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerWhiteStone},
				},
			},
			want: []string{"1,1", "1,3", "2,2", "2,3", "3,1", "3,2"},
		},
		{
			// Full board, no open spots.
			have: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerX, MarkerWhiteStone, MarkerWhiteStone},
					[]Marker{MarkerX, MarkerWhiteStone, MarkerWhiteStone},
					[]Marker{MarkerX, MarkerX, MarkerX},
				},
			},
			want: []string{},
		},
	}

	for _, test := range tests {
		if got := test.have.OpenPositions(); !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("b.OpenPositions() = %v, want %v\ndiff: %s",
				got, test.want, cmp.Diff(got, test.want))
		}
	}

}

func TestBoardDecodeMove(t *testing.T) {
	tests := []struct {
		board     *Board
		rowLabels []string
		colLabels []string
		move      string
		want      Coord
		wantOK    bool
	}{
		{
			// Board has no labels, so decode needs to handle as a coord string
			// but it's not a valid coord string, so it should fail.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			move:   "A1",
			want:   Coord{Row: 1, Col: 1},
			wantOK: false,
		},
		{
			// Board has no labels, so decode needs to handle as a coord string.
			// Value does not contain supported separator.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			move:   "2:1",
			want:   Coord{Row: 1, Col: 0},
			wantOK: false,
		},
		{
			// Board has no labels, so decode needs to handle as a coord string.
			// Value is valid.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			move:   "2,1",
			want:   Coord{Row: 1, Col: 0},
			wantOK: true,
		},
		{
			// Board has no labels, so decode needs to handle as a coord string.
			// Value is bounds values.
			board: &Board{
				rows:       3,
				cols:       4,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			move:   "3,4",
			want:   Coord{Row: 2, Col: 3},
			wantOK: true,
		},
		{
			// Board has no labels, so decode needs to handle as a coord string
			// Value is outside board bounds.
			board: &Board{
				rows:       9,
				cols:       9,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			move:   "11,-51",
			want:   Coord{Row: 0, Col: 0},
			wantOK: false,
		},
		{
			// Board has labels, so decode needs to use the labels to reverse map.
			// Value is not in the range of labels
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			rowLabels: []string{"A", "B"},
			colLabels: []string{"1", "2"},
			move:      "C1",
			want:      Coord{Row: 0, Col: 0},
			wantOK:    false,
		},
		{
			// Board has labels, so decode needs to use the labels to reverse map.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			rowLabels: []string{"A", "B"},
			colLabels: []string{"1", "2"},
			move:      "A1",
			want:      Coord{Row: 0, Col: 0},
			wantOK:    true,
		},
	}

	for _, test := range tests {
		test.board.SetLabels(test.rowLabels, test.colLabels)
		got, ok := test.board.decodeMove(test.move)

		if test.wantOK != ok {
			t.Errorf("decodeMove(%q) = (%+v, %v), want %v", test.move, got, ok, test.wantOK)
			continue
		}

		// If we didn't want a failure and didn't get a failure, test the actual result.
		if ok && !got.equals(test.want) {
			t.Errorf("decodeMove(%q) = %+v, want %+v", test.move, got, test.want)
		}
	}
}

func TestBoardGenerateAllWinningCoordinateSets(t *testing.T) {
	tests := []struct {
		board *Board
		want  CoordsList
	}{
		{
			// An invalid game, target size larger than any dimension.
			board: &Board{
				rows:       2,
				cols:       3,
				targetSize: 4,
				cells: [][]Marker{
					[]Marker{MarkerEmpty},
				},
			},
			want: CoordsList{},
		},
		{
			// The smallest game. 1x1, 1 in a row. Only one winning sequence.
			board: &Board{
				rows:       1,
				cols:       1,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty},
				},
			},
			want: CoordsList{
				Coords{
					Coord{Row: 0, Col: 0},
				},
			},
		},
		{
			// 2x2, 1 in a row. Should be 4 result, one of each cell.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 1,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			want: CoordsList{
				Coords{
					Coord{Row: 0, Col: 0},
				},
				Coords{
					Coord{Row: 1, Col: 0},
				},
				Coords{
					Coord{Row: 0, Col: 1},
				},
				Coords{
					Coord{Row: 1, Col: 1},
				},
			},
		},
		{
			// 2x2, 2 in a row. Should be 6 result, 2 horizontal,
			// 2 vertical, 2 diagonal.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 2,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			want: CoordsList{
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 0, Col: 1},
				},
				Coords{
					Coord{Row: 1, Col: 0},
					Coord{Row: 1, Col: 1},
				},
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 1, Col: 0},
				},
				Coords{
					Coord{Row: 0, Col: 1},
					Coord{Row: 1, Col: 1},
				},
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 1, Col: 1},
				},
				Coords{
					Coord{Row: 0, Col: 1},
					Coord{Row: 1, Col: 0},
				},
			},
		},
		{
			// A common game, 3x3 3 in a row. Should be 3 horizontal,
			// 3 vertical, and 2 diagonal.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
				},
			},
			want: CoordsList{
				// rows
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 0, Col: 1},
					Coord{Row: 0, Col: 2},
				},
				Coords{
					Coord{Row: 1, Col: 0},
					Coord{Row: 1, Col: 1},
					Coord{Row: 1, Col: 2},
				},
				Coords{
					Coord{Row: 2, Col: 0},
					Coord{Row: 2, Col: 1},
					Coord{Row: 2, Col: 2},
				},
				// cols
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 1, Col: 0},
					Coord{Row: 2, Col: 0},
				},
				Coords{
					Coord{Row: 0, Col: 1},
					Coord{Row: 1, Col: 1},
					Coord{Row: 2, Col: 1},
				},
				Coords{
					Coord{Row: 0, Col: 2},
					Coord{Row: 1, Col: 2},
					Coord{Row: 2, Col: 2},
				},
				// diags
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 1, Col: 1},
					Coord{Row: 2, Col: 2},
				},
				Coords{
					Coord{Row: 0, Col: 2},
					Coord{Row: 1, Col: 1},
					Coord{Row: 2, Col: 0},
				},
			},
		},
		{
			// Uneven dimensions, 4x3, 3 in a row. Should be
			// 3 horizontal, 6 vertical, and 8 diag.
			board: &Board{
				rows:       4,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty, MarkerEmpty},
				},
			},
			want: CoordsList{

				// rows
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 0, Col: 1},
					Coord{Row: 0, Col: 2},
				},
				Coords{
					Coord{Row: 1, Col: 0},
					Coord{Row: 1, Col: 1},
					Coord{Row: 1, Col: 2},
				},
				Coords{
					Coord{Row: 2, Col: 0},
					Coord{Row: 2, Col: 1},
					Coord{Row: 2, Col: 2},
				},
				Coords{
					Coord{Row: 3, Col: 0},
					Coord{Row: 3, Col: 1},
					Coord{Row: 3, Col: 2},
				},
				// cols
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 1, Col: 0},
					Coord{Row: 2, Col: 0},
				},
				Coords{
					Coord{Row: 1, Col: 0},
					Coord{Row: 2, Col: 0},
					Coord{Row: 3, Col: 0},
				},
				Coords{
					Coord{Row: 0, Col: 1},
					Coord{Row: 1, Col: 1},
					Coord{Row: 2, Col: 1},
				},
				Coords{
					Coord{Row: 1, Col: 1},
					Coord{Row: 2, Col: 1},
					Coord{Row: 3, Col: 1},
				},
				Coords{
					Coord{Row: 0, Col: 2},
					Coord{Row: 1, Col: 2},
					Coord{Row: 2, Col: 2},
				},
				Coords{
					Coord{Row: 1, Col: 2},
					Coord{Row: 2, Col: 2},
					Coord{Row: 3, Col: 2},
				},

				// diags
				Coords{
					Coord{Row: 0, Col: 0},
					Coord{Row: 1, Col: 1},
					Coord{Row: 2, Col: 2},
				},
				Coords{
					Coord{Row: 1, Col: 0},
					Coord{Row: 2, Col: 1},
					Coord{Row: 3, Col: 2},
				},
				Coords{
					Coord{Row: 2, Col: 0},
					Coord{Row: 1, Col: 1},
					Coord{Row: 0, Col: 2},
				},
				Coords{
					Coord{Row: 3, Col: 0},
					Coord{Row: 2, Col: 1},
					Coord{Row: 1, Col: 2},
				},
			},
		},
		// TODO(rsned): Other cases to test the generation?
	}

	for _, test := range tests {
		// To ensure the cmp.Equal works, apply the sorting to each set of
		// coordinates and the list of coordinates.
		for _, c := range test.want {
			slices.SortFunc(c, coordCompare)
		}
		slices.SortFunc(test.want, coordSliceCompare)

		got := test.board.generateAllWinningCoordinateSets()
		if !cmp.Equal(got, test.want, cmpopts.EquateEmpty()) {
			t.Errorf("generateAllWinningCoordinateSets(%d, %d, %d) = %+v, want: %+v\ndiff: %+v",
				test.board.rows, test.board.cols, test.board.targetSize,
				got, test.want, cmp.Diff(test.want, got))

		}
	}
}

// boardShowCoords is a helper to generate a blank board, then populate the given
// coordinates and print the resulting board for some simple visual debugging.
func boardShowCoords(rows, cols int, coords []Coord) {
	board := &Board{
		rows:  rows,
		cols:  cols,
		cells: make([][]Marker, rows),
	}

	for i := 0; i < rows; i++ {
		board.cells[i] = make([]Marker, cols)
		for j := 0; j < cols; j++ {
			board.cells[i][j] = MarkerEmpty
		}
	}

	for _, c := range coords {
		board.cells[c.Row][c.Col] = MarkerX
	}

	fmt.Printf("%v\n", board.String())
}

func TestBoardCheckOutcome(t *testing.T) {
	tests := []struct {
		board  *Board
		player *Player
		coords Coords
		want   bool
	}{
		{
			//  Coords to short, so should fail the checkOutcome.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerWhiteStone, MarkerX, MarkerWhiteStone},
					[]Marker{MarkerWhiteStone, MarkerX, MarkerX},
					[]Marker{MarkerX, MarkerWhiteStone, MarkerWhiteStone},
				},
			},
			coords: Coords{
				Coord{Row: 0, Col: 0},
				Coord{Row: 1, Col: 1},
			},
			player: Player1,
			want:   false,
		},
		{
			// Board with players in a draw. Testing vs a diagonal.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerWhiteStone, MarkerX, MarkerWhiteStone},
					[]Marker{MarkerWhiteStone, MarkerX, MarkerX},
					[]Marker{MarkerX, MarkerWhiteStone, MarkerWhiteStone},
				},
			},
			coords: Coords{
				Coord{Row: 0, Col: 0},
				Coord{Row: 1, Col: 1},
				Coord{Row: 2, Col: 2},
			},
			player: Player1,
			want:   false,
		},
		{
			//  Player 2 has a diagonal win.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerWhiteStone, MarkerX, MarkerWhiteStone},
					[]Marker{MarkerWhiteStone, MarkerWhiteStone, MarkerX},
					[]Marker{MarkerX, MarkerX, MarkerWhiteStone},
				},
			},
			coords: Coords{
				Coord{Row: 0, Col: 0},
				Coord{Row: 1, Col: 1},
				Coord{Row: 2, Col: 2},
			},
			player: Player2,
			want:   true,
		},
		{
			//  Player 1 has a horizontal win.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerX, MarkerX, MarkerX},
					[]Marker{MarkerWhiteStone, MarkerX, MarkerX},
					[]Marker{MarkerX, MarkerWhiteStone, MarkerWhiteStone},
				},
			},
			coords: Coords{
				Coord{Row: 0, Col: 0},
				Coord{Row: 0, Col: 1},
				Coord{Row: 0, Col: 2},
			},
			player: Player1,
			want:   true,
		},
	}

	for _, test := range tests {
		if got := test.board.checkOutcome(test.coords, test.player); got != test.want {
			t.Errorf("checkOutcome(%+v, %+v) = %v, want %v",
				test.coords, test.player, got, test.want)
		}
	}

}

func TestBoardOutcome(t *testing.T) {
	tests := []struct {
		board     *Board
		p1Outcome Outcome
		p2Outcome Outcome
	}{
		{
			// Empty board, should be incomplete for both.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 2,
				cells: [][]Marker{
					[]Marker{MarkerEmpty, MarkerEmpty},
					[]Marker{MarkerEmpty, MarkerEmpty},
				},
			},
			p1Outcome: OutcomeIncomplete,
			p2Outcome: OutcomeIncomplete,
		},
		{
			// Board with player1 X markers in a winning state.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 2,
				cells: [][]Marker{
					[]Marker{MarkerX, MarkerEmpty},
					[]Marker{MarkerWhiteStone, MarkerX},
				},
			},
			p1Outcome: OutcomeWin,
			p2Outcome: OutcomeLoss,
		},
		{
			// Board with player2 O markers in a winning state.
			board: &Board{
				rows:       2,
				cols:       2,
				targetSize: 2,
				cells: [][]Marker{
					[]Marker{MarkerWhiteStone, MarkerWhiteStone},
					[]Marker{MarkerEmpty, MarkerX},
				},
			},
			p1Outcome: OutcomeLoss,
			p2Outcome: OutcomeWin,
		},
		{
			// Board with players in a draw.
			board: &Board{
				rows:       3,
				cols:       3,
				targetSize: 3,
				cells: [][]Marker{
					[]Marker{MarkerWhiteStone, MarkerX, MarkerWhiteStone},
					[]Marker{MarkerWhiteStone, MarkerX, MarkerX},
					[]Marker{MarkerX, MarkerWhiteStone, MarkerWhiteStone},
				},
			},
			p1Outcome: OutcomeDraw,
			p2Outcome: OutcomeDraw,
		},
	}

	for _, test := range tests {
		test.board.winTests = test.board.generateAllWinningCoordinateSets()

		gotp1, gotp2 := test.board.Outcome()
		if gotp1 != test.p1Outcome {
			t.Errorf("Outcome() = %v, %v, want player 1 %v",
				gotp1, gotp2, test.p1Outcome)
		}
		if gotp2 != test.p2Outcome {
			t.Errorf("Outcome() = %v, %v, want player 2 %v",
				gotp1, gotp2, test.p2Outcome)
		}

	}
}

// String() isn't tested since it's just a change-detector test.
