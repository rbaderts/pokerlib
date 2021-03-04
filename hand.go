package pokerlib

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"sort"
	"strings"
)

/*

   HandRank encodes a hand ranking in 32 bit.    It is encoded as follows:

      4 bits per card index.   So first least significant 20 bit is
	     room for the indexes of all 5 cards in a hand.

      The next 4 bits:   HAND_KIND_SHIFT-23 encode the 4 bit hand kind

        Ex:    For the hand K 7 8 4 2 (not-suited).   The value would be:

		0000 0000 0000 1101 0111 1000 0100 0010



*/

type HandRank uint

const HAND_KIND_MASK = 0xFFF00000
const HAND_KIND_SHIFT = 20
const CARD_MASK = 0x0000000F

type HandKind int8 // Pair, Straight, etc.

var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgGreen).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

var handEvalLog *log.Logger

const (
	HighCard      HandKind = 0b0000
	Pair          HandKind = 0b0001
	TwoPair       HandKind = 0b0010
	ThreeOfAKind  HandKind = 0b0011
	Straight      HandKind = 0b0100
	Flush         HandKind = 0b0101
	FullHouse     HandKind = 0b0111
	FourOfAKind   HandKind = 0b1000
	StraightFlush HandKind = 0b1001
)

func init() {
	filename := os.Getenv("POKERLIB_HANDEVAL_LOG")

	if filename != "" {
		f, err := os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		handEvalLog = log.New(f, "", log.LstdFlags)
	}

}

func GetHandKind(rank HandRank) HandKind {
	return HandKind(rank&HAND_KIND_MASK) >> HAND_KIND_SHIFT
}

func (this HandRank) GetCard(pos int) Index {
	shiftby := (5 - pos - 1) * 4
	index := Index((this >> shiftby) & CARD_MASK)
	//	index := Index((this&HIGH_CARD_MASK) >> HIGH_CARD_SHIFT)
	return index
}

func (this HandKind) String() string {

	switch this {
	case HighCard:
		return "HighCard"
	case Pair:
		return "Pair"
	case TwoPair:
		return "2 Pair"
	case ThreeOfAKind:
		return "3 of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "4 of A Kind"
	case StraightFlush:
		return "Straight Flush"
	}
	return ""
}

type Hand []Card

func (this HandRank) Describe() string {

	r := HandKind((int(this) & HAND_KIND_MASK) >> HAND_KIND_SHIFT)
	if r == Pair {
		return fmt.Sprintf("a %s of %s's with a %s kicker",
			red(r.String()),
			yellow(this.GetCard(0).String()),
			yellow(this.GetCard(2).String()))
	} else if r == TwoPair {
		return fmt.Sprintf("%s, %s's and %s's, with a %s kicker",
			red(r.String()),
			this.GetCard(0).String(),
			this.GetCard(2).String(),
			this.GetCard(4).String())
	} else if r == ThreeOfAKind {
		return fmt.Sprintf("Trip %s's with a %s kicker",
			yellow(this.GetCard(0).String()),
			this.GetCard(3).String())
	} else if r == Straight {
		return fmt.Sprintf("a %s high %s",
			this.GetCard(0).String(), red(r.String()))
	} else if r == Flush {
		return fmt.Sprintf("a %s high %s",
			this.GetCard(0).String(), red(r.String()))
	} else if r == HighCard {
		return fmt.Sprintf("%s high",
			this.GetCard(0).String())
	} else if r == FullHouse {
		return fmt.Sprintf("a %s %s's full of %s's",
			red(r.String()),
			this.GetCard(0).String(),
			this.GetCard(3).String())
	} else if r == FourOfAKind {
		return fmt.Sprintf("4 %s's",
			yellow(this.GetCard(0).String()))
	} else if r == StraightFlush {
		return fmt.Sprintf("A straight flush! %s high",
			this.GetCard(0).String())
	}
	return r.String() + " Not yet implemented"

}

func PrintHand(cards []Card) string {

	var builder strings.Builder
	for i, c := range cards {

		builder.WriteString(fmt.Sprintf("%v ", c.String()))
		if i < len(cards)-1 {
			builder.WriteString(fmt.Sprintf(" - "))
		}
	}

	return builder.String()

}

/**
 */
