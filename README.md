# pokerlib

A library of poker primitives in go

Features:

   * Basic types for Cards, Hands, Decks
   * Complete hand ranking.  Calculates best 5 card hand given any number of cards.
   * Hand rank encoding 
   * English language hand descriptions
        Ex:   Full house, Kings full of 8's,  An Ace high flush, etc.]]
   * Hand evaluation logging
   * Todo:  Partial hand win probabilities

## Rank representation:  

Each 5 card hand has a unique encoding that fits into
32 bits, this value is called the Hand Rank.  
The 32 bits contains info about the type of hand (Flush,
straight, etc), as well as the index of each of the 5 cards.
The encoding is arranged such that a better hand will have 
bigger value (interpreting the handrank bits as a 32-bit
unsigned int with simple BCD encoding)
better hand
will always have a higher  
 
 
 The 32 bit value
taken as a number 
 five card hand is encodedEach hand rank is encoded into a 4 byte value.   A higher ranking hand's 
4 byte value will have a higher absolute value than a lower ranking hand (when
interpreted as an unsigned integer).   

Each hand type, high-card thru straight flush is represented by an 8 bit value.   
The hand type is shifted to be the most signifcant bits in the rank, thus the hand 
type determines the better hand no matter the other 20bits bytes.    

For hands with the same type the other 20 bits come into play.

Each card in the hand is represented as a 4 bit value.  The least signaficant 20 bits
of the rank represent the 5 cards in order of hand evaluation signficance.    I.e.
for a pair, the 2 cards in the pair are the left most 2 cards, no matter if they are Kings or 2's...



    Example:    hand:   A♥ A♣ K♦ K♣ Q♥   (2 Pair, Ace's and King's)

                (2 pair)   (Ace)  (Ace)  (King) (King) (Queen)
        rank:   00000011   1110   1110   1101   1101   1100

