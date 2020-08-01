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
	hr := GetHandRank(v)
	s := hr.String()

	fmt.Printf("h1 value: %d\n", v)
	fmt.Printf("rank = %b\n", v)
	fmt.Printf("hr = %s\n", s)

	hand[4] = Card{8, Hearts}
	v = Rank(hand)
	hr = GetHandRank(v)
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
