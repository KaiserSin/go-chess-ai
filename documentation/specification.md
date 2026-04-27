# Specification Document

## Programming Language

The project is written in Go.
The documentation language is English.

## Course Information

Study programme is tietojenkäsittelytieteen kandidaatti (TKT).

Programming languages suitable for peer review are Go.

## Problem Statement

The project implements a chess game with an artificial intelligence opponent.
The bot must evaluate positions, calculate future moves, and choose legal responses against a human player.
All chess logic is implemented in the project without external chess libraries.

## Algorithms and Data Structures

The main search is a minimax-style game tree search with alpha-beta pruning.
The AI uses time-limited iterative deepening and searches as deeply as it can within the configured move budget.
The current search also uses aspiration windows, quiescence search, and move ordering.
The evaluation function considers material and positional factors such as piece placement, pawn structure, king safety, and simple endgame bonuses.

The project uses the following data structures

- bitboards based on `uint64` for efficient board representation
- a game tree explored through recursive search
- Go slices for storing generated legal moves and root move results

## Program Inputs

The program uses a graphical desktop interface.
The user selects pieces and target squares on the board with the mouse.
The program validates the move, updates the board state, and then lets the AI calculate and apply its reply.

## Expected Time and Space Complexity

The time complexity of plain Minimax is `O(b^d)`, where `b` is the branching factor and `d` is the completed search depth reached within the time budget.
Alpha-beta pruning keeps the same worst-case order, but it reduces the amount of explored work in practice and in the best case improves the search toward `O(b^(d/2))`.

Move ordering, aspiration windows, and iterative deepening are used to make the practical search more efficient and stable.
They do not remove the exponential worst-case shape of the search tree.
Quiescence search extends the search in tactical positions, so its practical cost depends on the number of available captures and promotions.
The static evaluation function scans the fixed-size chess board and runs in `O(64)`, which is effectively `O(1)` for normal chess.

The space complexity is `O(d)`.
The main memory cost comes from the recursion stack because the search stores the current line rather than the full tree.

## Core of the Project

The core of the project is the chess AI search and evaluation.
Correct move generation and the desktop interface are necessary for the full application, but the main technical focus is the search for the best move and the evaluation of board positions.

## Sources

- www.chessprogramming.org
- www.chessprogramming.org/Alpha-Beta
- www.chessprogramming.org/Move_Ordering
- www.chessprogramming.org/Iterative_Deepening
- www.chessprogramming.org/Quiescence_Search
- en.wikipedia.org/wiki/Minimax
- en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
