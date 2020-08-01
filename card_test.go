package pokerlib

import (
	"fmt"
	"testing"
)

func TestCard(t *testing.T) {

	c1 := Card{2, Spades}
	c2 := Card{Ace, Clubs}

	fmt.Printf("2 spades  = %d\n", int(c1.AbsoluteValue()))
	fmt.Printf("Ace clubs = %d\n", int(c2.AbsoluteValue()))

}

func TestAbsoluteValues(t *testing.T) {

	hand1 := make([]Card, 5)

	hand1[0] = Card{2, Hearts}
	hand1[1] = Card{3, Clubs}
	hand1[2] = Card{4, Hearts}
	hand1[3] = Card{5, Hearts}
	hand1[4] = Card{6, Spades}

	hand2 := make([]Card, 5)

	hand2[0] = Card{2, Hearts}
	hand2[1] = Card{12, Clubs}
	hand2[4] = Card{4, Hearts}
	hand2[3] = Card{5, Hearts}
	hand2[2] = Card{6, Spades}

	hand3 := make([]Card, 5)

	hand3[0] = Card{2, Hearts}
	hand3[1] = Card{12, Clubs}
	hand3[3] = Card{4, Hearts}
	hand3[2] = Card{5, Hearts}
	hand3[4] = Card{6, Spades}

	h1Rank := Rank(hand1)
	h2Rank := Rank(hand2)
	h3Rank := Rank(hand3)

	if h1Rank == h2Rank {
		t.Error("hands should not have same rank\n")
	}
	if h2Rank != h3Rank {
		t.Error("h2 and h3 should have same rank\n")
	}
	fmt.Printf("h1 abs value = %d\n", h1Rank)
	fmt.Printf("h2 abs value = %d\n", h2Rank)

}

func TestDeck(t *testing.T) {

	hand1 := make([]Card, 7)

	hand1[0] = Card{2, Hearts}
	hand1[1] = Card{3, Clubs}
	hand1[2] = Card{4, Hearts}
	hand1[3] = Card{5, Hearts}
	hand1[4] = Card{6, Spades}
	hand1[5] = Card{8, Hearts}
	hand1[6] = Card{14, Hearts}

	var h1 Hand = hand1

	hand2 := make([]Card, 7)

	hand2[0] = Card{2, Hearts}
	hand2[1] = Card{12, Clubs}
	hand2[2] = Card{4, Hearts}
	hand2[3] = Card{5, Hearts}
	hand2[4] = Card{6, Spades}
	hand2[5] = Card{8, Hearts}
	hand2[6] = Card{14, Hearts}

	var h2 Hand = hand2

	if h1.Equals(h2) {
		t.Error("h1 and h2 are not equal\n")
	} else {
		fmt.Printf("hand equals worked\n")
	}

	if h1.ContainsCard(Card{8, Clubs}) {
		t.Error("h1 does not contain an 8 of clubs\n")
	} else {
		fmt.Printf("negative hand ContainsCard worked\n")
	}
	if !h2.ContainsCard(Card{8, Hearts}) {
		t.Error("h2 does contain an 8 of hearts\n")
	} else {
		fmt.Printf("hand ContainsCard worked\n")
	}
	/*

		deck := NewDeck()
		fmt.Printf("%s", deck)
		deck.Shuffle()
		fmt.Printf("%s", deck)

		deck.DrawCard()
		deck.DrawCard()
		deck.DrawCard()
		deck.DrawCard()
		deck.DrawCard()



		fmt.Printf("%s", deck)
	*/

}
