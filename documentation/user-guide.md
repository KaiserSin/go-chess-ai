# User Guide

This guide explains how to install the required tools, run the chess application, play a game, and repeat the main project checks.

## Requirements

The project is written in Go and currently uses Go `1.26`.
Install Go before running the application or the tests.

Official Go installation guide

- go.dev/doc/install

After installation, check that Go is available from the terminal

```bash
go version
```

The command should print the installed Go version.

## Running The Application

Run the application from the repository root directory

```bash
make run
```

This command starts the Ebiten desktop application from `cmd/chess-desktop`.
The program opens a graphical chess window.

## Playing A Game

When the program starts, it shows a side selection menu.
Choose either `Play as White` or `Play as Black`.

After the game starts

1. Click one of your pieces to select it.
2. The board shows legal target squares for the selected piece.
3. Click one of the legal target squares to make the move.
4. Wait for the AI to calculate and make its reply.
5. Continue playing until the game ends.

If you choose to play as Black, the board is shown from Black's perspective and the AI makes the first move as White.

The game ends automatically on checkmate, stalemate, insufficient material, the fifty-move rule, or threefold repetition.

## Pawn Promotion

When one of your pawns reaches the final rank, the game opens a promotion overlay.
Click the piece you want to promote to.

The available promotion choices are

- Queen
- Rook
- Bishop
- Knight

The move is completed after the promotion piece is selected.

## Accepted Inputs

The application does not read chess moves from text files or command-line move input during play.
The game input is graphical

- mouse clicks on the side selection buttons
- mouse clicks on board squares
- mouse clicks on promotion choices

Terminal commands are used only for running the application and project checks.

## Development Commands

Show all available Make targets

```bash
make help
```

Run all Go tests

```bash
make test
```

Run the extended deterministic AI test suite

```bash
make test-ai-extended
```

Run the core coverage check

```bash
make coverage
```

Generate `coverage.out` and `coverage.html`

```bash
make coverage-html
```

Run static checks with `go vet` and `golangci-lint`

```bash
make lint
```
