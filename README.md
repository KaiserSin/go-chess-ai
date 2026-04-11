# Go Chess AI

This project is a chess application written in Go. It includes the chess rules, a desktop board built with Ebiten, and the base structure for a future AI player.

## Install Go

Official Go installation guide

- https://go.dev/doc/install

### macOS

1. Download the macOS installer package from `https://go.dev/dl/`
2. Open the `.pkg` file and follow the installer steps
3. Reopen Terminal
4. Check that Go is installed

```bash
go version
```

### Linux

1. Download the Linux archive from `https://go.dev/dl/`
2. Remove any old Go folder from `/usr/local/go`
3. Extract the new archive into `/usr/local`
4. Add `/usr/local/go/bin` to your `PATH`
5. Open a new terminal or reload your shell profile
6. Check that Go is installed

```bash
go version
```

Example `PATH` line

```bash
export PATH=$PATH:/usr/local/go/bin
```

### Windows

1. Download the Windows `.msi` installer from `https://go.dev/dl/`
2. Open the installer and follow the steps
3. Reopen Command Prompt or PowerShell
4. Check that Go is installed

```bash
go version
```

## Run the project

After Go is installed, start the application with

```bash
make run
```

To see the other available commands, run

```bash
make help
```
