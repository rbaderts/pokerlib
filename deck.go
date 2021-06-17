package pokerlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgryski/go-pcgr"
	"hash/maphash"
	"io/ioutil"
	_ "math/rand"
	"net"
	"strings"
	"time"
	_ "time"
)

type Deck struct {
	cards     [52]*Card
	cardIndex map[Card]int
	top       int
	count     int
	rnd       *pcgr.Rand
}

func ReadDeck(filename string) *Deck {
	file, _ := ioutil.ReadFile(filename)
	deck := Deck{}
	err := json.Unmarshal([]byte(file), &deck)

	if err != nil {
		fmt.Printf("Error reading file %s\n", filename)

	}
	return &deck
}

func NewDeck() *Deck {
	deck := new(Deck)
	deck.rnd = &pcgr.Rand{uint64(time.Now().UnixNano()), 0x00004443}
	deck.cardIndex = make(map[Card]int, 0)

	index := 0
	for rank := Two; rank <= Ace; rank++ {
		for suit := 1; suit <= 4; suit++ {
			card := Card{Index(rank), Suit(suit)}
			deck.appendCard(card)
			index++
		}
	}
	return deck

}

func (this *Deck) Copy() *Deck {
	deck := new(Deck)
	deck.rnd = this.rnd
	deck.cardIndex = make(map[Card]int, 0)

	for i, c := range this.cards {
		deck.cards[i] = c
	}
	deck.count = this.count
	for c, i := range this.cardIndex {
		deck.cardIndex[c] = i
	}
	return deck

}

func (this *Deck) appendCard(card Card) {

	this.cards[this.count] = &card
	this.cardIndex[card] = this.count
	this.count += 1

}

func (this *Deck) removeCard(card Card) {

	index := this.cardIndex[card]
	if index == this.top {
		for i := this.top + 1; i < 52; i++ {
			if this.cards[i] != nil {
				this.top = i
				break
			}
		}
	}

	this.cards[index] = nil
	this.count -= 1

}

func (this *Deck) BorrowRandom() Card {
	var card *Card
	for {

		index := int(this.rnd.Bound(uint32(52)))
		card = this.cards[index]
		if card == nil {
			continue
		}
		this.cards[index] = nil
		break
	}

	return *card

}

func (this *Deck) ReturnCard(card Card) error {
	i := this.cardIndex[card]
	if this.cards[i] != nil {
		return errors.New(fmt.Sprintf("Card %v, not borrowed\n", card))
	}
	this.cards[i] = &card
	return nil
}

func (this *Deck) DrawCard() Card {
	card := this.cards[this.top]
	this.cards[this.top] = nil
	this.top += 1
	return *card

}

func (this *Deck) BurnCard() {
	this.DrawCard()
}

func (this *Deck) buildCardIndex() {
	this.cardIndex = make(map[Card]int, 0)
	for i, c := range this.cards {
		this.cardIndex[*c] = i
	}
}
func (this *Deck) Shuffle() {

	rnd := pcgr.Rand{uint64(time.Now().UnixNano()), 0x00004441}

	this.shuffle(rnd, 52, func(i, j int) {
		this.cards[i], this.cards[j] = this.cards[j], this.cards[i]
	})

	this.buildCardIndex()
}

func (this *Deck) shuffle(r pcgr.Rand, n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	// Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	// Shuffle really ought not be called with n that doesn't fit in 32 bits.
	// Not only will it take a very long time, but with 2³¹! possible permutations,
	// there's no way that any PRNG can have a big enough internal state to
	// generate even a minuscule percentage of the possible permutations.
	// Nevertheless, the right API signature accepts an int n, so handle it as best we can.
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(r.Bound(uint32(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(r.Bound(uint32(i + 1)))
		swap(i, j)
	}
}

func (this Deck) String() string {

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Deck (totalCards:%d, top: %d\n", this.count, this.top))
	for i, c := range this.cards {
		str := "none"
		if c != nil {
			str = c.String()
		}
		builder.WriteString(fmt.Sprintf("%d - %v : ", i, str))
	}
	builder.WriteString("\n")
	return builder.String()

}
func getMacAddrHash() uint64 {
	s, _ := getMacAddr()
	var h maphash.Hash
	h.Write(([]byte)(s))
	return h.Sum64()

}
func getMacAddr() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "FALLBACK", err
	}
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			return a, nil
		}
	}
	return "FALLBACK", nil
}
