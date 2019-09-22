package game

import "math/rand"

func ShuffleDeck() []Card {
	var deck []Card

	for s := SuitClubs; s <= SuitSpades; s++ {
		for r := RankNine; r <= RankAce; r++ {
			deck = append(deck, Card{s, r})
		}
	}

	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}
