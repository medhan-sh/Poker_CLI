package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	if err != nil {
		fmt.Fprint(cli.out, "you're so silly!")
		return
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

// PlayBlackjack runs the blackjack game loop
func (cli *CLI) PlayBlackjack() {
	bjGame, ok := cli.game.(*BlackjackGame)
	if !ok {
		fmt.Fprint(cli.out, "Error: game is not a blackjack game\n")
		return
	}

	fmt.Fprintf(cli.out, "Welcome to Blackjack!\n")
	fmt.Fprintf(cli.out, "You have %d chips.\n", bjGame.Chips())
	fmt.Fprintf(cli.out, "Commands: bet <amount>, hit, stand\n")

	for bjGame.Chips() > 0 {
		fmt.Fprintf(cli.out, "\n--- New Hand ---\n")
		fmt.Fprintf(cli.out, "Chips: %d\n", bjGame.Chips())

		// Get bet
		bet := cli.getBet(bjGame)
		if bet < 0 {
			// Player wants to quit
			fmt.Fprintf(cli.out, "Thanks for playing! You leave with %d chips.\n", bjGame.Chips())
			return
		}

		// Play the round
		result, err := bjGame.PlayRound(bet)
		if err != nil {
			fmt.Fprintf(cli.out, "Error: %s\n", err)
			continue
		}

		// If result is empty, the hand is still in progress (player needs to hit/stand)
		if result == "" {
			result = cli.playHand(bjGame)
		}

		// Record win if player won
		if result == "win" || result == "player_blackjack" {
			bjGame.Finish("player")
		}

		// Ask to play again
		if bjGame.Chips() <= 0 {
			fmt.Fprintf(cli.out, "You're out of chips! Game over.\n")
			return
		}

		fmt.Fprintf(cli.out, "Play another hand? (y/n): ")
		answer := strings.TrimSpace(cli.readLine())
		if answer != "y" && answer != "yes" {
			fmt.Fprintf(cli.out, "Thanks for playing! You leave with %d chips.\n", bjGame.Chips())
			return
		}
	}
}

func (cli *CLI) getBet(bjGame *BlackjackGame) int {
	for {
		fmt.Fprintf(cli.out, "Enter bet (or 'quit' to exit): ")
		input := strings.TrimSpace(cli.readLine())

		if input == "quit" || input == "q" {
			return -1
		}

		bet, err := strconv.Atoi(input)
		if err != nil || bet <= 0 || bet > bjGame.Chips() {
			fmt.Fprintf(cli.out, "Invalid bet. You have %d chips. Try again.\n", bjGame.Chips())
			continue
		}

		return bet
	}
}

func (cli *CLI) playHand(bjGame *BlackjackGame) string {
	for {
		fmt.Fprintf(cli.out, "Your action (hit/stand): ")
		input := strings.TrimSpace(strings.ToLower(cli.readLine()))

		switch input {
		case "hit", "h":
			result, err := bjGame.Hit()
			if err != nil {
				fmt.Fprintf(cli.out, "Error: %s\n", err)
				continue
			}
			if result != "" {
				return result // bust, blackjack, etc.
			}
			// Continue - player can hit again or stand
		case "stand", "s":
			result, err := bjGame.Stand()
			if err != nil {
				fmt.Fprintf(cli.out, "Error: %s\n", err)
				continue
			}
			return result
		default:
			fmt.Fprintf(cli.out, "Unknown command. Use 'hit' or 'stand'.\n")
		}
	}
}
