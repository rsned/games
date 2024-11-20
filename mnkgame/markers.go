package mnkgame

// Marker represents the various pre-defined markers that may appear in the games.
type Marker string

const (
	filledBlackCircle = "⚫" // U+26AB MEDIUM BLACK CIRCLE
	filledWhiteCircle = "⭘" // U+2B58 HEAVY CIRCLE
	blackX            = "🗙" // U+1F5D9 CANCELLATION X

	winMarkerUpArrowWhite    = "▵" // U+25B5 - WHITE UP-POINTING SMALL TRIANGLE
	winMarkerDownArrowWhite  = "▿" // U+25BF - WHITE DOWN-POINTING SMALL TRIANGLE
	winMarkerRightArrowWhite = "▹" // U+25B9 - WHITE RIGHT-POINTING SMALL TRIANGLE
	winMarkerLeftArrowWhite  = "◃" // U+25C3 - WHITE LEFT-POINTING SMALL TRIANGLE

	winMarkerUpArrowBlack    = "▴" // U+25B4 - BLACK UP-POINTING SMALL TRIANGLE
	winMarkerDownArrowBlack  = "▾" // U+25BE - BLACK DOWN-POINTING SMALL TRIANGLE
	winMarkerRightArrowBlack = "▸" // U+25B8 - BLACK RIGHT-POINTING SMALL TRIANGLE
	winMarkerLeftArrowBlack  = "◂" // U+25C2 - BLACK LEFT-POINTING SMALL TRIANGLE
)

// Predefine some markers.
//
// TODO(rsned): Add more markers to choose from.
var (
	MarkerEmpty      = Marker(" ")
	MarkerX          = Marker(blackX)
	MarkerWhiteStone = Marker(filledWhiteCircle)
	MarkerBlackStone = Marker(filledBlackCircle)
)
