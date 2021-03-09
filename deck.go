package pokerlib

import (
	"encoding/json"
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

/*
  A deck is a set of 52 cards that can be shuffled and drawn on.   Cards
  are drawn sequentialy and are removed from the deck */

type DeckCard struct {
	card     Card
	next     *DeckCard
	prev     *DeckCard
	index    int
	borrowed bool
}

type Deck struct {
	//RemainingCards [52]Card `json:"cards"`
	//Position       int
	topCard    *DeckCard
	bottomCard *DeckCard
	cardIndex  map[Card]*DeckCard
	orderIndex []*DeckCard
	totalCards int
	rnd        *pcgr.Rand
}

func ReadDeck(filename string) *Deck {
	file, _ := ioutil.ReadFile(filename)
	deck := Deck{}
	err := json.Unmarshal([]byte(file), &deck)

	if err != nil {
		fmt.Printf("Error reading file %s\n", filename)

	}
	//deck.Position = 0
	return &deck
}

func NewDeck() *Deck {
	deck := new(Deck)
	deck.rnd = &pcgr.Rand{uint64(time.Now().UnixNano()), 0x00004443}
	//deck.Position = 0
	///	deck.RemainingCards = make([]Card, 52)
	deck.orderIndex = make([]*DeckCard, 0)
	deck.cardIndex = make(map[Card]*DeckCard, 0)

	index := 0
	for rank := Two; rank <= Ace; rank++ {
		for suit := 1; suit <= 4; suit++ {
			card := Card{Index(rank), Suit(suit)}
			deck.appendCard(card)

			index++

			//deck.RemainingCards[index] = card
			//index++
		}
	}
	return deck

}

func (this *Deck) Copy() *Deck {
	deck := new(Deck)
	//deck.rnd = &pcgr.Rand{uint64(time.Now().UnixNano()), 0x00004443}
	deck.rnd = this.rnd
	deck.orderIndex = make([]*DeckCard, 0)
	deck.cardIndex = make(map[Card]*DeckCard, 0)

	dc := this.topCard
	for {
		deck.appendCard(dc.card)
		if dc.next == nil {
			break
		}
		dc = dc.next
	}

	return deck

}

func (this *Deck) appendCard(card Card) {

	deckCard := &(DeckCard{card, nil, nil, 0, false})
	if this.topCard == nil {
		this.topCard = deckCard
		this.bottomCard = deckCard
	} else {
		this.bottomCard.next = deckCard
		deckCard.prev = this.bottomCard
		this.bottomCard = deckCard
	}
	this.cardIndex[card] = deckCard
	this.orderIndex = append(this.orderIndex, deckCard)
	deckCard.index = len(this.orderIndex) - 1
	this.totalCards += 1

}

func (this *Deck) removeCard(card Card) {

	//fmt.Printf("removeCard %v\n", card)
	deckcard := this.cardIndex[card]
	//fmt.Printf("cared index = %d\n", deckcard.index)
	//index := deckcard.index

	if deckcard.prev != nil {
		deckcard.prev.next = deckcard.next
	}
	if deckcard.next != nil {
		deckcard.next.prev = deckcard.prev
	}
	if deckcard == this.topCard {
		this.topCard = deckcard.next
		this.topCard.prev = nil
	}
	if deckcard == this.bottomCard {
		this.bottomCard = deckcard.prev
		this.bottomCard.next = nil
	}
	this.orderIndex[deckcard.index] = nil
	this.totalCards -= 1

}

// mark all borrowed cards unborrowed
func (this *Deck) Restore() {
	var deckcard *DeckCard = this.topCard
	for {
		deckcard.borrowed = false
		if deckcard.next == nil {
			break
		}
		deckcard = deckcard.next
	}

}

func (this *Deck) BorrowRandom() Card {
	index := int(this.rnd.Bound(uint32(this.totalCards)))

	var deckcard *DeckCard
	for {
		deckcard = this.orderIndex[index]
		if deckcard == nil || deckcard.borrowed == true {
			index += 1
			if index >= len(this.orderIndex)-1 {
				index = 0
			}
		} else {
			break
		}
	}
	card := deckcard.card
	deckcard.borrowed = true
	//	this.removeCard(card)
	return card

}

func (this *Deck) DrawCard() Card {
	deckcard := this.topCard
	card := deckcard.card
	this.topCard = this.topCard.next
	this.topCard.prev = nil
	this.orderIndex[deckcard.index] = nil
	delete(this.cardIndex, card)
	this.totalCards -= 1
	return card
}

func (this *Deck) BurnCard() {
	//	this.Position += 1
	this.DrawCard()
}

/*
func (this *Deck) Shuffle() {
	//	r := rand.New(rand.NewSource(time.Now().Unix()))
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(this.RemainingCards), func(i, j int) {
		this.RemainingCards[i], this.RemainingCards[j] = this.RemainingCards[j], this.RemainingCards[i]
	})
}
*/

func (this *Deck) Shuffle() {

	//rnd := pcgr.Rand{ uint64(time.Now().UnixNano()), getMacAddrHash())
	rnd := pcgr.Rand{uint64(time.Now().UnixNano()), 0x00004441}

	this.shuffle(rnd, this.totalCards, func(i, j int) {

		iDeckCard := this.orderIndex[i]
		jDeckCard := this.orderIndex[j]
		iCard := iDeckCard.card
		jCard := jDeckCard.card
		iDeckCard.card = jCard
		jDeckCard.card = iCard
		this.cardIndex[iCard] = jDeckCard
		this.cardIndex[jCard] = iDeckCard
		//		save := iDeckCard.index
		///		iDeckCard.index = jDeckCard.index
		//		jDeckCard.index = save
		iDeckCard.index, jDeckCard.index = jDeckCard.index, iDeckCard.index
		///		this.RemainingCards[i], this.RemainingCards[j] = this.RemainingCards[j], this.RemainingCards[i]
	})
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

	//builder.WriteString(fmt.Sprintf("Deck (totalCards:%d\n", this.totalCards))
	d := this.topCard
	for {
		if d == nil {
			break
		}
		card := d.card
		builder.WriteString(fmt.Sprintf("%v : ", card.String()))
		d = d.next
	}
	builder.WriteString("\n")
	return builder.String()

}
func getMacAddrHash() uint64 {
	s, _ := getMacAddr()
	var h maphash.Hash
	h.Write(([]byte)(s))
	return h.Sum64()

	/*
		h.W
		io.WriteString(h, s)
		return h.Sum(nil)
		fmt.Printf("%x", h.Sum(nil))
	*/

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
