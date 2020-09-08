# pokerlib
A library of poker primitives in go

Features:

   * Basic types for Cards, Hands, Decks
   * Complete hand ranking.   Determines best 5 card hand given any # of cards
   * <32 byte encoding of a 5 card hand where bigger value mean bigger hand
   * Support for english descriptions of hands.    
        Ex:   Full house, Kings full of 8's,  An Ace high flush, etc.]]

Implements ranking of 5 or more card hands.    
Each hand rank is represented as a 32-byte value.   A higher 5 card poker
hand will have a greater value number than a lower 5 card poker hand.



