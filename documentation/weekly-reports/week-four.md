# Weekly Report 4

During the fourth week the project moved from a playable local chess board to a real human versus AI game. The main result of this week is that the program now has its first working chess engine and the user can play against it through the desktop interface.

The largest part of the work was done in the AI layer. I first added a position evaluation function and then started building the actual search on top of it. The search now uses alpha beta pruning, iterative deepening, move ordering, a transposition table, quiescence search, and root level goroutines. Because of this the AI is now able to choose legal moves by searching future positions instead of only reacting through hardcoded rules.

I also improved the evaluation of positions. At first the evaluation only looked at material, but later I added piece square tables and simple positional terms. The current evaluation now includes material, piece placement, pawn structure, king safety, and simple endgame bonuses. This makes the AI less naive than a pure material counter, even though it is still far from a strong chess engine.

Another important part of the week was the user interface. I added a start menu where the player can choose to play as White or Black. I also added AI depth selection and limited the usable depth to a small range so the program stays responsive. When the player chooses Black, the board is rotated so the pieces are shown from the correct side. I also changed the move hint markers so legal target squares are shown as simple dots.

The game flow is now closer to a real chess game against a computer. After the player makes a move, the UI first shows that move on the board and only after that the AI replies. This makes the game easier to follow. The hot seat style local play is no longer the main scenario. The program now clearly starts from the side selection menu and then continues as a human versus AI game.

Documentation also moved forward this week. I added the implementation document and updated it to describe the current program structure, the main algorithms, the current complexity analysis, the limitations of the project, and the use of language models during the work.

Testing remained an important part of the project. I added and updated tests for the AI search, evaluation, gameplay service, and UI helpers. This helped keep the project stable while the search logic and evaluation became more complex. At the end of the week the full automated test suite still passed.

Next week I want to keep improving the playing strength of the AI and continue polishing the project documentation. The most useful next steps would be better evaluation, benchmark style performance measurements, and maybe the first opening book support.

## Time tracking

| Day | Hours | Description |
| - | - | - |
| To be added | - | AI search, evaluation, UI improvements, and documentation |
| Total | - | |
