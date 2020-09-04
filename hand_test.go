package pokerlib

import (
	"fmt"
	"strings"
	"testing"
)

func TestSort(t *testing.T) {
	hand := make([]Card, 5)

	fmt.Printf("TestHand\n")
	hand[0] = Card{7, Hearts}
	hand[1] = Card{4, Clubs}
	hand[2] = Card{11, Hearts}
	hand[3] = Card{14, Hearts}
	hand[4] = Card{2, Spades}

	SortCards(hand)

	if hand[0].Index != Two ||
		hand[1].Index != Four ||
		hand[2].Index != Seven ||
		hand[3].Index != Jack ||
		hand[4].Index != Ace {
		t.Errorf("SortCards not working as expected")
	}

}
func TestHand(t *testing.T) {

	hand := make([]Card, 5)

	fmt.Printf("TestHand\n")
	hand[0] = Card{2, Hearts}
	hand[1] = Card{3, Clubs}
	hand[2] = Card{4, Hearts}
	hand[3] = Card{5, Hearts}
	hand[4] = Card{6, Spades}

	if isStraight(hand) == false {
		t.Error("Straight not recognized\n")
	}
	if isFlush(hand) == true {
		for i, c := range hand {
			fmt.Printf("card %d is a %s\n", i, c.Suit.String())

		}
		t.Error("Flush incorrectly recognized\n")
	}

	hand[0] = Card{2, Hearts}
	hand[1] = Card{3, Hearts}
	hand[2] = Card{4, Hearts}
	hand[3] = Card{5, Hearts}
	hand[4] = Card{7, Hearts}

	if isStraight(hand) == true {
		t.Error("Straight incorrectly recognized\n")
	}
	if isFlush(hand) == false {
		t.Error("Flush not recognized\n")
	}

	cards, v := Rank(hand)
	hr := GetHandKind(v)
	s := hr.String()

	fmt.Printf("h1 cards: %v\n", PrintHand(cards))
	fmt.Printf("h1 value: %d\n", v)
	fmt.Printf("rank = %b\n", v)
	fmt.Printf("hr = %s\n", s)

}

func TestSuitsEqual(t *testing.T) {

	hand := make([]Card, 5)

	fmt.Printf("TestSuitsEqual\n")
	hand[0] = Card{2, Hearts}
	hand[1] = Card{3, Hearts}
	hand[2] = Card{4, Hearts}
	hand[3] = Card{5, Hearts}
	hand[4] = Card{6, Hearts}

	cards1, v1 := Rank(hand)

	hand[0] = Card{2, Clubs}
	hand[1] = Card{3, Clubs}
	hand[2] = Card{4, Clubs}
	hand[3] = Card{5, Clubs}
	hand[4] = Card{6, Clubs}

	cards2, v2 := Rank(hand)

	if v2 != v1 {
		t.Error("Problem: Different suited identical hand have different rank\n")
	}
	_ = cards1
	_ = cards2
}

func TestPush(t *testing.T) {

	fmt.Printf("TestPush\n")

	hand1 := Hand([]Card{
		Card{Ace, Hearts},
		Card{Ace, Clubs},
		Card{King, Diamonds},
		Card{Queen, Spades},
		Card{Ten, Clubs},
		Card{Nine, Clubs},
		Card{Three, Hearts},
	})

	hand2 := Hand([]Card{
		Card{Ace, Spades},
		Card{Ace, Diamonds},
		Card{King, Hearts},
		Card{Queen, Clubs},
		Card{Ten, Hearts},
		Card{Eight, Clubs},
		Card{Two, Hearts},
	})

	AssertEquals(t, hand1, hand2)

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

	_, r1 := Rank(highestSetHand)
	fmt.Printf("Set description:  %s\n", r1.Describe())

	lowestStraight := Hand([]Card{
		Card{Ace, Hearts},
		Card{Two, Clubs},
		Card{Three, Diamonds},
		Card{Four, Clubs},
		Card{Five, Hearts},
	})
	_, r2 := Rank(lowestStraight)
	fmt.Printf("Straight description:  %s\n", r2.Describe())

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

	_, r3 := Rank(highestFlush)
	fmt.Printf("Flush description:  %s\n", r3.Describe())

	fullBoat := Hand([]Card{
		Card{Ace, Hearts},
		Card{Ace, Spades},
		Card{Ace, Diamonds},
		Card{Jack, Hearts},
		Card{Jack, Diamonds},
	})

	_, r4 := Rank(fullBoat)
	fmt.Printf("Fullhouse description:  %s\n", r4.Describe())

	AssertGreater(t, lowest2PairHand, highestPairHand)
	AssertGreater(t, lowest2PairHand, highestPairHand)
	AssertGreater(t, lowestSetHand, highest2PairHand)
	AssertGreater(t, lowestStraight, highestSetHand)
	AssertGreater(t, lowestFlush, highestStraight)

}

func AssertEquals(t *testing.T, h1 Hand, h2 Hand) {
	c1, r1 := Rank(h1)
	c2, r2 := Rank(h2)

	fmt.Printf("val: %.24b, for hand %v, desc: %s\n", r1, h1, r1.Describe())
	fmt.Printf("cards: %v\n", c1)
	fmt.Printf("%s\n", GetBinaryRankString(r1))
	fmt.Printf("val: %.24b, for hand %v, desc: %s\n", r2, h2, r2.Describe())
	fmt.Printf("cards: %v\n", c2)
	fmt.Printf("%s\n", GetBinaryRankString(r2))

	if r1 != r2 {
		fmt.Printf("Error hand %v is not equals to %v\n", h1, h2)
		//fmt.Printf("Hand %v is equals to %v\n", h1, h2)
	} else {
		//fmt.Printf("Hand %v is equals to %v\n", h1, h2)
		//t.Errorf("Hand %v is equals to %v\n", h1, h2)
		fmt.Printf("Hand %v is equals to %v\n", h1, h2)
	}
}

func AssertGreater(t *testing.T, h1 Hand, h2 Hand) {
	_, r1 := Rank(h1)
	_, r2 := Rank(h2)

	fmt.Printf("val: %.24b, for hand %v, desc: %s\n", r1, h1, r1.Describe())
	fmt.Printf("%s\n", GetBinaryRankString(r1))
	fmt.Printf("val: %.24b, for hand %v, desc: %s\n", r2, h2, r2.Describe())
	fmt.Printf("%s\n", GetBinaryRankString(r2))
	if r1 <= r2 {
		t.Errorf("Error hand %v greater than %v\n", h1, h2)
	}
}

func GetBinaryRankString(rank HandRank) string {

	var builder strings.Builder
	mask := 0x0000000F
	r := int(rank)
	for i := 0; i < 6; i++ {
		res := (r >> ((5 - i) * 4)) & mask
		builder.WriteString(fmt.Sprintf("%.4b ", res))
	}
	builder.WriteString("\n")
	return builder.String()

}
