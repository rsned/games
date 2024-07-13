package mnkgame

// Marker represents the various pre-defined markers that may appear in the games.
type Marker string

const (
	filledBlackCircle = "⬤" // U+2B24 Black Large Circle
	filledWhiteCircle = "◯" // U+25EF Large Circle
	blackX            = "🗙" // U+1F5D9 Cancellation X
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
