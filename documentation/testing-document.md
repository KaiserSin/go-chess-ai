# Testing Document

The project test suite verifies correctness of the chess rules, the gameplay flow, and the AI behavior.
The tests focus on functional logic rather than the desktop UI.

## Testing Strategy

The automated tests concentrate on three areas:

- chess rule correctness in the domain layer
- gameplay flow in the application layer
- AI correctness through representative search scenarios

UI and presentation tests are intentionally excluded from the main suite.
The core value of this project is in the chess logic and the search algorithm, so the tests focus on those parts directly.

## Test Layout

The folder layout stays consistent with the project architecture:

- black-box tests are kept under `internal/tests/...`
- white-box tests stay next to the code only when package-private access is required

This keeps behavior-oriented tests separate from implementation-specific tests without moving tests between layers unnecessarily.

## Main Commands

The main automated verification commands are

```text
go test ./...
go test ./internal/...
make test-ai-extended
./scripts/coverage.sh
```

## Coverage

The coverage command reports coverage for the current core packages:

- `internal/domain/chess`
- `internal/application/gameplay`
- `internal/application/ai`

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
They use the public `ai.BestMove` entry point and verify the bot from the outside.

The fast black-box AI suite covers

- returning a legal move in the initial position
- returning `HasMove = false` in checkmate and stalemate terminal positions
- returning the only legal move when the side to move is under check
- choosing an immediate checkmate when one is available
- finding a forced mate inside the fixed production depth of `3` plies
- choosing a legal promotion move in a promotion-ready position
- avoiding a poisoned capture in a tactical position

These scenarios check correctness rather than playing strength alone.
A chess bot must recognize terminal states, avoid illegal play, and convert a forced mate when the mate is inside its search depth.

### AI extended suite

The project also has a separate deterministic extended AI suite behind the `extended` build tag.
It is run with `make test-ai-extended`.

The extended suite uses a fixed corpus of opening positions, tactical builder positions, and terminal positions.
For each non-terminal position it checks that repeated calls to `BestMove` return the same result, that the move is legal, and that the move can be applied successfully.
For terminal positions it checks that the AI returns `HasMove = false`.

### AI white-box test

One internal AI test is kept under `internal/application/ai/search_internal_test.go`.
It uses the package-private `bestMove(position, depth)` helper directly.

The white-box AI test covers

- a curated forced-mate position that requires depth `5`, which is deeper than the fixed runtime depth of `3`

This test verifies deeper search behavior without changing the public API or the production search depth.

## Rationale

The strongest correctness evidence for this project comes from representative chess situations.
For that reason, the suite emphasizes public rule tests, gameplay integration tests, curated AI positions with forced mates and tactical traps, and a separate deterministic extended AI corpus.

This gives stronger support for the correctness of the project than a suite focused on rendering details, text formatting, or internal helper functions.
