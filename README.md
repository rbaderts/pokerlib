# pokerlib
A library of poker primitives in go

Features:

   * Basic types for Cards, Hands, Decks
   * Calculates best 5 card hand given 7 cards.
   * <32 byte encoding of hand ranks who's unsigned value preserve ranking order
   * Support for english descriptions of hands
   * 

Implements ranking of 5 or more card hands.    
Each hand rank is represented as a 32-byte value.   A higher 5 card poker
hand will have a greater value number than a lower 5 card poker hand.



