package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuit_Left(t *testing.T) {
	checkSuit := func(s, expected Suit) func(*testing.T) {
		return func(t *testing.T) {
			require.Equal(t, expected, s.Left())
		}
	}

	t.Run("Clubs", checkSuit(SuitClubs, SuitSpades))
	t.Run("Diamonds", checkSuit(SuitDiamonds, SuitHearts))
	t.Run("Spades", checkSuit(SuitSpades, SuitClubs))
	t.Run("Hearts", checkSuit(SuitHearts, SuitDiamonds))
}

func TestCard_Score(t *testing.T) {
	checkScore := func(c Card, expected int, trump, lead Suit) func(*testing.T) {
		return func(t *testing.T) {
			require.Equal(t, expected, c.Score(trump, lead))
		}
	}

	t.Run("Trump", func(t *testing.T) {
		t.Run("Right Bower", checkScore(Card{SuitDiamonds, RankJack}, ScoreRightBower, SuitDiamonds, SuitDiamonds))
		t.Run("Left Bower", checkScore(Card{SuitDiamonds.Left(), RankJack}, ScoreLeftBower, SuitDiamonds, SuitDiamonds))
		t.Run("Ace", checkScore(Card{SuitDiamonds, RankAce}, TrumpScoreOffset+int(RankAce), SuitDiamonds, SuitDiamonds))
		t.Run("King", checkScore(Card{SuitDiamonds, RankKing}, TrumpScoreOffset+int(RankKing), SuitDiamonds, SuitDiamonds))
		t.Run("Queen", checkScore(Card{SuitDiamonds, RankQueen}, TrumpScoreOffset+int(RankQueen), SuitDiamonds, SuitDiamonds))
		t.Run("10", checkScore(Card{SuitDiamonds, RankTen}, TrumpScoreOffset+int(RankTen), SuitDiamonds, SuitDiamonds))
		t.Run("9", checkScore(Card{SuitDiamonds, RankNine}, TrumpScoreOffset+int(RankNine), SuitDiamonds, SuitDiamonds))
	})

	t.Run("Lead Suit", func(t *testing.T) {
		t.Run("Ace", checkScore(Card{SuitDiamonds, RankAce}, int(RankAce), SuitClubs, SuitDiamonds))
		t.Run("King", checkScore(Card{SuitDiamonds, RankKing}, int(RankKing), SuitClubs, SuitDiamonds))
		t.Run("Queen", checkScore(Card{SuitDiamonds, RankQueen}, int(RankQueen), SuitClubs, SuitDiamonds))
		t.Run("Jack", checkScore(Card{SuitDiamonds, RankJack}, int(RankJack), SuitClubs, SuitDiamonds))
		t.Run("10", checkScore(Card{SuitDiamonds, RankTen}, int(RankTen), SuitClubs, SuitDiamonds))
		t.Run("9", checkScore(Card{SuitDiamonds, RankNine}, int(RankNine), SuitClubs, SuitDiamonds))

		t.Run("Left Bower", checkScore(Card{SuitDiamonds, RankJack}, ScoreLeftBower, SuitHearts, SuitDiamonds))
	})

	t.Run("Other Suits", func(t *testing.T) {
		for r := RankNine; r <= RankAce; r++ {
			t.Run(fmt.Sprintf("Rank %d", r), checkScore(Card{SuitDiamonds, r}, 0, SuitClubs, SuitClubs))
		}
	})
}

func TestCard_BetterThan(t *testing.T) {
	checkBeats := func(c, other Card, trump, lead Suit) func(*testing.T) {
		return func(t *testing.T) {
			require.True(t, c.BetterThan(&other, trump, lead))
		}
	}

	checkBeatsAllInSuit := func(c Card, target, trump, lead Suit) func(*testing.T) {
		return func(t *testing.T) {
			for r := RankAce; r >= RankNine; r-- {
				t.Run(fmt.Sprintf("Beats Rank %d in Suit %d", r, target), checkBeats(c, Card{target, r}, trump, lead))
			}
		}
	}

	t.Run("Right Bower", func(t *testing.T) {
		bower := Card{SuitClubs, RankJack}

		for s := SuitDiamonds; s <= SuitSpades; s++ {
			t.Run(fmt.Sprintf("Beats everythign in suit %d", s), checkBeatsAllInSuit(bower, s, SuitClubs, SuitClubs))
		}
	})

	t.Run("Left Bower", func(t *testing.T) {
		bower := Card{SuitSpades, RankJack}

		t.Run("Beats Ace of Trump", checkBeats(bower, Card{SuitClubs, RankAce}, SuitClubs, SuitClubs))
		t.Run("Beats King of Trump", checkBeats(bower, Card{SuitClubs, RankKing}, SuitClubs, SuitClubs))
		t.Run("Beats Queen of Trump", checkBeats(bower, Card{SuitClubs, RankQueen}, SuitClubs, SuitClubs))
		t.Run("Beats Ten of Trump", checkBeats(bower, Card{SuitClubs, RankTen}, SuitClubs, SuitClubs))
		t.Run("Beats Nine of Trump", checkBeats(bower, Card{SuitClubs, RankNine}, SuitClubs, SuitClubs))

		for s := SuitDiamonds; s <= SuitSpades; s++ {
			t.Run(fmt.Sprintf("Beats everythign in suit %d", s), checkBeatsAllInSuit(bower, s, SuitClubs, SuitClubs))
		}
	})

	t.Run("Lead beats non-trump", checkBeats(Card{SuitClubs, RankAce}, Card{SuitSpades, RankAce}, SuitDiamonds, SuitClubs))
}
