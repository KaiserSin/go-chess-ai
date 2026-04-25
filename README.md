# Go Chess AI

Go Chess AI is a desktop chess application written in Go.
The project includes a full chess rule implementation, a graphical board built with Ebiten, and a human versus AI game mode.

The chess logic is implemented without external chess libraries.
The AI uses a fixed search depth of `3` and is based on alpha-beta search with iterative deepening, aspiration windows, quiescence search, move ordering, and a hand-tuned positional evaluation function.

## Features

- playable desktop chess board with mouse input
- side selection before the game starts
- legal move validation for all pieces
- check, checkmate, and stalemate detection
- castling, en passant, pawn promotion, and common draw rules
- AI opponent with fixed-depth search and positional evaluation
- automated tests for chess rules, gameplay flow, and AI behavior

## Documentation

- [Specification document](documentation/specification.md)
- [User guide](documentation/user-guide.md)
- [Implementation document](documentation/implementation-document.md)
- [Testing document](documentation/testing-document.md)
- [Architecture overview](documentation/architecture.md)
- Weekly reports
  - [Week 1](documentation/weekly-reports/week-one.md)
  - [Week 2](documentation/weekly-reports/week-two.md)
  - [Week 3](documentation/weekly-reports/week-three.md)
  - [Week 4](documentation/weekly-reports/week-four.md)
  - [Week 5](documentation/weekly-reports/week-five.md)

## Quick Start

Install Go `1.26`, clone the repository, and run the desktop application from the repository root

```bash
make run
```

For full installation, gameplay, input, and command instructions, see the [User guide](documentation/user-guide.md).

Run the main automated test suite with

```bash
make test
```

Other development commands are listed in the user guide and in `make help`.

## Project Structure

```text
cmd/chess-desktop/              desktop entry point
internal/domain/chess/          chess rules and board state
internal/application/ai/        search and evaluation logic
internal/application/gameplay/  gameplay service used by the UI
internal/presentation/ebiten/   desktop rendering and input handling
internal/tests/                 black-box domain and application tests
documentation/                  course documentation and weekly reports
```
