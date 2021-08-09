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

func (index Index) String() string {

	switch index {
	case Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten:
		return strconv.Itoa(int(index))
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case Ace:
		return "Ace"
	}
	return "Unknown: " + strconv.Itoa(int(index))
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

func (suid Suit) String() string {

	switch suid {
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

func (card Card) GetAbsoluteValue() CardAbsoluteValue {
	return CardToCardAbsoluteValue(card)
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

func (code CardCode) String() string {
	return fmt.Sprintf("%0x", int(code))
}

func (card Card) GetCardCode() CardCode {

	v := int(card.Suit - 1)
	val := v | (int(card.Index) << 2)

	return CardCode(val)
}

func (card Card) GetIndexValue() CardCode {
	return CardCode(card.Index)
}

func (card Card) Equals(c Card) bool {
	if card.Index == c.Index && card.Suit == c.Suit {
		return true
	}
	return false
}

func (card Card) String() string {

	uc := uint(0x1F0A0)

	uc += 0x10 * (uint(card.Suit) - 1)
	uc += uint(card.Index)

	rankStr := ""
	if card.Index == 10 {
		rankStr = "T"
	} else if card.Index == 11 {
		rankStr = "J"
	} else if card.Index == 12 {
		rankStr = "Q"
	} else if card.Index == 13 {
		rankStr = "K"
	} else if card.Index == 14 {
		rankStr = "A"
	} else {
		rankStr = strconv.Itoa(int(card.Index))
	}

	suitStr := ""
	switch card.Suit {
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
