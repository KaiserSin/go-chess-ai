# Weekly Report 2

During the second week I moved from planning to real implementation. The project still follows the DDD and SOLID structure from the architecture document, but now it is not only a plan anymore. The main rule is the same: chess logic must stay inside the domain and must not depend on UI code, framework code, or application code. This gave me a clean base for building the game step by step.

The main work was done in `internal/domain/chess`. At first I implemented the chess domain as one package with full game rules. After that I improved the structure and split the code into logical folders inside the chess domain. Now the root `internal/domain/chess` package is a thin facade. It re-exports the main public types and constructors and stays the main entry point for the rest of the project. The real implementation is now grouped by responsibility, which makes the code easier to read and easier to extend.

The `model` folder contains the main chess value objects and domain types. There I keep `Side`, `PieceType`, `Piece`, `Square`, `Move`, `CastlingRights`, `Status`, `OutcomeReason`, and `Outcome`. This folder also contains the domain errors. The `position` folder contains the board state and almost all chess rules tied to one position. There I implemented `Board`, `Position`, `PositionBuilder`, move validation, move application, move generation, attack detection, castling rules, en passant logic, promotion logic, and draw checks that depend only on the current position. The `game` folder contains the aggregate level flow. There `Game` stores the current position, records repetition history, checks whether the game is finished, and resolves the final outcome after every move.

I also kept low level helpers in private internal folders. In `internal/domain/chess/internal/geom` I put raw board geometry such as movement offsets, pawn direction, back rank, and promotion rank. In `internal/domain/chess/internal/bitboard` I put raw bitboard helpers for iterating and counting bits. This means the domain code is now split in a way that follows logic and responsibility. Small pure helpers are separated from real chess rules, and the public entry point still stays simple.

The chess domain now supports the main rules needed for a playable game. It supports normal legal moves for all pieces, move validation, check detection, checkmate, stalemate, castling, en passant, and pawn promotion with a chosen piece. It also supports draw by threefold repetition, draw by the 50-move rule, and draw by insufficient material. The logic for legal moves works in two steps. First the code generates pseudo legal moves. After that it filters them so the side to move cannot leave its own king in check. Because of this the domain can decide not only if a move is possible, but also whether the game must continue or end.

Another important part of the week was code cleanup. At first the logic worked, but some files were too large and mixed several responsibilities. I simplified this by moving the code into the `model`, `position`, and `game` folders. This made the code closer to DDD and SOLID ideas. `model` holds the core domain values, `position` holds board state and rules, and `game` holds the aggregate root flow. The root `chess` package stays thin and does not contain chess logic anymore. This structure is easier to understand than one large folder with mixed files.

Testing also became a large part of the work. I kept external black box tests in `internal/tests/domain/chess`. These tests check the public behavior of the chess domain through the root `chess` package. I also moved package private white box tests next to the code that owns the logic. Because of this there are now internal tests in `model`, `position`, `internal/geom`, and `internal/bitboard`. This helps keep the tests close to the implementation while the public tests still check the domain from the outside.

I also updated the `Makefile` and the coverage process. After the domain was split into subpackages, a simple one-package coverage command was not enough anymore. I changed the coverage flow so it now checks the whole chess tree and still keeps the strict `100.0%` statement coverage target. This means the new structure did not lower code quality. The project now has a full chess domain with tests and full coverage, and this gives a strong base for the next steps.

Next week I want to build on top of this domain and start the `application/gameplay` layer. After that I want to connect it to the first minimal playable board and UI. The chess rules are now ready, so the next goal is to make the application use them in a clean way and show the game on screen.

## Time tracking

| Day | Hours | Description |
| Sunday | 3 | first chess version |
| Saturday | 10 | all chess logics + tests |

