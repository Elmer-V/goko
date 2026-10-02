# gompare

> One-line description of what this project does. (placeholder)

![Go Version](https://img.shields.io/badge/go-1.26-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## About

Short paragraph describing the project. Replace this with a real summary of
what `gompare` is and why it exists.

## Features

- Feature one
- Feature two
- Feature three

## Installation

```bash
go install github.com/Elmer-V/gompare@latest
```

Or build from source:

```bash
git clone https://github.com/Elmer-V/gompare.git
cd gompare
go build
```

## Usage

```bash
gompare <mode> <id>
```

| Mode      | Alias | Description                |
|-----------|-------|----------------------------|
| `set`     | `s`   | Description of set mode    |
| `check`   | `c`   | Description of check mode  |

### Example

```bash
gompare set P42280
```

## Project structure

```
gompare/
├── main.go         # Entry point
├── set.go          # set mode
├── check.go        # check mode
├── downloader.go   # Fetches problem archives
├── dirwork.go      # Directory / unzip helpers
└── go.mod
```

## Contribating

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.
