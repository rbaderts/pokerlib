package pokerlib

import (
	"encoding/json"
	"fmt"
	"github.com/dgryski/go-pcgr"
	"hash/maphash"
	"io/ioutil"
	_ "math/rand"
	"net"
	"strconv"
	"time"
	_ "time"
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
	card := this.RemainingCards[this.Position]
	this.Position += 1
	fmt.Printf("\ndraw %s\n", card.String())
	return card
}

func (this *Deck) BurnCard() {
	this.Position += 1
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
	rnd := pcgr.Rand{uint64(time.Now().UnixNano()), 0x00004444}

	this.shuffle(rnd, len(this.RemainingCards), func(i, j int) {
		this.RemainingCards[i], this.RemainingCards[j] = this.RemainingCards[j], this.RemainingCards[i]
	})
}

func (this *Deck) shuffle (r pcgr.Rand, n int, swap func(i, j int)) {
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

func (this *Deck) String() string {
	str := ""
	str += "Position: "
	str += strconv.Itoa(this.Position)
	str += " - "
	for _, card := range this.RemainingCards {
		str += fmt.Sprintf("%v", card.String())
		str += ":"
	}
	str += "\n"
	return str
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
