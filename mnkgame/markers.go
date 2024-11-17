package mnkgame

// Marker represents the various pre-defined markers that may appear in the games.
type Marker string

const (
	filledBlackCircle = "⬤" // U+2B24 BLACK LARGE CIRCLE
	//filledBlackCircle = "⚫" // U+26AB MEDIUM BLACK CIRCLE
	filledWhiteCircle = "⭘" // U+2B58 HEAVY CIRCLE
	blackX            = "🗙" // U+1F5D9 CANCELLATION X

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
