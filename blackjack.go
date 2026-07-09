package poker

import (
	"fmt"
	"io"
)

type BlackjackHand struct {
	Cards []Card
}

// AddCard adds a card to the hd
func (h *BlackjackHand) AddCard(c Card) {
	h.Cards = append(h.Cards, c)
}

// Value returns the best blackjack value of the hand (handles soft aces)
func (h *BlackjackHand) Value() int {
	total := 0
	aces := 0

	for _, card := range h.Cards {
		total += card.BlackjackValue()
		if card.Rank == 14 { // Ace
			aces++
		}
	}

	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}

	return total
}

// IsBust returns true if the hand value exceeds 21
func (h *BlackjackHand) IsBust() bool {
	return h.Value() > 21
}

func (h *BlackjackHand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

func (h *BlackjackHand) String() string {
	result := ""
	for i, card := range h.Cards {
		if i > 0 {
			result += " "
		}
		result += card.String()
	}
	result += fmt.Sprintf(" (%d)", h.Value())
	return result
}

// BlackjackGame represents a game of blackjack
type BlackjackGame struct {
	store  PlayerStore
	out    io.Writer
	deck   *Deck
	player *BlackjackHand
	dealer *BlackjackHand
	chips  int
	bet    int
}

func NewBlackjackGame(store PlayerStore, out io.Writer) *BlackjackGame {
	return &BlackjackGame{
		store: store,
		out:   out,
		chips: 1000, // Starting chips
	}
}

func (g *BlackjackGame) Start(numberOfPlayers int) {
	// Not used in blackjack - we use PlayRound instead
}

func (g *BlackjackGame) Finish(winner string) {
	g.store.RecordWin(winner)
}

// Chips returns the player's current chip count
func (g *BlackjackGame) Chips() int {
	return g.chips
}

func (g *BlackjackGame) PlayRound(bet int) (string, error) {
	if bet <= 0 || bet > g.chips {
		return "", fmt.Errorf("invalid bet: have %d chips, bet %d", g.chips, bet)
	}

	g.bet = bet
	g.deck = NewDeck()
	g.player = &BlackjackHand{}
	g.dealer = &BlackjackHand{}

	// Deal initial cards
	g.player.AddCard(g.deck.Deal())
	g.dealer.AddCard(g.deck.Deal())
	g.player.AddCard(g.deck.Deal())
	g.dealer.AddCard(g.deck.Deal())

	// Show initial state
	fmt.Fprintf(g.out, "\nYour hand: %s\n", g.player)
	fmt.Fprintf(g.out, "Dealer shows: %s ?\n", g.dealer.Cards[0])

	// Check for dealer blackjack
	if g.dealer.IsBlackjack() {
		fmt.Fprintf(g.out, "Dealer hand: %s\n", g.dealer)
		if g.player.IsBlackjack() {
			g.chips += bet // Push - return bet
			fmt.Fprintf(g.out, "Push! Both have blackjack.\n")
			return "push", nil
		}
		g.chips -= bet
		fmt.Fprintf(g.out, "Dealer has blackjack! You lose %d chips.\n", bet)
		return "dealer_blackjack", nil
	}

	if g.player.IsBlackjack() {
		winnings := bet + bet/2
		g.chips += winnings
		fmt.Fprintf(g.out, "Blackjack! You win %d chips!\n", winnings)
		return "player_blackjack", nil
	}

	return "", nil
}

func (g *BlackjackGame) Hit() (string, error) {
	if g.player == nil {
		return "", fmt.Errorf("game not started")
	}

	g.player.AddCard(g.deck.Deal())
	fmt.Fprintf(g.out, "Your hand: %s\n", g.player)

	if g.player.IsBust() {
		g.chips -= g.bet
		fmt.Fprintf(g.out, "Bust! You lose %d chips.\n", g.bet)
		return "bust", nil
	}

	return "", nil
}

func (g *BlackjackGame) Stand() (string, error) {
	if g.player == nil {
		return "", fmt.Errorf("game not started")
	}

	// Reveal dealer's hidden card
	fmt.Fprintf(g.out, "Dealer hand: %s\n", g.dealer)

	// Dealer hits until 17 or higher
	for g.dealer.Value() < 17 {
		g.dealer.AddCard(g.deck.Deal())
		fmt.Fprintf(g.out, "Dealer hits: %s\n", g.dealer)
	}

	// Determine winner
	playerValue := g.player.Value()
	dealerValue := g.dealer.Value()

	if g.dealer.IsBust() {
		g.chips += g.bet
		fmt.Fprintf(g.out, "Dealer busts! You win %d chips!\n", g.bet)
		return "win", nil
	}

	if playerValue > dealerValue {
		g.chips += g.bet
		fmt.Fprintf(g.out, "You win %d chips!\n", g.bet)
		return "win", nil
	}

	if playerValue == dealerValue {
		fmt.Fprintf(g.out, "Push! Your bet is returned.\n")
		return "push", nil
	}

	g.chips -= g.bet
	fmt.Fprintf(g.out, "Dealer wins! You lose %d chips.\n", g.bet)
	return "lose", nil
}
