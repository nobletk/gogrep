# GoGrep

A CLI tool to search any given input files, selecting lines that
match one or more patterns.

## Table of Contents

- [Features](#features)
- [Usage](#usage)
    - [Examples](#examples)
- [Getting started](#getting-started)
    - [Clone the repo](#clone-the-repo)

## Features

* Empty expression, complex pattern search.

* Recursive search in a directory tree.

* Invert search excluding pattern from matching lines.

* Digit "\d" and Word "\w" pattern matching.

* Beginning of the line "^" and end of the line "$" pattern matching.

* Case insensitive search.

* Support Stdin input pipe.

## Usage

```
> gogrep [OPTION...] PATTERN [FILE...]
```

### Examples

* Empty expression

```
// Test the difference between the original
// and the gogrep result
~> gogrep "" test.txt | diff test.txt -
```

* Pattern search

```
~> gogrep Jovi rockbands.txt
Bon Jovi
```

* `-r` or `--recursive`: Recursive search in a directory tree

```
~> gogrep -r Nirvana test-data/*
test-data/rockbands.txt:Nirvana
test-data/test-subdir/BFS1985.txt:Since Bruce Springsteen, Madonna, way before Nirvana
test-data/test-subdir/BFS1985.txt:On the radio was Springsteen, Madonna, way before Nirvana
test-data/test-subdir/BFS1985.txt:And bring back Springsteen, Madonna, way before Nirvana
test-data/test-subdir/BFS1985.txt:Bruce Springsteen, Madonna, way before Nirvana
```

* `-v` or `--invert-match`: Inverted search excluding pattern 

```
~> gogrep -v Nirvana test-data/* | gogrep -v Madonna
rockbands.txt:Nirvana
```

* Digit "\d" pattern search

```
~> gogrep "\d" test-subdir/BFS1985.txt
Her dreams went out the door when she turned 24
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
```

* Word "\w" pattern search

```
~> gogrep "\w" symbols.txt
pound
dollar
```

* Beginning of the line "^" pattern matching

```
~> gogrep ^A rockbands.txt
AC/DC
Aerosmith
Accept
April Wine
Autograph
```

* End of the line "$" pattern matching

```
~> gogrep na$ rockbands.txt
Nirvana
```

* Case insensitive search

```
~> gogrep -i A rockbands.txt | wc -l
58
```


* Support Stdin input pipe.

```
~> gogrep J test-data/rockbands.txt | gogrep -v Jovi
Judas Priest
Junkyard
```

## Getting started

### Clone the repo

```
~> git clone https://github.com/nobletk/gogrep
# then build the binary
~> make build
```
