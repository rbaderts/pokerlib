package pokerlib

import ()

type Hands map[int][2]Card

type Odds struct {
	Wins float32 `json:"wins"`
	Ties float32 `json:"ties"`
}

/**

  Computes win & ties odds given a set of hand cards and a set of common cards

*/
func CalculateOdds(deck *Deck, hands Hands, commonCards []Card, depth int) map[int]*Odds {

	//wins := make(map[int]int, 0)
	//ties := make(map[int]int, 0)

	results := make(map[int]*Odds, 0)
	for k, _ := range hands {
		results[k] = new(Odds)
	}

	//result := make(map[int]float32)

	cardsToDraw := 5 - len(commonCards)
	//	tested := make([][]Card, 0)

	fiveCards := make([]Card, 5)
	for i := 0; i < len(commonCards); i++ {
		fiveCards[i] = commonCards[i]
	}
	startIndex := len(commonCards)

	handsEvaluated := 0
	//total := 10000

	workingDeck := deck.Copy()

	checkedHands := make(map[uint]bool, 0)
	for count := 0; count < depth; count++ {

		cards := getNCards(workingDeck, cardsToDraw)
		h := NewCardSet(cards)
		hash := h.Hash()
		_, ok := checkedHands[hash]
		if ok {
			for _, c := range cards {
				workingDeck.ReturnCard(c)
			}
			continue // already checked
		}
		checkedHands[hash] = true

		//fmt.Printf("Evaulating with additional cards: %v\n", cards)
		//fmt.Printf("deck now; %v\n", workingDeck)
		for i := 0; i < len(cards); i++ {
			fiveCards[startIndex+i] = cards[i]
		}

		topRankedHands := make([]int, 0)
		topRank := uint(0)
		for k, h := range hands {
			handCards := make([]Card, 7)
			handCards[0] = h[0]
			handCards[1] = h[1]
			handCards[2] = fiveCards[0]
			handCards[3] = fiveCards[1]
			handCards[4] = fiveCards[2]
			handCards[5] = fiveCards[3]
			handCards[6] = fiveCards[4]

			_, rank := Rank(handCards)
			if uint(rank) > uint(topRank) {
				topRank = uint(rank)
				topRankedHands = make([]int, 1)
				topRankedHands[0] = k
			} else if uint(rank) == uint(topRank) {
				topRankedHands = append(topRankedHands, k)
			}

		}

		if len(topRankedHands) == 1 {
			results[topRankedHands[0]].Wins = results[topRankedHands[0]].Wins + 1
		} else {
			for _, k := range topRankedHands {
				results[k].Ties = results[k].Ties + 1
			}
		}

		for _, c := range cards {
			workingDeck.ReturnCard(c)

		}
		handsEvaluated += 1

	}

	for _, result := range results {
		result.Wins = result.Wins / float32(handsEvaluated)
		result.Ties = result.Ties / float32(handsEvaluated)
	}

	return results

}

func getNCards(deck *Deck, n int) []Card {

	result := make([]Card, n)

	for i := 0; i < n; i++ {
		result[i] = deck.BorrowRandom()
	}
	//fmt.Printf("The Deck is now\n")
	//fmt.Printf("deck : %v\n", deck)
	///deck.Restore()
	return result

}
