# goko

> Automation cli tool for jutge.org problems of PRO1

![Go Version](https://img.shields.io/badge/go-1.26-blue)
![License](https://img.shields.io/badge/license-BSD--3--Clause-blue)

## About

A cli tool built to save some time doing repetitive setups and tests of the Jutge problems. Written in idiomatic(not yet. The only idiomatic parts are the ones I've stolen from here and there online) golang because C++ sucks. Works on Linux and probably MacOS(not sure, never tried, but in theory the file paths should be the same, but you would always need to pass an "app to open with" argument because the "xdg-open" fallback doesn't work on MacOS), not on windows. 

## Features

- Downloads and unpacks input examples and test cases when provided a public problem ID. 
- Builds the written C++ program using clang(you may change it for p1++ on the source level if needed. p1++ must be a pseudo binary(a file with a shell script in your $PATH), not an alias for it)
- Runs through the public test cases, so you don't have to send the thing to Jutge over and over again. 

## Installation

```bash
go install github.com/Elmer-V/goko@latest
```
Needs clang(clang++) and diff to be on the machine in runtime to function. Obviously needs go to be installed too for compilation. 

Or build from source:

```bash
git clone https://github.com/Elmer-V/goko.git
cd goko
go build
```

## Usage

```bash
goko <mode> <id>
```

| Mode      | Alias | Description                |
|-----------|-------|----------------------------|
| `set`     | `s`   | Needs a problem ID and then can take an app name to open the newly created empty cpp file with    |
| `check`   | `c`   | Takes a problem ID if run from a directory above the problem or, if run directly from the problem directory, works without any arguments |

### Example

```bash
goko set P42280
goko s P42280 codium
goko c P42280
goko c
```

## Project structure

```
goko/
├── main.go         # Entry point
├── set.go          # set mode
├── check.go        # check mode
├── downloader.go   # Fetches problem archives
├── dirwork.go      # Directory / unzip helpers
├── go.mod
├── .gitignore
├── LICENSE
└── README.md
```

## Contribating

1. While you may try to commit something, I'm going to be too tired and lazy to review/merge, so you'd better clone the thing if you need it.

## License

Distributed under the BSD 3-Clause License. See `LICENSE` for more information.
