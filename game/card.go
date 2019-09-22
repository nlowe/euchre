package game

import "fmt"

// Suit indicates a given suit in a Euchre Deck
type Suit int

const (
	SuitMask = 0b11

	SuitClubs    Suit = 0b00
	SuitDiamonds Suit = 0b01
	SuitSpades   Suit = 0b11
	SuitHearts   Suit = 0b10
)

func (s Suit) String() string {
	switch s {
	case SuitClubs:
		return "Clubs"
	case SuitDiamonds:
		return "Diamonds"
	case SuitSpades:
		return "Spades"
	case SuitHearts:
		return "Hearts"
	}

	return "UNKNOWN"
}

// Left returns the suit of the opposite color for finding the Left Bower
func (s Suit) Left() Suit {
	return ^s & SuitMask
}

// Rank is the rank of the card in a given Suit. Do not use the Rank for scoring,
// instead, call card.Score() with the trump and lead suit
type Rank int

const (
	RankNine Rank = 9 + iota
	RankTen
	RankJack
	RankQueen
	RankKing
	RankAce
)

func (r Rank) String() string {
	switch r {
	case RankNine:
		return "9"
	case RankTen:
		return "10"
	case RankJack:
		return "Jack"
	case RankQueen:
		return "Queen"
	case RankKing:
		return "King"
	case RankAce:
		return "Ace"
	}

	return "UNKNOWN"
}

const (
	TrumpScoreOffset = int(RankAce) + 1

	ScoreLeftBower = TrumpScoreOffset + int(RankAce) + 1 + iota
	ScoreRightBower
)

// Card represents a card in a Euchre deck
type Card struct {
	Suit Suit
	Rank Rank
}

func (c *Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())
}

// Score returns the score of the given card with the provided trump and lead suits
func (c *Card) Score(trump, lead Suit) int {
	switch c.Suit {
	case trump:
		if c.Rank == RankJack {
			return ScoreRightBower
		}

		return int(c.Rank) + TrumpScoreOffset
	case trump.Left():
		if c.Rank == RankJack {
			return ScoreLeftBower
		}

		// Fallthrough in case the left suit was also lead
		fallthrough
	case lead:
		// If this isn't trump or the left bower, use the Rank for the score
		return int(c.Rank)
	default:
		// Cards that aren't trump and aren't what was lead are worthless
		return 0
	}
}

// BetterThan returns true iff this card is better than the provided card, given
// the specified trump and lead suits
func (c *Card) BetterThan(other *Card, trump, lead Suit) bool {
	// Use >= here instead of > to make tests easier to write
	// Realistically speaking we will never be comparing a card to itself
	return c.Score(trump, lead) >= other.Score(trump, lead)
}
