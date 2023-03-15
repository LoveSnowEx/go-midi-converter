# go-midi-converter

This is a simple tool for generate midi file from numbered musical notation.

## Installation

```sh
go mod tidy
```

## Run

You can write numbered musical notation at `in.txt`,

For example:

```text
1155665-4433221-5544332-
5544332-1155665-4433221-
```

and run go program with command:

```sh
go run .
```

then the output wil be stored at `out.mid`
