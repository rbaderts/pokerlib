package pokerlib

import (
	"fmt"
	"sort"
)

/*

   HandRank is the absolute value of the hand.    It is encoded as follows:

      4 bits per card index.   So first least significant HAND_KIND_SHIFT bit is room for the indexes
         of all 5 cards in a hand.  However for pairs, trips, two pairs on the cards
         not in any of those sets are present.     So a KK78T   would have cards:   7 8 T

      The next 4 bits:   HAND_KIND_SHIFT-23 encode the 4 bit hand kind

        Ex:    For the hand K 7 8 4 2 (not-suited).   The value would be:

*/

type HandRank int

const HAND_KIND_MASK = 0xFFF80000
const HAND_KIND_SHIFT = 20

type HandKind int8 // Pair, Straight, etc.

const (
	HighCard      HandKind = 0b00000001
	Pair          HandKind = 0b00000010
	TwoPair       HandKind = 0b00000011
	ThreeOfAKind  HandKind = 0b00000100
	Straight      HandKind = 0b00000101
	Flush         HandKind = 0b00000110
	FullHouse     HandKind = 0b00000111
	FourOfAKind   HandKind = 0b00001000
	StraightFlush HandKind = 0b00001001
)

func GetHandKind(rank HandRank) HandKind {
	fmt.Printf("GetHandRank: rank = %b\n", rank)
	return HandKind(rank&HAND_KIND_MASK) >> HAND_KIND_SHIFT
}

func (this HandKind) String() string {
	fmt.Printf("HandKind.String() = %b\n", this)
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
	return r.String()

}
/**
 */
func Rank(cards []Card) HandRank {

	fmt.Printf("Cards = %v\n", cards)
	var topRank HandRank = 0
	if len(cards) > 5 {
		for i, _ := range cards {
			var subset = make([]Card, len(cards)-1)
			count := 0
			for j, c := range cards {
				if j != i {
					subset[count] = c
					count += 1
				}
			}
			r := Rank(subset)
			if r > topRank {
				topRank = r
			}
		}
		return topRank
	}
	r := RankHand(cards)
	fmt.Printf("hand value = %.24b\n", r)
	return r

}

/*
   Called only for 5 cards
 */
func RankHand(cards []Card) HandRank {

	sortedCards := make([]Card, 5)
	// Histogram is a map of ranks to frequency
	histos := buildHistogram(cards)
	fmt.Printf("histos = %v\n", histos)

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
	fmt.Printf("isStraight = %v\n", isStraight)
	value := 0

	if isStraight && isFlush {
		value = int(StraightFlush) << HAND_KIND_SHIFT
		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}

	} else if topCount == 4 {
		mainCard := histos[0].Index
		i := 0
		for _, c := range cards {
			if c.Index == mainCard {
				sortedCards[i] = c
				i += 1
			} else {
				sortedCards[4] = c
			}
		}

		value = int(FourOfAKind) << HAND_KIND_SHIFT
		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}

	} else if topCount == 3 && secondCount == 2 {
		topCard := histos[0].Index

		i := 0
		j := 3
		for _, c := range cards {
			if c.Index == topCard {
				sortedCards[i] = c
				i += 1
			} else {
				sortedCards[j] = c
				j += 1
			}
		}
		value = int(FullHouse) << HAND_KIND_SHIFT
		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))

		}

	} else if isFlush {

		value = int(Flush) << HAND_KIND_SHIFT
		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}
		return HandRank(value)

	} else if isStraight {

		value = int(Straight) << HAND_KIND_SHIFT
		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}

	} else if topCount == 3 {
		topCard := histos[0].Index

		i := 0
		j := 0

		extraCards := make([]Card, 2)
		for _, c := range cards {
			if c.Index == topCard {
				sortedCards[i] = c
				i += 1
			} else {
				extraCards[j] = c
				j += 1
			}
		}

		for i, c := range extraCards {
			sortedCards[3+i] = c
		}

		value = int(ThreeOfAKind) << HAND_KIND_SHIFT

		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}

	} else if topCount == 2 && secondCount == 2 {

		topCard := histos[0].Index
		nextCard := histos[1].Index
		i := 0
		j := 2
		for _, c := range cards {
			if c.Index == topCard {
				sortedCards[i] = c
				i += 1
			} else if c.Index == nextCard {
				sortedCards[j] = c
				j += 1
			} else {
				sortedCards[4] = c
			}
		}

		value = int(TwoPair) << HAND_KIND_SHIFT
		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}
	} else if topCount == 2 {
		topCard := histos[0].Index
		i := 0
		j := 0
		extraCards := make([]Card, 3)
		for _, c := range cards {
			if c.Index == topCard {
				sortedCards[i] = c
				i += 1
			} else {
				extraCards[j] = c
				j += 1
			}
		}
		for i, c := range extraCards {
			sortedCards[2+i] = c
		}

		value = int(Pair) << HAND_KIND_SHIFT

		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}
	} else {
		copy(sortedCards, cards)
		SortCards(sortedCards)
		value = int(HighCard) << HAND_KIND_SHIFT

		for i, c := range sortedCards {
			value |= (int(c.GetCardValue()) << ((len(sortedCards) - i -1) * 4))
		}
	}

	//highCard := getHighCard(cards)

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
	fmt.Printf("isStraight called\n")

	if len(cards) != 5 {
		fmt.Printf("Error isStraight called on a set of more than 5 cards\n")
		return false
	}
	low := 0
	last := 0
	for i, c := range cards {

		if i == 0 {
			fmt.Printf("i == 0\n")
			index := c.Index
			if index == Ace {
				low = 0
			} else {
				low = int(c.Index)
			}
			last = low
			fmt.Printf("low == %d\n" ,low)
		} else {
			if int(c.Index) != last+1 {
				return false
			}
			last += 1
		}
	}
	return true
	/*
	min := 100
	max := -4
	hasAce := false
	for _, c := range cards {
		index := c.Index
		if index == Ace {
			if hasAce == true {
				return false;
			}
			hasAce = true
			index := -1
		}
		if int(index) == min {
			return false
		}
		if int(index) == max {
			return false
		}
		if int(index) < min {
			min = int(index)
		}
		if int(index) > max {
			max = int(index)
		}
	}
	fmt.Printf("max = %d, min = %d\n", max, min)


	if max-min == 4 {
		return true
	}
	return false
	*/
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

/*
func (this Hand) GetCardValue() GetCardValue {

	h := make([]Card, len(this))
	fmt.Printf("hand = %v\n", this)
	copy(h, this)
	SortHandByValue(h)
	fmt.Printf("sorted hand = %v\n", h)

	var value int64
	value = 0
	for i, c := range h {
		value |= (int64(c.GetCardValue()) << (i * 6))
	}

	return GetCardValue(value)
}
*/

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
