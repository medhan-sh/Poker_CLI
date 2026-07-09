// cmd/cli/main.go
package main

import (
	"fmt"
	poker "http_server"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play Blackjack!")
	game := poker.NewBlackjackGame(store, os.Stdout)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayBlackjack()
}
