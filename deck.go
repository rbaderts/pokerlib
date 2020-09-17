package pokerlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

/*
  A deck is a set of 52 cards that can be shuffled and drawn on.   Cards
  are drawn sequentialy and are removed from the deck */
type Deck struct {
	RemainingCards [52]Card `json:"cards"`
	Position       int
}

func ReadDeck(filename string) *Deck {
	file, _ := ioutil.ReadFile(filename)
	deck := Deck{}
	err := json.Unmarshal([]byte(file), &deck)


	if err != nil {
		fmt.Printf("Error reading file %s\n", filename)

	}
	deck.Position = 0
	return &deck
}

func NewDeck() *Deck {
	deck := new(Deck)
	deck.Position = 0
///	deck.RemainingCards = make([]Card, 52)
	index := 0
	for rank := Two; rank <= Ace; rank++ {
		for suit := 1; suit <= 4; suit++ {
			card := Card{Index(rank), Suit(suit)}
			deck.RemainingCards[index] = card
			index++
		}
	}
	return deck

}

func (this *Deck) DrawCard() Card {
	card := this.RemainingCards[this.Position];
	this.Position += 1
	fmt.Printf("\ndraw %s\n", card.String())
	return card
}

func (this *Deck) BurnCard() {
	this.Position += 1
}

func (this *Deck) Shuffle() {
	//	r := rand.New(rand.NewSource(time.Now().Unix()))
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(this.RemainingCards), func(i, j int) {
		this.RemainingCards[i], this.RemainingCards[j] = this.RemainingCards[j], this.RemainingCards[i]
	})
}

func (this *Deck) String() string {
	str := ""
	str += "Position: "
	str += strconv.Itoa(this.Position)
	str += " - ";
	for _, card := range this.RemainingCards {
		str += fmt.Sprintf("%v", card.String())
		str += ":"
	}
	str += "\n"
	return str
}
