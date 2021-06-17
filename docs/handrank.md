## Hand rank representation:  

The rank of a 5 card hand is represented by a 32-bit number, the higher
the number the higher the rank.   THe 32 bit # also encodes the hand details.
The 32 bits contains info about the type of hand (Flush, straight, etc), as well 

The encoding is arranged such that a better hand will have a bigger integer value 
(interpreting the handrank bits as a 32-bit unsigned int with simple BCD encoding) 
 
 
32 bit encoding: 
 
      bits 27-20 kind (high card, pair, ...., Straight Flush)
      bits 19-16 card (2-A)
      bits 15-12 card
      bits 11-8  card
      bits 7-4   card
      bits 3-0   card


Each card in the hand is represented as a 4 bit value.  The least signaficant 20 bits
of the rank represent the 5 cards in order (from left to right) of hand evaluation 
signficance.    I.e.  for a pair, the 2 cards in the pair are the left most 2 cards,
no matter if they are Kings or 2's...


    Example:    hand:   A♥ A♣ K♦ K♣ Q♥   (2 Pair, Ace's and King's)

                (2 pair)   (Ace)  (Ace)  (King) (King) (Queen)
        rank:   00000011   1110   1110   1101   1101   1100




## Hand odds

The Monte Carlo method is used.     

Win and Tie probabilities can be calculated given a pair of hole cards plus a set
of additional cards and a deck.    


