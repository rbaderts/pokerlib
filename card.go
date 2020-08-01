package pokerlib

import (
	"fmt"
	"strconv"
)

type Index int

const (
	Two Index = iota + 1
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

type Card struct {
	Index Index `json:"index"`
	Suit  Suit  `json:"suit"`
}

type AbsoluteValue int

func (this AbsoluteValue) String() string {
	return fmt.Sprintf("%0x", int(this))
}

func (this Card) AbsoluteValue() AbsoluteValue {
	// 4 bits:    for card index

	v := 0
	v = int(this.Index)
	//v |= int(this.Suit)

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
	} else if this.Index == 1 {
		rankStr = "A"
	} else {
		rankStr = strconv.Itoa(int(this.Index))
	}

	suitStr := ""
	if this.Suit == Hearts {
		suitStr = "\u2665"
	} else if this.Suit == Diamonds {
		suitStr = "\u2666"
	} else if this.Suit == Clubs {
		suitStr = "\u2663"
	} else if this.Suit == Spades {
		suitStr = "\u2660"
	}

	return fmt.Sprintf("%s(%c)", rankStr+suitStr, uc)

}
