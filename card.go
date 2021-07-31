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

    A card can be uniquely identified using 6 bits

    card representations:

    	Absolute Value:    ((Suit-1) * 13) + (Index-2):     range 0 - 52
		Encoded:           2 bits for Suit, 4 bits for Index


   card Index values are 2=2, Ace=14
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
	Two Index = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
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

type CardAbsoluteValue uint16

func (this Card) GetAbsoluteValue() CardAbsoluteValue {
	return CardToCardAbsoluteValue(this)
}

func CardToCardAbsoluteValue(c Card) CardAbsoluteValue {

	v := 13 * (int(c.Suit) - 1)
	v += int(c.Index) - 2

	return CardAbsoluteValue(v)
}

func CardAbsoluteValueToCard(v CardAbsoluteValue) Card {

	suit := ((int32(v) - 1) / 13) + 1
	index := ((int32(v) - 1) % 13) + 2

	return Card{Index(index), Suit(suit)}

}

type CardCode int

func (this CardCode) String() string {
	return fmt.Sprintf("%0x", int(this))
}

func (this Card) GetCardCode() CardCode {

	v := int(this.Suit - 1)
	val := v | (int(this.Index) << 2)

	return CardCode(val)
}

func (this Card) GetIndexValue() CardCode {
	return CardCode(this.Index)
}

func (this Card) Equals(c Card) bool {
	if this.Index == c.Index && this.Suit == c.Suit {
		return true
	}
	return false
}

func (this Card) String() string {

	uc := uint(0x1F0A0)

	uc += 0x10 * (uint(this.Suit) - 1)
	uc += uint(this.Index)

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
