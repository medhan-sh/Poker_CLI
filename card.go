package poker

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}

type Rank int

func (r Rank) String() string {
	switch r {
	case 14:
		return "A"
	case 13:
		return "K"
	case 12:
		return "Q"
	case 11:
		return "J"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

func (c Card) BlackjackValue() int {
	switch c.Rank {
	case 14: // Ace
		return 11
	case 13: // King
		return 10
	case 12: // Queen
		return 10
	case 11: // Jack
		return 10
	default:
		return int(c.Rank)
	}
}

type Deck struct {
	cards []Card
	index int
}

func NewDeck() *Deck {
	cards := make([]Card, 0, 52)
	for suit := Hearts; suit <= Spades; suit++ {
		for rank := 2; rank <= 14; rank++ {
			cards = append(cards, Card{Rank: Rank(rank), Suit: suit})
		}
	}

	d := &Deck{cards: cards, index: 0}
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
	d.index = 0
}

func (d *Deck) Deal() Card {
	if d.index >= len(d.cards) {
		d.Shuffle()
	}
	card := d.cards[d.index]
	d.index++
	return card
}

// Remaining returns the number of cards left
func (d *Deck) Remaining() int {
	return len(d.cards) - d.index
}
