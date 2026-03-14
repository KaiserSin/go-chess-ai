# Specification Document

## Programming language and peer review
I am writing the project in the Go language. Also I can review other projects in Python, Java, Kotlin, C++ and JS.

## Algorithms and data structures
The main algorithm of the project is Minimax with alpha beta pruning. I will use the following data structures
* Bitboards. This means using the uint64 data type for very fast storage of the chess board and pieces state.
* Game tree for checking possible game states.
* Built in Go arrays for storing lists of available moves.

## What problem I solve
I am creating an artificial intelligence for a chess game. The bot must be able to play against a human calculate the situation several steps ahead and make strong moves. All game logic is written completely from scratch without using ready made chess libraries.

## Program inputs
The program has a graphical interface. The user clicks on the board and places a piece on the desired square. Then the program checks this move according to the rules of the game and applies it to the virtual board. After that the algorithm calculates its response and the bot also makes its move graphically on the screen.

## Expected time and space complexity
The time complexity of the basic Minimax algorithm is O(b^d) where b is the average number of move options in one situation and d is the search depth. In chess b is usually 35. Alpha beta pruning in the best case reduces this complexity to O(b^(d/2)) because the algorithm ignores bad branches and does not waste time checking them.

The space complexity is O(d) where d is the maximum search depth. Memory is spent mostly on the recursion call stack because the algorithm only stores the current board path and not the entire possible game tree.
## List of sources
* www.chessprogramming.org
* en.wikipedia.org/wiki/Minimax
* en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning

## Core of the project
The core of my project is the implementation of artificial intelligence using the Minimax algorithm with alpha beta pruning. Creating correct moves and a graphical interface are necessary to run the game but the main task and focus of all development is exactly the fast search for the best move and the evaluation of the board position.

## Administrative information
My study program is TKT. The project documentation language is English.