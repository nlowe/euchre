package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShuffleDeck(t *testing.T) {
	deck := ShuffleDeck()

	t.Run("Contains 24 Cards", func(t *testing.T) {
		require.Len(t, deck, 24)
	})

	t.Run("Contains No Duplicates", func(t *testing.T) {
		for i := 0; i < len(deck)-1; i++ {
			for j := i + 1; j < len(deck); j++ {
				assert.False(
					t,
					deck[i].Suit == deck[j].Suit && deck[i].Rank == deck[j].Rank,
					fmt.Sprintf(
						"Expected no duplicates but found a duplicate %d=%s %d=(%s)",
						i, deck[i].String(),
						j, deck[j].String(),
					),
				)
			}
		}
	})
}
