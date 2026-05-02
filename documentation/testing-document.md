# Testing Document

The project test suite verifies correctness of the chess rules, the gameplay flow, and the AI behavior.
The tests focus on functional logic rather than the desktop UI.

## Testing Strategy

The automated tests concentrate on three areas

- chess rule correctness in the domain layer
- gameplay flow in the application layer
- AI correctness through representative search scenarios

UI and presentation tests are intentionally excluded from the main suite.
The core value of this project is in the chess logic and the search algorithm, so the tests focus on those parts directly.

## Test Layout

The folder layout stays consistent with the layered project structure

- black-box tests are kept under `internal/tests/...`
- white-box tests stay next to the code only when package-private access is required

This keeps behavior-oriented tests separate from implementation-specific tests without moving tests between layers unnecessarily.

The separate `internal/tests` folder is used for tests that look at the project from the outside.
These tests call the same public functions that the rest of the program uses.
Because of this they show that the real behavior works, not only private helper functions.
Some AI tests stay next to the AI code because they need direct access to internal search and evaluation logic.
This keeps normal behavior tests clean while still making it possible to test difficult AI details directly.

## Main Commands

The main automated verification commands are

```text
go test ./...
go test ./internal/...
make test-ai-extended
make coverage
./scripts/coverage.sh
```

## Coverage

The coverage command reports coverage for the current core packages

- `internal/domain/chess/...`
- `internal/application/gameplay`
- `internal/application/ai`

Current coverage report from `make coverage`

```text
Package coverage
domain     85.4%
gameplay   77.1%
ai         92.2%

Combined coverage
total statements 86.7%
```

The coverage script combines the core package reports and enforces at least `75.0%` combined statement coverage.
The domain coverage includes the public chess facade and its implementation subpackages under `internal/domain/chess/...`.
Presentation code is not included in this target because the project tests focus on chess rules, gameplay state, and AI behavior rather than rendering.

The focus stays on representative correctness scenarios rather than rebuilding low-level coverage-only tests.

## Test Scope

### Domain tests

The domain tests are kept under `internal/tests/domain/chess`.
They verify the public behavior of the chess rules.

The domain suite covers

- legal and illegal move handling
- castling legality and castling state updates
- en passant legality and its move window
- pawn promotion requirements and valid promotion choices
- check, checkmate, and stalemate detection
- threefold repetition
- the fifty-move rule
- insufficient material detection
- public draw and query helpers

These tests validate the rules that the rest of the project depends on.

### Gameplay tests

The gameplay tests are kept under `internal/tests/application/gameplay` and `internal/application/gameplay`.
They verify that the application service uses the chess domain correctly.

The gameplay suite covers

- initial board snapshot contents
- square selection behavior
- legal move application and side-to-move changes
- promotion flow and promotion choice handling
- AI move application through the gameplay service
- error handling when the game is already finished
- error handling when promotion is still pending

These tests check the integration between the chess rules and the service layer that the UI uses.

### AI fast black-box tests

The main AI behavior tests are kept under `internal/tests/application/ai`.
They use the public `ai.BestMoveWithin` entry point with a short test budget and verify the bot from the outside.

The fast black-box AI suite covers

- returning a legal move in the initial position
- returning `HasMove = false` in checkmate and stalemate terminal positions
- returning the only legal move when the side to move is under check
- choosing an immediate checkmate when one is available
- finding a forced mate within the test time budget
- choosing a legal promotion move in a promotion-ready position
- avoiding a poisoned capture in a tactical position

These scenarios check correctness rather than playing strength alone.
A chess bot must recognize terminal states, avoid illegal play, and convert tactical wins when the time-limited search can reach them.

### AI extended suite

The project also has a separate extended AI suite behind the `extended` build tag.
It is run with `make test-ai-extended`.

The extended suite uses a fixed corpus of opening positions, tactical builder positions, and terminal positions.
For each non-terminal position it checks that `BestMoveWithin` returns a legal move and that the move can be applied successfully.
For terminal positions it checks that the AI returns `HasMove = false`.

### AI white-box tests

White-box AI tests are kept under `internal/application/ai`.
They use package-private helpers directly when this makes the tested behavior clearer.

The white-box AI tests cover

- returning a legal fallback move when the search deadline is already expired
- matching results between aspiration search and full-window search at a controlled depth of `3`
- finding a forced mate in a controlled depth-`3` search
- evaluating material advantage directly
- checking that White and Black perspectives produce opposite evaluation signs
- rewarding stronger piece placement on the board
- penalizing weak pawn structure
- rewarding king shield in middlegame positions
- rewarding active king placement in endgames

These tests verify internal search and evaluation behavior without exposing extra public configuration.

## Representative Scenarios

The most important test inputs are chosen to cover correctness properties that matter for a chess AI project.

- Initial position AI move verifies that the AI can search a normal high-branching chess position and still return a legal move.
- Terminal checkmate and stalemate positions verify that the AI recognizes finished games and does not invent moves when none should be played.
- Only legal move while in check verifies that move generation and AI selection respect forced defensive positions.
- Immediate checkmate verifies that the AI chooses a winning tactical move when mate is directly available.
- Forced mate within the test budget verifies that the public time-limited search can find a forced win in a representative tactical position.
- Controlled depth-`3` white-box search verifies that depth-specific search behavior still works independently of the production time budget.
- Poisoned capture position verifies that the AI is not only greedy about material and can avoid a capture that loses tactically.
- Promotion-ready position verifies that the AI can return a legal promotion move and that the promoted piece appears on the expected square.
- Direct evaluation positions verify that the heuristic rewards material, piece placement, pawn structure, king safety, and endgame king activity in controlled positions.
- Castling, en passant, promotion, repetition, fifty-move rule, and insufficient material domain tests verify the non-trivial chess rules that the AI and gameplay service rely on.
- Gameplay promotion and AI service tests verify that the application layer handles user-facing state transitions such as pending promotion, finished games, and AI replies.
- Extended AI corpus verifies legal AI calls on opening, tactical, and terminal positions.

## Rationale

The strongest correctness evidence for this project comes from representative chess situations.
For that reason, the suite emphasizes public rule tests, gameplay integration tests, curated AI positions with forced mates and tactical traps, and a separate deterministic extended AI corpus.

This gives stronger support for the correctness of the project than a suite focused on rendering details, text formatting, or internal helper functions.
