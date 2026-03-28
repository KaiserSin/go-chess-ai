# Weekly Report 3

During the third week the project became much easier to use and observe. The chess logic from the earlier work is now connected to a working desktop interface. Because of this the program can already be started, the board can be seen on screen, and the game can be played locally by two sides on one machine.

The main work this week was in the `application` and `presentation` layers. I added the gameplay service that connects the chess rules to the user interface. It now supports selecting a piece, showing legal target squares, making moves, and handling pawn promotion. I also finished the basic Ebiten board so the player can click pieces and move them on a real board instead of testing only through code.

The user interface is still simple, but it already works well enough for this stage. The board is drawn on screen, the pieces are visible, and the current game state can be observed directly. This means the core functionality of the project can now be run and checked in practice, which was one of the main goals for this week.

Testing and code quality also moved forward. The project now has unit tests not only for the chess domain, but also for the gameplay layer and small helper parts of the presentation layer. I started the testing document and added the current coverage report there. I also added lint support with `make lint`, so code quality can be checked more easily while developing.

Next week I want to prepare the project for peer review, but I also want to finally start the algorithm part of the project. The main goal of this work is still the chess AI, so the next important step is to begin the first version of the Minimax search with alpha beta pruning on top of the current playable game. I also want to improve code clarity, keep the tests in good shape, continue the implementation and testing documentation, and add missing useful tests if needed.

## Time tracking

| Day | Hours | Description |
| - | - | - |
| Saturday | 10 | UI, gameplay, testing, documentation, and code quality improvements |
| Total | 10 | |