func Rank(cards []Card) ([]Card, HandRank) {

	if len(cards) == 7 {
		if handEvalLog != nil {
			handEvalLog.Printf("All cards:  %v\n", cards)
		}

	}

	checkedSets := make(map[HandRank][]Card, 0)
	topCards, top := DoRank(cards, checkedSets)

	if handEvalLog != nil {
		handEvalLog.Printf("Checked sets: \n")
	}
	for r, s := range checkedSets {
		if handEvalLog != nil {
			handEvalLog.Printf("rank: %d : cards: %v : %s\n", r, s, r.Describe())
		}
	}
	if handEvalLog != nil {
		handEvalLog.Printf("Top Rank: %v\n", top)
	}

	return topCards, top

}

func DoRank(cards []Card, checkedSets map[HandRank][]Card) ([]Card, HandRank) {

	if len(cards) > 5 {
		var top HandRank
		var topCards []Card
		for i, _ := range cards {
			var subset = make([]Card, len(cards)-1)
			count := 0
			for j, c := range cards {
				if j != i {
					subset[count] = c
					count += 1
				}
			}
			topC, r := DoRank(subset, checkedSets)
			if r > top {
				top = r
				topCards = topC
			}
		}
		return topCards, top
	}
	r := RankHand(cards)
	checkedSets[r] = cards
	return cards, RankHand(cards)
}

/*
   Returns the difference in value, a positive result means rank1 was higher than rank2
   A boolean is also returned indicating if the winning rank was determined by a kicker
*/

/*
func CompareHandRanks(rank1 HandRank, rank2 HandRank) (int, bool) {
	hk1 := GetHandKind(hk1)
	hk2 := GetHandKind(hk2)
	if hk1 == hk2 && hk1 == Pair
}
*/

/*
   Called only for 5 cards
*/
func RankHand(cards []Card) HandRank {

	sortedCards := make([]Card, 5)
	orderedCards := make([]Card, 5)
	// Histogram is a map of ranks to frequency
	histos := buildHistogram(cards)
	//	fmt.Printf("histos = %v\n", histos)

	topCount := 0
	secondCount := 0
	for i, h := range histos {
		if i == 0 {
			topCount = h.Count
		}
		if i == 1 {
			secondCount = h.Count
		}
	}

	copy(sortedCards, cards)
	SortCards(sortedCards)

	isFlush := isFlush(sortedCards)
	isStraight := isStraight(sortedCards)
	value := 0

	if isStraight && isFlush {
		value = int(StraightFlush) << HAND_KIND_SHIFT
		i := 4

		if sortedCards[0].Index == Two && sortedCards[4].Index == Ace {
			orderedCards[0] = sortedCards[3]
			orderedCards[1] = sortedCards[2]
			orderedCards[2] = sortedCards[1]
			orderedCards[3] = sortedCards[0]
			orderedCards[4] = sortedCards[4]
		} else {

			for _, c := range sortedCards {
				orderedCards[i] = c
				i -= 1
			}
		}

		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(orderedCards) - i - 1) * 4))
		}

	} else if topCount == 4 {
		mainCard := histos[0].Index
		i := 1
		for _, c := range sortedCards {
			if c.Index == mainCard {
				orderedCards[i] = c
				i++
			} else {
				orderedCards[0] = c
			}
		}

		value = int(FourOfAKind) << HAND_KIND_SHIFT
		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(orderedCards) - i - 1) * 4))
		}

	} else if topCount == 3 && secondCount == 2 {
		topCard := histos[0].Index

		i := 0
		j := 3
		for _, c := range sortedCards {
			if c.Index == topCard {
				orderedCards[i] = c
				i += 1
			} else {
				orderedCards[j] = c
				j += 1
			}
		}
		value = int(FullHouse) << HAND_KIND_SHIFT
		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(orderedCards) - i - 1) * 4))

		}

	} else if isFlush {

		value = int(Flush) << HAND_KIND_SHIFT
		i := 4
		for _, c := range sortedCards {
			orderedCards[i] = c
			i--
		}
		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(orderedCards) - i - 1) * 4))
		}
		return HandRank(value)

	} else if isStraight {

		value = int(Straight) << HAND_KIND_SHIFT

		if sortedCards[0].Index == Two && sortedCards[4].Index == Ace {
			orderedCards[0] = sortedCards[3]
			orderedCards[1] = sortedCards[2]
			orderedCards[2] = sortedCards[1]
			orderedCards[3] = sortedCards[0]
			orderedCards[4] = sortedCards[4]
		} else {
			i := 4
			for _, c := range sortedCards {
				orderedCards[i] = c
				i--
			}
		}
		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(orderedCards) - i - 1) * 4))
		}

	} else if topCount == 3 {
		topCard := histos[0].Index

		i := 0
		j := 4
		//extraCards := make([]Card, 2)
		for _, c := range sortedCards {
			if c.Index == topCard {
				orderedCards[i] = c
				i += 1
			} else {
				orderedCards[j] = c
				j -= 1
			}
		}

		value = int(ThreeOfAKind) << HAND_KIND_SHIFT

		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i - 1) * 4))
		}

	} else if topCount == 2 && secondCount == 2 {

		topCard := histos[0].Index
		nextCard := histos[1].Index
		i := 0
		j := 2
		for _, c := range sortedCards {
			if c.Index == topCard {
				orderedCards[j] = c
				j += 1
			} else if c.Index == nextCard {
				orderedCards[i] = c
				i += 1
			} else {
				orderedCards[4] = c
			}
		}

		value = int(TwoPair) << HAND_KIND_SHIFT
		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(orderedCards) - i - 1) * 4))
		}
	} else if topCount == 2 {
		topCard := histos[0].Index
		i := 0
		j := 4
		//extraCards := make([]Card, 3)
		for _, c := range sortedCards {
			if c.Index == topCard {
				orderedCards[i] = c
				i += 1
			} else {
				orderedCards[j] = c
				j -= 1
			}
		}

		value = int(Pair) << HAND_KIND_SHIFT

		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i - 1) * 4))
		}
	} else {
		i := 4
		for _, c := range sortedCards {
			orderedCards[i] = c
			i -= 1
		}
		value = int(HighCard) << HAND_KIND_SHIFT

		for i, c := range orderedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i - 1) * 4))
		}
	}

	return HandRank(value)

}

