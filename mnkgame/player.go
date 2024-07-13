package mnkgame

type playerType int

const (
	playerTypeHuman playerType = iota
	playerTypeComputerRandom
	playerTypeComputerAI
)

// Player holds basic fields about a player in this game, primarily what
// type of player and what marker it uses.
type Player struct {
	id          string
	displayName string

	playerType playerType

	marker Marker
}

// SetHuman updates the player type to be a human.
func (p *Player) SetHuman() {
	p.playerType = playerTypeHuman
}

// SetComputer sets the player type to be a computer making random moves.
func (p *Player) SetComputer() {
	p.playerType = playerTypeComputerRandom
}

func (p *Player) String() string {
	return p.displayName
}

// Predefine some players that can be used in games.
var (
	Player1 = &Player{
		id:          "1",
		displayName: "Player 1",
		marker:      MarkerX,
		playerType:  playerTypeHuman,
	}

	Player2 = &Player{
		id:          "2",
		displayName: "Player 2",
		marker:      MarkerWhiteStone,
		playerType:  playerTypeHuman,
	}

	PlayerComputer1 = &Player{
		id:          "1001",
		displayName: "Computer Player Player 1",
		marker:      MarkerWhiteStone,
		playerType:  playerTypeComputerRandom,
	}

	PlayerComputer2 = &Player{
		id:          "1002",
		displayName: "Computer Player Player 2",
		marker:      MarkerBlackStone,
		playerType:  playerTypeComputerRandom,
	}
)

// TODO(rsned): Add in some mechanism for the player to choose its move, be it a
// human reading from STDIN, or a computer player uise rand.Intn(), etc.
