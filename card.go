/*
  all rights reserved

    : Rick Badertscher (rick.badertscher@gmail.com)
*/
package pokerlib

import (
	"fmt"
	"strconv"
)

/*
   card Index values are 2=1, Ace=13
*/
type Index int

func (this Index) String() string {

	switch this {
	case Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten:
		return strconv.Itoa(int(this))
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case Ace:
		return "Ace"
	}
	return "Unknown: " + strconv.Itoa(int(this))
}

const (
	Two   Index = 2
	Three Index = 3
	Four  Index = 4
	Five  Index = 5
	Six   Index = 6
	Seven Index = 7
	Eight Index = 8
	Nine  Index = 9
	Ten   Index = 10
	Jack  Index = 11
	Queen Index = 12
	King  Index = 13
	Ace   Index = 14
)

type Suit int

const (
	Spades Suit = iota + 1
	Hearts
	Diamonds
	Clubs
)

func (this Suit) String() string {

	switch this {
	case Clubs:
		return "clubs"
	case Hearts:
		return "hearts"
	case Diamonds:
		return "diamonds"
	case Spades:
		return "spades"
	}
	return "NA"
}

type Card struct {
	Index Index `json:"index"`
	Suit  Suit  `json:"suit"`
}

type AbsoluteValue int

func (this AbsoluteValue) String() string {
	return fmt.Sprintf("%0x", int(this))
}

func (this Card) GetCardValue() AbsoluteValue {
	// 4 bits:    for card index

	v := int(this.Index)

	return AbsoluteValue(v)
}

func (this Card) Equals(c Card) bool {
	if this.Index == c.Index && this.Suit == c.Suit {
		return true
	}
	return false
}

func (this Card) String() string {

	uc := 0x1F0A0

	uc += 0x10 * (int(this.Suit) - 1)
	uc += int(this.Index)

	rankStr := ""
	if this.Index == 10 {
		rankStr = "T"
	} else if this.Index == 11 {
		rankStr = "J"
	} else if this.Index == 12 {
		rankStr = "Q"
	} else if this.Index == 13 {
		rankStr = "K"
	} else if this.Index == 14 {
		rankStr = "A"
	} else {
		rankStr = strconv.Itoa(int(this.Index))
	}

	suitStr := ""
	switch this.Suit {
	case Hearts:
		suitStr = "\u2665"
	case Diamonds:
		suitStr = "\u2666"
	case Clubs:
		suitStr = "\u2663"
	case Spades:
		suitStr = "\u2660"
	}

	return fmt.Sprintf("%s(%c)", rankStr+suitStr, uc)

}
