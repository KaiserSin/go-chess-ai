# Implementation document

The project is a chess program written in Go.
It has full chess rules, a desktop user interface built with Ebiten, and a human versus AI mode.

## Features

The program supports

- full move validation and game state handling
- check, checkmate, stalemate, castling, en passant, promotion, and draw rules
- a desktop chess board with mouse interaction
- a start menu where the user chooses a side
- a fixed AI search depth of `3`
- an AI based on alpha-beta search, iterative deepening, root-level goroutines, quiescence search, and a simple hand-tuned evaluation function

## Project structure

The program is split into clear layers and packages.
The dependency direction is

`presentation -> application -> domain`

Main parts of the program

- `cmd/chess-desktop` starts the desktop program
- `internal/domain/chess` contains the chess rules, board state, move generation, move validation, and game results
- `internal/application/gameplay` connects the UI to the chess domain and handles square selection, move attempts, promotion flow, snapshots, and AI move requests
- `internal/application/ai` contains the search and evaluation code. It includes alpha-beta search, iterative deepening, quiescence search, move ordering, and position evaluation
- `internal/presentation/ebiten` contains the desktop UI, rendering, input translation, board mapping, and theme values
- `internal/infrastructure/bootstrap` wires the application together and starts the desktop program

This structure keeps the chess rules separate from the UI code.
Because of this, the core game logic can be tested without the rendering layer.

## Time and space complexities

The most important algorithms and operations in the project are listed below.

- Evaluation function runs in `O(64)`. On a fixed chess board this is effectively `O(1)`. The function scans the board and calculates material and positional terms
- Move ordering runs in `O(m)`, where `m` is the number of legal moves in the current position. It groups moves into categories such as promotions, captures, and quiet moves
- Alpha-beta search has worst-case time complexity `O(b^d)`, where `b` is the branching factor and `d` is the search depth. In practice, move ordering and pruning reduce the number of visited nodes
- Iterative deepening has the same asymptotic order as the deepest search, but it repeats some work from shallow searches. In practice it still helps because shallow searches can make the final move choice more stable
- Quiescence search does not have one simple fixed bound in practice. Its cost depends on the number of tactical continuations. It increases search work, but reduces simple tactical blunders

The main space cost is

- recursion stack with complexity `O(d)`

## Performance and O-analysis comparison

This comparison is mostly theoretical.
The project does not have a separate benchmark document or benchmark suite.

- Plain minimax searches all branches up to depth `d`. Alpha-beta prunes many branches. The worst-case order is still `O(b^d)`, but alpha-beta is much faster in practice
- Without move ordering, alpha-beta still works, but pruning is weaker. With move ordering, stronger moves are searched first, so cutoffs happen earlier and the running time is better
- Without quiescence search, the evaluation can stop in unstable tactical positions. With quiescence search, the engine continues searching tactical moves and avoids some simple blunders, even though it increases the amount of work
- Root-level goroutines do not change the asymptotic complexity, but they can reduce wall-clock time on multi-core hardware

## Possible flaws and improvements

There are still several clear limitations in the project.

The AI uses a hand-tuned evaluation function.
It is better than a pure material count, but it is still simple.
The program does not use an opening book, endgame tablebases, killer heuristic, history heuristic, or time-based search.
Because of this, the AI is still much weaker than a mature chess engine.

The current search depth is intentionally fixed at `3`.
This keeps response time predictable and simplifies the menu flow, but it also limits playing strength.

The project also does not have a dedicated benchmark document.
Complexity analysis can be given, but there is still no proper measured comparison between different search configurations or evaluation versions.

Possible future improvements

- add an opening book
- add time-based search and time management
- add killer and history heuristics
- improve the evaluation function with more positional terms
- create benchmark measurements for different depths and search variants
- add endgame tablebases or more endgame-specific knowledge

## Use of language models

Large language models were used in this project.

A Gemini model was used for

- checking the quality of written text
- giving hints about possible improvement directions
- explaining search algorithms in simple terms

The final code and documentation were still checked in the repository and validated locally.
The project passes

```text
go test ./...
```

## References

Only sources that had real relevance for the project are listed here.

- Chessprogramming Wiki https://www.chessprogramming.org/Minimax
- Chessprogramming Wiki https://www.chessprogramming.org/Alpha-Beta
- Chessprogramming Wiki https://www.chessprogramming.org/Move_Ordering
- Chessprogramming Wiki https://www.chessprogramming.org/Iterative_Deepening
- Chessprogramming Wiki https://www.chessprogramming.org/Quiescence_Search
- Wikipedia https://en.wikipedia.org/wiki/Minimax
- Wikipedia https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
