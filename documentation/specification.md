# Specification Document

## Programming Language

The project is written in Go.

## Problem Statement

The project implements a chess game with an artificial intelligence opponent.
The bot must evaluate positions, calculate future moves, and choose legal responses against a human player.
All chess logic is implemented in the project without external chess libraries.

## Algorithms and Data Structures

The main search algorithm is Minimax with alpha-beta pruning.
The project uses the following data structures:

- bitboards based on `uint64` for efficient board representation
- a game tree explored through recursive search
- Go slices for storing generated legal moves

## Program Inputs

The program uses a graphical desktop interface.
The user selects pieces and target squares on the board with the mouse.
The program validates the move, updates the board state, and then lets the AI calculate and apply its reply.

## Expected Time and Space Complexity

The time complexity of plain Minimax is `O(b^d)`, where `b` is the branching factor and `d` is the search depth.
Alpha-beta pruning reduces the amount of explored work in practice and in the best case improves the search toward `O(b^(d/2))`.

The space complexity is `O(d)`.
The main memory cost comes from the recursion stack because the search stores the current line rather than the full tree.

## Core of the Project

The core of the project is the chess AI search and evaluation.
Correct move generation and the desktop interface are necessary for the full application, but the main technical focus is the search for the best move and the evaluation of board positions.

## Sources

- https://www.chessprogramming.org
- https://en.wikipedia.org/wiki/Minimax
- https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
