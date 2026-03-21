# Weekly Report 2

During the second week I moved from planning to real implementation. The project now follows a DDD and SOLID style structure with separate `domain`, `application`, `presentation`, and `infrastructure` parts. The main idea is still the same. Chess rules must stay inside the domain and must not depend on UI code or framework code. This gave me a clear base for the next steps of the project.

The main practical work was done in `internal/domain/chess`. In this folder I implemented the core chess domain. Now the code is split into smaller files with clearer responsibilities. In `side.go`, `piece.go`, `state.go`, and `outcome.go` I keep the main domain types such as side, piece type, piece, castling rights, status, and game outcome. In `square.go` I implemented board coordinates, square parsing, and square formatting. In `move.go` I added the move model with source square, target square, and optional promotion piece. In `board.go` I implemented board storage with hidden bitboards, piece placement, piece removal, and piece lookup. I kept bitboards because they are a good fit for future AI work.

The biggest part of the work was in `position.go`, `movegen.go`, and `game.go`. In `position.go` I implemented the full state of a chess position: the board, side to move, castling rights, en passant square, halfmove clock, and fullmove number. The same file applies moves and validates position state. In `movegen.go` I implemented move generation for all chess pieces. The logic first builds pseudo legal moves and then filters them so the king is never left in check. In `game.go` I made `Game` the main entry point for the domain. It stores the current `Position`, applies moves, tracks repetition history, checks if the game is finished, and updates the final outcome after every move.

The domain now supports the main rules needed for a real playable chess game. It supports normal legal moves for every piece, check detection, move validation, checkmate, stalemate, castling, en passant, and pawn promotion with a selected piece. I also implemented the main draw rules. In `draw.go` I implemented insufficient material detection. In `history.go` I added position history tracking for threefold repetition. The game also supports the 50-move rule. Because of this the current chess domain can already decide if the game should continue or if it must end in win or draw.

I also added support code around the domain. In `builder.go` I implemented a `PositionBuilder` that makes it easier to create custom chess positions for tests. In `errors.go` I added domain errors for bad moves and bad positions. In `layout.go` I put shared board helpers and named squares that are used in pawn movement and castling logic. This means the project is not only planned anymore. The core chess logic is already implemented in code and follows the architecture that was defined earlier.

Later I also simplified the code for readability. The old large mixed type file was split into smaller files. Repeated side and square logic was moved into helpers. Move validation and castling logic became more direct. The goal of this cleanup was to keep all chess rules correct but make the code easier to read and easier to extend.

Another big part of the work was testing. I created two layers of tests. The external black box tests are in `internal/tests/domain/chess` and they check the public behavior of the chess domain. I split them by topic so they are easier to read: text and value objects, position and move validation, castling, en passant, promotion, status and end states, repetition, draw rules, and builder validation. I also added internal white box tests in `internal/domain/chess` for private helpers and hard to reach branches. The test names use simple English so they are easier to understand.

I also worked on code quality and verification. I added a `Makefile` with commands for normal tests, coverage, and HTML coverage output. The chess domain now has `100.0%` statement coverage. This does not mean the project is finished, but it shows that the core domain logic was tested very carefully. It also gives me a strong base before I move to the application layer and the user interface.

Next week I will build on this domain and start the `application/gameplay` layer. After that I want to make the first minimal playable board and UI. The goal is to connect the already finished chess logic to a simple interface so the game can be played on screen. When this is ready I will have a strong base for the next step, which is the chess bot and search algorithm.

## Time tracking

| Day | Hours | Description |
| Sunday | 3 | first chess version |
| Saturday | 9 | all chess logics + tests |
| TBD | TBD | Add exact hours for week two work |
