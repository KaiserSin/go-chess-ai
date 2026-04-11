# Implementation document

This document describes the project in its current state.
The project is a chess program written in Go.
It has full chess rules, a desktop user interface built with Ebiten, and a human versus AI mode.

At the moment the program supports

- full move validation and game state handling
- check, checkmate, stalemate, castling, en passant, promotion, and draw rules
- a desktop chess board with mouse interaction
- a start menu where the user chooses a side and the AI depth
- AI depth selection in the range `1..5`
- an AI based on alpha-beta search, iterative deepening, a transposition table, root-level goroutines, quiescence search, and a simple hand-tuned evaluation function

## Project structure

The program is split into clear layers and packages.
The dependency direction is

`presentation -> application -> domain`

Main parts of the program

- `cmd/chess-desktop` starts the desktop program
- `internal/domain/chess` contains the chess rules, board state, move generation, move validation, and game results
- `internal/application/gameplay` connects the UI to the chess domain and handles square selection, move attempts, promotion flow, snapshots, and AI move requests
- `internal/application/ai` contains the search and evaluation code. It includes alpha-beta search, iterative deepening, transposition table use, quiescence search, move ordering, and position evaluation
- `internal/presentation/ebiten` contains the desktop UI, rendering, input translation, board mapping, and theme values
- `internal/infrastructure/bootstrap` wires the application together and starts the desktop program

This structure keeps the chess rules separate from the UI code.
Because of this, the core game logic can be tested without the rendering layer.

## Time and space complexities

The most important algorithms and operations in the project are listed below.

- Evaluation function runs in `O(64)`. On a fixed chess board this is effectively `O(1)`. The function scans the board and calculates material and positional terms
- Move ordering runs in `O(m)`, where `m` is the number of legal moves in the current position. It groups moves into categories such as hash move, promotions, captures, and quiet moves
- Transposition table probe and store are on average `O(1)`. The table is implemented as a hash map for one search call
- Alpha-beta search has worst-case time complexity `O(b^d)`, where `b` is the branching factor and `d` is the search depth. In practice, move ordering and pruning reduce the number of visited nodes
- Iterative deepening has the same asymptotic order as the deepest search, but it repeats some work from shallow searches. In practice it helps because shallow searches improve move ordering and transposition table use
- Quiescence search does not have one simple fixed bound in practice. Its cost depends on the number of tactical continuations. It increases search work, but reduces simple tactical blunders

The main space costs are

- recursion stack with complexity `O(d)`
- transposition table with complexity `O(t)`, where `t` is the number of stored positions during one search

So the total search memory use is approximately

- `O(d + t)`

## Performance and O-analysis comparison

This comparison is mostly theoretical.
The project does not yet have a separate benchmark document or benchmark suite.

- Plain minimax searches all branches up to depth `d`. Alpha-beta prunes many branches. The worst-case order is still `O(b^d)`, but alpha-beta is much faster in practice
- Without move ordering, alpha-beta still works, but pruning is weaker. With move ordering, stronger moves are searched first, so cutoffs happen earlier and the running time is better
- Without a transposition table, repeated positions are searched again. With a transposition table, already analyzed positions can be reused, which reduces repeated work
- Without quiescence search, the evaluation can stop in unstable tactical positions. With quiescence search, the engine continues searching tactical moves and avoids some simple blunders, even though it increases the amount of work
- Root-level goroutines do not change the asymptotic complexity, but they can reduce wall-clock time on multi-core hardware

## Possible flaws and improvements

There are still several clear limitations in the project.

The AI uses a hand-tuned evaluation function.
It is better than a pure material count, but it is still simple.
The program does not yet use an opening book, endgame tablebases, killer heuristic, history heuristic, or time-based search.
Because of this, the AI is still much weaker than a mature chess engine.

The current search depth is intentionally limited to `5` in the user interface.
This makes the program easier to use and keeps response time reasonable, but it also limits playing strength.

The project also does not yet have a dedicated benchmark document.
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

The main model used was an Gemini.
It was used for

- checking the quality of written text
- giving hints about possible improvement directions
- explaining search algorithms in simple terms

The final code and documentation were still checked in the repository and validated locally.
At the time of writing, the project passes

```text
go test ./...
```

## References

Only sources that had real relevance for the project are listed here.

- Chessprogramming Wiki https://www.chessprogramming.org/Minimax
- Chessprogramming Wiki https://www.chessprogramming.org/Alpha-Beta
- Chessprogramming Wiki https://www.chessprogramming.org/Move_Ordering
- Chessprogramming Wiki https://www.chessprogramming.org/Transposition_Table
- Chessprogramming Wiki https://www.chessprogramming.org/Iterative_Deepening
- Chessprogramming Wiki https://www.chessprogramming.org/Quiescence_Search
- Wikipedia https://en.wikipedia.org/wiki/Minimax
- Wikipedia https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning
