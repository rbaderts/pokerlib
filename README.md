# pokerlib

A library of poker primitives in go

Features:

   * Foundation data types for Card, Hand and Deck
   * Hand ranking.  Calculates best 5 card hand out of any set of cards.
   * Hand rank encoding 
   * English language hand descriptions
        Ex:   Full house, Kings full of 8's,  An Ace high flush, etc.]]
   * Hand evaluation logging
   * Odss calculations

## Rank representation:  

Each 5 card hand has a unique 32-bit encoding, this value is called the Hand Rank.  
The 32 bits contains info about the type of hand (Flush, straight, etc), as well 
as the index of each of the 5 cards.  The encoding is arranged such that a better 
hand will have bigger value (interpreting the handrank bits as a 32-bit unsigned 
int with simple BCD encoding) better hand will always have a higher  
 
 
32 bit encoding: 
 
      bits 27-20 contains the hand kind 
      bits 19-16 card
      bits 15-12 card
      bits 11-8  card
      bits 7-4   card
      bits 3-0   card

Each card in the hand is represented as a 4 bit value.  The least signaficant 20 bits
of the rank represent the 5 cards in order of hand evaluation signficance.    I.e.
for a pair, the 2 cards in the pair are the left most 2 cards, no matter if they 
are Kings or 2's...


    Example:    hand:   A♥ A♣ K♦ K♣ Q♥   (2 Pair, Ace's and King's)

                (2 pair)   (Ace)  (Ace)  (King) (King) (Queen)
        rank:   00000011   1110   1110   1101   1101   1100




## Hand odds

Win and Tie probabilities can be calculated given a set of hand hold cards, 
an additional set of 0-5 common cards and a deck.

The Monte Carlo method is used.     
