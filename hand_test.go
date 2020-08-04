package pokerlib

import (
	"fmt"
	"testing"
)

func TestHand(t *testing.T) {

	hand := make([]Card, 5)

	fmt.Printf("TestHand\n")
	hand[0] = Card{2, Hearts}
	hand[1] = Card{3, Clubs}
	hand[2] = Card{4, Hearts}
	hand[3] = Card{5, Hearts}
	hand[4] = Card{6, Spades}

	fmt.Printf("straight = %v\n", isStraight(hand))
	fmt.Printf("flush = %v\n", isFlush(hand))

	hand[0] = Card{2, Hearts}
	hand[1] = Card{3, Hearts}
	hand[2] = Card{4, Hearts}
	hand[3] = Card{5, Hearts}
	hand[4] = Card{7, Hearts}

	fmt.Printf("straight = %v\n", isStraight(hand))
	fmt.Printf("flush = %v\n", isFlush(hand))

	v := Rank(hand)
	hr := GetHandKind(v)
	s := hr.String()

	fmt.Printf("h1 value: %d\n", v)
	fmt.Printf("rank = %b\n", v)
	fmt.Printf("hr = %s\n", s)

	hand[4] = Card{8, Hearts}
	v = Rank(hand)
	hr = GetHandKind(v)
	fmt.Printf("h2 value: %d\n", v)
	s = hr.String()

}

func TestSuitsEqual(t *testing.T) {

	hand := make([]Card, 5)

	fmt.Printf("TestSuitsEqual\n")
	hand[0] = Card{2, Hearts}
	hand[1] = Card{3, Hearts}
	hand[2] = Card{4, Hearts}
	hand[3] = Card{5, Hearts}
	hand[4] = Card{6, Hearts}

	v1 := Rank(hand)

	hand[0] = Card{2, Clubs}
	hand[1] = Card{3, Clubs}
	hand[2] = Card{4, Clubs}
	hand[3] = Card{5, Clubs}
	hand[4] = Card{6, Clubs}

	v2 := Rank(hand)

	if v2 != v1 {
		fmt.Printf("Problem: Different suited identical hand have different rank\n")
	}
}

func TestFullHands(t *testing.T) {

	fmt.Printf("FullHandTest\n")

	highestPairHand := Hand([]Card{
		Card{Ace, Hearts},
		Card{Ace, Clubs},
		Card{King, Diamonds},
		Card{Queen, Clubs},
		Card{Ten, Hearts},
	})

	lowest2PairHand := Hand([]Card{
		Card{Two, Hearts},
		Card{Two, Clubs},
		Card{Three, Diamonds},
		Card{Three, Clubs},
		Card{Four, Hearts},
	})

	highest2PairHand := Hand([]Card{
		Card{Ace, Hearts},
		Card{Ace, Clubs},
		Card{King, Diamonds},
		Card{King, Clubs},
		Card{Queen, Hearts},
	})

	lowestSetHand := Hand([]Card{
		Card{Two, Hearts},
		Card{Two, Clubs},
		Card{Two, Diamonds},
		Card{Three, Clubs},
		Card{Four, Hearts},
	})

	highestSetHand := Hand([]Card{
		Card{Ace, Hearts},
		Card{Ace, Clubs},
		Card{Ace, Diamonds},
		Card{King, Clubs},
		Card{Queen, Hearts},
	})

	lowestStraight := Hand([]Card{
		Card{Ace, Hearts},
		Card{Two, Clubs},
		Card{Three, Diamonds},
		Card{Four, Clubs},
		Card{Five, Hearts},
	})

	highestStraight := Hand([]Card{
		Card{Ten, Hearts},
		Card{Jack, Clubs},
		Card{Queen, Diamonds},
		Card{King, Clubs},
		Card{Ace, Hearts},
	})

	lowestFlush := Hand([]Card{
		Card{Two, Hearts},
		Card{Three, Hearts},
		Card{Four, Hearts},
		Card{Five, Hearts},
		Card{Seven, Hearts},
	})

	highestFlush := Hand([]Card{
		Card{Ace, Hearts},
		Card{King, Hearts},
		Card{Queen, Hearts},
		Card{Jack, Hearts},
		Card{Nine, Hearts},
	})
	_ = highestFlush

	AssertGreater(lowest2PairHand, highestPairHand)
	AssertGreater(lowestSetHand, highest2PairHand)
	AssertGreater(lowestStraight, highestSetHand)
	AssertGreater(lowestFlush, highestStraight)

}

func AssertGreater(h1 Hand, h2 Hand) {
	r1 := Rank(h1)
	r2 := Rank(h2)

	fmt.Printf("val: %.24b, for hand %v\n",  r1, h1)
	fmt.Printf("val: %.24b, for hand %v\n",  r2, h2)
	if r1 <= r2 {
		fmt.Printf("Error hand %v greater than %v\n", h1, h2)
	}
}
