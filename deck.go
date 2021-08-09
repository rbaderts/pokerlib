package pokerlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "math/rand"
	"strings"
	"time"

	"github.com/dgryski/go-pcgr"
)

type Deck struct {
	Cards     [52]*Card    `json:"cards"`
	CardIndex map[Card]int `json:"cardIndex"`
	Top       int          `json:"top"`
	Count     int          `json:"count"`
	Rnd       *pcgr.Rand   `json:"rnd"`
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
	deck.Rnd = &pcgr.Rand{State: uint64(time.Now().UnixNano()), Inc: 0x00004443}
	deck.CardIndex = make(map[Card]int)

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

func (dec *Deck) Copy() *Deck {
	deck := new(Deck)
	deck.Rnd = dec.Rnd
	deck.CardIndex = make(map[Card]int)

	for i, c := range dec.Cards {
		deck.Cards[i] = c
	}
	deck.Count = dec.Count
	for c, i := range dec.CardIndex {
		deck.CardIndex[c] = i
	}
	return deck

}

func (deck *Deck) appendCard(card Card) {

	deck.Cards[deck.Count] = &card
	deck.CardIndex[card] = deck.Count
	deck.Count += 1

}

func (deck *Deck) removeCard(card Card) {

	index := deck.CardIndex[card]
	if index == deck.Top {
		for i := deck.Top + 1; i < 52; i++ {
			if deck.Cards[i] != nil {
				deck.Top = i
				break
			}
		}
	}

	deck.Cards[index] = nil
	deck.Count -= 1

}

func (deck *Deck) BorrowRandom() Card {
	var card *Card
	for {

		index := int(deck.Rnd.Bound(uint32(52)))
		card = deck.Cards[index]
		if card == nil {
			continue
		}
		deck.Cards[index] = nil
		break
	}

	return *card

}

func (deck *Deck) ReturnCard(card Card) error {
	i := deck.CardIndex[card]
	if deck.Cards[i] != nil {
		return fmt.Errorf("Card %v, not borrowed", card)
	}
	deck.Cards[i] = &card
	return nil
}

func (deck *Deck) DrawCard() Card {
	card := deck.Cards[deck.Top]
	deck.Cards[deck.Top] = nil
	deck.Top += 1
	return *card

}

func (deck *Deck) BurnCard() {
	deck.DrawCard()
}

func (deck *Deck) buildCardIndex() {
	deck.CardIndex = make(map[Card]int)
	for i, c := range deck.Cards {
		deck.CardIndex[*c] = i
	}
}
func (deck *Deck) Shuffle() {

	rnd := pcgr.Rand{State: uint64(time.Now().UnixNano()), Inc: uint64(0x00004441)}

	deck.shuffle(rnd, 52, func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})

	deck.buildCardIndex()
}

func (deck *Deck) shuffle(r pcgr.Rand, n int, swap func(i, j int)) {
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

func (deck Deck) String() string {

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Deck (totalCards:%d, top: %d\n", deck.Count, deck.Top))
	for i, c := range deck.Cards {
		str := "none"
		if c != nil {
			str = c.String()
		}
		builder.WriteString(fmt.Sprintf("%d - %v : ", i, str))
	}
	builder.WriteString("\n")
	return builder.String()

}

/*
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

*/
