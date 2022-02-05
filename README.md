# Nokia Defence

A Tower Defence game for Web and Desktop written as if for the Nokia 3310:

- 84x48 screen resolution
- only 2 colours light/dark green
- monophonic music only
- max 12 buttons allowed

This is a submission for [Nokia 3310 Jam 4](https://itch.io/jam/nokiajam4) made using the [Ebiten](https://ebiten.org/) library.

## For game testers

<!-- TODO: add a link to the latest downloads page -->

Game controls:
- F: toggle full-screen
- Q: quit the game
- Space: move up

## For programmers

Make sure you have [Go 1.17 or later](https://go.dev/) to contribute to the game

To build the game yourself, run: `go build .` it will produce an nokia-defence file and on Windows nokia-defence.exe.

To run the tests, run: `go test ./...` but there are no tests yet.

The project has a very simple, flat structure, the first place to start looking is the main.go file.
