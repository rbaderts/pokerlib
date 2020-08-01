package pokerlib

import (
	"fmt"
	"sort"
)

type HandRank int

const HAND_RANK_MASK = 0xFFF80000

const (
	HighCard      HandRank = 0x0
	Pair          HandRank = 0x1
	TwoPair       HandRank = 0x2
	ThreeOfAKind  HandRank = 0x3
	Straight      HandRank = 0x4
	Flush         HandRank = 0x5
	FullHouse     HandRank = 0x6
	FourOfAKind   HandRank = 0x7
	StraightFlush HandRank = 0x8
)

func GetHandRank(rank int) HandRank {
	fmt.Printf("GetHandRank: rank = %b\n", rank)
	return HandRank(rank&HAND_RANK_MASK) >> 24
}

func (this HandRank) String() string {
	fmt.Printf("HandRank.String() = %b\n", this)
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

type RankedHand struct {
	rank  HandRank
	value AbsoluteValue
}

/*
Hand Rankings:

      HighCard:   2-14
      Pair:    2*10 - 14*10    (20-140)
      TwoPair:  140 + (10*HiPairC) +  HiPairC    (
                  163-293

8's & 7' = 140+8*10+7 = 227,
2's & 3's:    140 + 20 + 3 = 163
9's & 2's = 140 + 9*10 * 2 = 252
A's & K's = 140 + 14*10 + 13 = 293

     ThreeOfAKind:   300 + C:     302-314
     Straight:       320 + HC:   320-334
      Flush:         340 + HC:    340 - 354
      FullHouse:     360 + (10 * TripC) + PairC     (391-513)

222 - 33 =   360 + 2*14 + 3 = 391
AAA - KK =   360 + 10*14 + 13 = 513

      FourOfAKind:    520 + C :     522-534




*/

/*
HandRank:
    THe base hand rank is High Card -> straight Flush
     The specific rank of hands of the same base rank are 0 - N (0 being the lowest
         rank for that base hand rank, N being the highest.

   THe overall rank is 100 * the base hand rank value (1 for High Card 9 for straight flush
                 + The specific rank.    Thus all Trip sets for example would have a rank of 400-412
*/
/*
type Hand struct {
	HandRank  string
	RankValue int
	/* 0 - N number where N is the relative ranking for the hand.

	          trip(2) - rankvalue = 0
	          Straight(A-4) - rankvalue = 0
	          Straight(10-A) - rankvalue = 13
	          Flush(A high) - rank value = 8
	          Flush(6 high) - rank value = 0
	          FullHouse(AAA-KK) = top
	          FullHouse(222-33) = bottom
	             FullHouse: OversRank * 13 + UndersRank.
	    (222-33) 0*13 + (1) = 1
	   (222-AA) 0*13 + 12)) = 12
	   (333-22) 1*13 + 10)) = 13

	Rank int
}
*/

func Rank(cards []Card) int {

	topRank := 0
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
	fmt.Printf("hand value = %d\n", r)
	return r

}

// 5 cards

/*

   absolute hand ranking is represented by an int64.
   Higher hand ranks havea  bigger absolute hand rank
      absolute hand rank.
        23-20  hand type  (1010 = straight flush : 0000 = high card)
        19-16 highest ranked card
        15-12 second highest ranked
        11-8  third highest ranked
        7-4   fourth highest ranked
        3-0   fifth highest ranked

*/

func RankHand(cards []Card) int {

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

	isFlush := isFlush(cards)
	isStraight := isStraight(cards)
	value := 0

	if isStraight && isFlush {
		copy(sortedCards, cards)
		SortCards(sortedCards)

		value = int(StraightFlush) << 24
		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}
		return value
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

		value = int(FourOfAKind) << 24
		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}

		return value

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
		value = int(FullHouse) << 24
		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}

	} else if isFlush {
		copy(sortedCards, cards)
		SortCards(sortedCards)

		value = int(Flush) << 24
		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}
		return value

	} else if isStraight {
		copy(sortedCards, cards)
		SortCards(sortedCards)

		value = int(Straight) << 24
		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
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

		value = int(ThreeOfAKind) << 24

		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}

	} else if topCount == 2 && secondCount == 2 {
		//copy(sortedCards, cards)

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
		value = int(TwoPair) << 24
		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
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

		value = int(Pair) << 24

		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}
	} else {
		copy(sortedCards, cards)
		SortCards(sortedCards)
		value = int(HighCard) << 24

		for i, c := range sortedCards {
			value |= (int(c.AbsoluteValue()) << (i * 4))
		}
	}

	//highCard := getHighCard(cards)

	return value

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
	min := 100
	max := 0
	for _, c := range cards {
		if int(c.Index) == min {
			return false
		}
		if int(c.Index) == max {
			return false
		}
		if int(c.Index) < min {
			min = int(c.Index)
		}
		if int(c.Index) > max {
			max = int(c.Index)
		}
	}
	if max-min == 4 {
		return true
	}
	return false
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
func (this Hand) AbsoluteValue() AbsoluteValue {

	h := make([]Card, len(this))
	fmt.Printf("hand = %v\n", this)
	copy(h, this)
	SortHandByValue(h)
	fmt.Printf("sorted hand = %v\n", h)

	var value int64
	value = 0
	for i, c := range h {
		value |= (int64(c.AbsoluteValue()) << (i * 6))
	}

	return AbsoluteValue(value)
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

type ByAbsoluteValue []Card

func (a ByAbsoluteValue) Len() int           { return len(a) }
func (a ByAbsoluteValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAbsoluteValue) Less(i, j int) bool { return a[i].AbsoluteValue() < a[j].AbsoluteValue() }

func SortHandByValue(cards []Card) {
	sort.Sort(ByAbsoluteValue(cards))
}
