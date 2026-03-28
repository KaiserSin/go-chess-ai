# Testing Document

This document describes the current state of automated testing in the project.
It is intended to be updated during the course and to give a short overview of coverage and the most important implemented tests at each stage of development.

## Test Coverage Report

```text
Command: ./scripts/coverage.sh
Coverage scope: internal/domain/chess (merged domain coverage report)
total:                                              (statements)            100.0%
```

## Current Tests

The current automated tests are unit tests. They verify the domain logic of the chess engine, the basic behavior of the gameplay layer, and small helper functions in the presentation layer.

- Verified that the initial gameplay snapshot contains 64 squares, the correct side to move, and the expected starting pieces on selected squares.
- Verified that selecting a piece marks the selected square and legal target squares correctly, and that selecting the same square again clears the selection.
- Verified that a legal move updates the board state, changes the side to move, and clears the previous selection state.
- Verified that a promotion move opens four promotion options in the correct order and that choosing a promotion piece updates the final board state correctly.
- Verified that the position builder rejects invalid sides, invalid squares, invalid counters, invalid piece types, and illegal positions.
- Verified that square creation and square parsing accept valid coordinates and reject invalid coordinates.
- Verified that legal move validation rejects illegal and pinned moves.
- Verified that castling is accepted only in legal situations, moves the king and rook to the correct squares, and updates castling rights when the king or rook moves.
- Verified that en passant is accepted only in the correct situations and expires after the allowed move window.
- Verified that pawn promotion requires a valid promotion piece and rejects invalid promotion attempts.
- Verified that check, checkmate, and stalemate produce the correct game status and outcome.
- Verified that threefold repetition, the fifty-move rule, and insufficient material produce the correct draw outcome.
- Verified that text conversion methods for side, piece type, piece, square, move, status, and outcome reason return the expected string values.
- Verified that piece sprite loading finds all 12 current piece images and does not return a sprite for an unknown key.
- Verified that input translation maps screen coordinates to the correct board squares and promotion choices, and rejects clicks outside the board.
- Verified that the piece catalog returns the correct labels for known piece keys and a safe fallback label for unknown keys.
- Verified that the board view model keeps the board orientation correct, uses the correct square grid coordinates, highlights selection and legal targets, and places promotion options in the expected order.
- Verified that sprite placement scales piece images correctly to the size of one board square.
