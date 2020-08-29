package pokerlib

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestDumpDeck(t *testing.T) {

	deck := NewDeck()
	deck.Shuffle()

	file, _ := json.MarshalIndent(deck, "", " ")

	_ = ioutil.WriteFile("deck.json", file, 0644)

}
