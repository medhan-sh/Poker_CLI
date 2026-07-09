# Blackjack CLI

A command-line Blackjack game written in Go, built on top of the Learn Go with Tests poker project architecture.

## Features

- **Full Blackjack gameplay** — Hit, Stand, natural Blackjack detection, dealer auto-play
- **Soft ace handling** — Aces automatically convert from 11 to 1 to avoid busting
- **Chip system** — Start with 1000 chips, bet on each hand, track your winnings
- **Persistent win tracking** — Wins are recorded to a JSON file (`game.db.json`) via the existing `PlayerStore`
- **Dealer logic** — Dealer hits until 17 or higher (standard casino rules)
- **Blackjack payout** — Natural blackjack pays 3:2

## How to Play

### Run the game

```bash
go run ./cmd/cli/
```

### Commands

| Command | Description |
|---------|-------------|
| `bet <amount>` | Place your bet for the current hand |
| `hit` or `h` | Take another card |
| `stand` or `s` | End your turn and let the dealer play |
| `quit` or `q` | Exit the game (at bet prompt) |
| `y` / `n` | Play another hand or quit |

### Game Flow

1. You start with **1000 chips**
2. Each hand, enter your **bet**
3. You and the dealer are dealt **2 cards each** (one dealer card is hidden)
4. Choose to **hit** (take a card) or **stand** (stop)
5. If you bust (>21), you lose your bet
6. If you stand, the dealer reveals their hand and hits until **17 or higher**
7. Winner is determined:
   - **Blackjack** (21 with 2 cards) pays 3:2
   - Higher hand wins
   - Push (tie) returns your bet
   - Dealer bust = you win

## Project Structure

```
├── blackjack.go       # Blackjack game engine (hands, dealing, scoring)
├── card.go            # Card, Suit, Rank, Deck types
├── CLI.go             # CLI with PlayBlackjack() game loop
├── cmd/cli/main.go    # Entry point — wires up BlackjackGame
├── game.go            # Original TexasHoldEm game (unused in blackjack mode)
├── server.go          # HTTP server + PlayerStore interface
├── file_system.go     # File-based player store (JSON persistence)
├── league.go          # League/Player types
├── blind_alerter.go   # Blind alerting (unused in blackjack mode)
├── testing.go         # Test helpers and stubs
├── cli_test.go        # CLI tests
└── game.db.json       # Win records (auto-created)
```

## Architecture

The game integrates with the existing architecture from the Learn Go with Tests project:

- **`PlayerStore` interface** — Wins are recorded via `RecordWin("player")` on each winning hand
- **`Game` interface** — `BlackjackGame` implements `Start()` and `Finish()` to satisfy the interface
- **`CLI` struct** — The `PlayBlackjack()` method drives the interactive game loop, while `PlayPoker()` remains available for the original Texas Hold'em flow

## Running Tests

```bash
go test ./...
