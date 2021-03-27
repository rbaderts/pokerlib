package pokerlib

import (
	"fmt"
	"testing"
)

func TestOdds1(t *testing.T) {
	fmt.Printf("TestOdds 1\n")

	deck := NewDeck()
	deck.Shuffle()
	hands := make([][2]Card, 3)

	card1 := Card{13, Clubs}
	card2 := Card{14, Spades}
	card3 := Card{13, Spades}
	card4 := Card{3, Clubs}
	hands[0] = [2]Card{card1, card2}
	hands[1] = [2]Card{card3, card4}

	commonCards := []Card{}

	result := CalculateOdds(deck, hands, commonCards, 1000)

	fmt.Printf(" with 1000 iterations\n")
	for i := 0; i <= 2; i++ {
		fmt.Printf("hand %d, %v, win = %f, tie = %f\n", i+1, hands[i], result[i].Wins, result[i].Ties)
	}
	fmt.Printf(" with 10000 iterations\n")

	result = CalculateOdds(deck, hands, commonCards, 10000)
	for i := 0; i <= 2; i++ {
		fmt.Printf("hand %d, %v, win = %f, tie = %f\n", i+1, hands[i], result[i].Wins, result[i].Ties)
	}

}

func TestOdds(t *testing.T) {
	fmt.Printf("TestOdds\n")

	deck := NewDeck()
	deck.Shuffle()
	hands := make([][2]Card, 3)

	card1 := Card{7, Hearts}
	card2 := Card{9, Hearts}
	card3 := Card{14, Hearts}
	card4 := Card{11, Diamonds}
	card5 := Card{14, Clubs}
	card6 := Card{9, Clubs}
	card7 := Card{11, Clubs}
	card8 := Card{7, Spades}
	card9 := Card{4, Diamonds}
	hands[0] = [2]Card{card1, card2}
	hands[1] = [2]Card{card3, card4}
	hands[2] = [2]Card{card5, card6}

	commonCards := []Card{card7, card8, card9}

	deck.removeCard(card1)
	deck.removeCard(card2)
	deck.removeCard(card3)
	deck.removeCard(card4)
	deck.removeCard(card5)
	deck.removeCard(card6)
	deck.removeCard(card7)
	deck.removeCard(card8)
	deck.removeCard(card9)

	result := CalculateOdds(deck, hands, commonCards, 1000)

	fmt.Printf(" with 1000 iterations\n")
	for i := 0; i <= 2; i++ {
		fmt.Printf("hand %d, %v, win = %f, tie = %f\n", i+1, hands[i], result[i].Wins, result[i].Ties)
	}
	fmt.Printf(" with 10000 iterations\n")

	result = CalculateOdds(deck, hands, commonCards, 10000)
	for i := 0; i <= 2; i++ {
		fmt.Printf("hand %d, %v, win = %f, tie = %f\n", i+1, hands[i], result[i].Wins, result[i].Ties)
	}
}