func getHighCard(cards []Card) Index {
	highest := Index(Two)
	for _, c := range cards {
		if c.Index > highest {
			highest = c.Index
		}
	}
	return highest

}

func isFlush(cards []Card) bool {
	if len(cards) != 5 {
		fmt.Printf("Error isFlush called on a set of more than 5 cards\n")
		return false
	}
	suit := -1
	flush := true
	for _, c := range cards {
		if suit == -1 {
			suit = int(c.Suit)
		}
		if suit != int(c.Suit) {
			flush = false
			break
		}
	}
	return flush

}

func isStraight(cards []Card) bool {

	if len(cards) != 5 {
		fmt.Printf("Error isStraight called on a set of more than 5 cards\n")
		return false

	}
	low := 0
	last := 0

	for i, c := range cards {

		if i == 0 {
			low = int(c.Index)
			last = low
		} else if i == 4 && low == 2 && c.Index == Ace {
			return true
		} else {
			if int(c.Index) != last+1 {
				return false
			}
			last += 1
		}
	}
	return true
}

type Histogram struct {
	Count int
	Index Index
}

type ByCount []Histogram

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count }

func buildHistogram(cards []Card) []Histogram {

	res := make(map[Index]Histogram)
	for _, c := range cards {
		r := c.Index
		v, has := res[r]
		if !has {
			res[r] = Histogram{1, r}
		} else {
			n := v.Count
			res[r] = Histogram{n + 1, r}
		}
	}

	i := 0
	var counts []Histogram
	for _, h := range res {
		counts = append(counts, h)
		i += 1
	}
	sort.Sort(ByCount(counts))

	return counts

}

type ByRank []Card

func (a ByRank) Len() int           { return len(a) }
func (a ByRank) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRank) Less(i, j int) bool { return a[i].Index < a[j].Index }

func SortCards(cards []Card) {
	sort.Sort(ByRank(cards))
}

func (this Hand) ContainsCard(card Card) bool {
	for _, c := range this {
		if c.Equals(card) {
			return true
		}
	}
	return false

}

func (this Hand) Equals(h Hand) bool {
	if len(this) != len(h) {
		return false
	}
	for _, c := range this {
		if !h.ContainsCard(c) {
			return false
		}
	}

	return true
}

type ByGetCardValue []Card

func (a ByGetCardValue) Len() int           { return len(a) }
func (a ByGetCardValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGetCardValue) Less(i, j int) bool { return a[i].GetCardValue() < a[j].GetCardValue() }

func SortHandByValue(cards []Card) {
	sort.Sort(ByGetCardValue(cards))
}
