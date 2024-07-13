package mnkgame

import (
	"fmt"
	"slices"
)

// Coord is the coordinates of a cell in a game using a Top-Left origin.
type Coord struct {
	Row int
	Col int

	// TODO(rsned): Consiser adding a connections to other Coord for
	// values that have specific ways they can be reached from other
	// Coords.
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.Row, c.Col)
}

// less is used for sorting.
func (c Coord) less(other Coord) bool {
	if c.Row <= other.Row {
		return c.Col < other.Col
	}
	return false
}

// equals is used in comparisons.
func (c Coord) equals(other Coord) bool {
	return c.Row == other.Row && c.Col == other.Col
}

// coordEqual is the non type-bound check for equality.
func coordEqual(a, b Coord) bool {
	return a.equals(b)
}

// coordCompare is the cmp.Compare func for Coord types.
func coordCompare(a, b Coord) int {
	if a.Row < b.Row {
		return -1
	}
	if a.Row > b.Row {
		return 1
	}

	// Rows are the same so go by Col.
	if a.Col < b.Col {
		return -1
	}
	if a.Col == b.Col {
		return 0
	}

	return 1
}

// Coords is a slice of Coord values.
type Coords []Coord

// Add attempts to add the given value, skipping if the value already is in the slice.
func (c *Coords) Add(coord Coord) {
	for _, v := range *c {
		if v.equals(coord) {
			fmt.Printf("Can't add %+v to %+v, matches existing\n", v, c)
			return
		}
	}
	(*c) = append(*c, coord)
}

// coordSliceCompare compares the two slices for equality or ordering.
func coordSliceCompare(a, b Coords) int {
	return slices.CompareFunc(a, b, coordCompare)
}

// CoordsList is a list of list of Coord.
type CoordsList []Coords

// Add attempts to add the given value, skipping if the value already is in this.
//
// This method assuemes the input is in the same order as existing values.
// Permutations are NOT checked and are considered distinct.
func (c *CoordsList) Add(coords Coords) {
	for _, v := range *c {
		// TODO(rsned): This is assuming the slices are already sorted.
		if slices.EqualFunc(v, coords, coordEqual) {
			return
		}
	}
	(*c) = append(*c, coords)
}
