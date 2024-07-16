package main

import (
	"fmt"
	"math/rand"

	"github.com/rsned/games/mnkgame"
)

func main() {
	playerN := readInput("Do you wish to be player 1 or 2?", []string{"1", "2"})
	var player1Play, player2Play playFunc
	player1 := mnkgame.Player1
	player2 := mnkgame.Player2

	if playerN == "1" {
		player1.SetHuman()
		player1Play = humanPlayer
		player2.SetComputer()
		player2Play = randomPlayer
	} else {
		player1.SetComputer()
		player1Play = randomPlayer
		player2.SetHuman()
		player2Play = humanPlayer
	}

	game := mnkgame.TicTacToe(player1, player2)
	var move string

	// TODO(rsned): There's lots of repetition here, refactor player some more
	// and the Outcome method to make it easier for player to be more abstract
	// and this loop and above setup simpler.
	for {
		// Player 1
		fmt.Printf("\n%s\n", game.RenderBoard())
		move = player1Play(player1, game)
		game.ApplyMove(player1, move)

		if p1, _ := game.Outcome(); p1 != mnkgame.OutcomeIncomplete {
			fmt.Printf("\n%s\n", game.RenderBoard())
			fmt.Printf("Game Over. %s Wins.\n", player1.String())
			break
		}

		// Player 2
		fmt.Printf("\n%s\n", game.RenderBoard())
		move = player2Play(player2, game)
		game.ApplyMove(player2, move)

		if _, p2 := game.Outcome(); p2 != mnkgame.OutcomeIncomplete {
			fmt.Printf("\n%s\n", game.RenderBoard())
			fmt.Printf("Game Over. %s Wins.\n", player2)
			break
		}
	}
}

func readInput(prompt string, valid []string) string {
	var entry string
	vals := map[string]bool{}

	for _, v := range valid {
		vals[v] = true
	}

	for {
		fmt.Println(prompt)
		_, err := fmt.Scanf("%s", &entry)
		if vals[entry] {
			return entry
		}

		if err != nil {
			fmt.Print("Error reading entry please try again.\n", err)
			continue
		}

		fmt.Printf("Invalid entry %q, please try again.\n", entry)
	}
}

type playFunc func(*mnkgame.Player, *mnkgame.MNKGame) string

func humanPlayer(player *mnkgame.Player, game *mnkgame.MNKGame) string {
	moves := game.PotentialMoves()
	move := readInput(fmt.Sprintf("Select square: %+v", moves), moves)
	return move

}

func randomPlayer(player *mnkgame.Player, games *mnkgame.MNKGame) string {
	moves := games.PotentialMoves()
	move := moves[rand.Intn(len(moves))]
	fmt.Printf("%s plays %s\n", player, move)
	return move
}
