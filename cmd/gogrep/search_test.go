package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessPaths(t *testing.T) {
	tests := []struct {
		cfg           config
		name          string
		path          string
		pattern       string
		expectedOut   string
		expectedCount int
	}{
		{
			name:          "Simple One Letter Pattern",
			path:          "test-data/rockbands.txt",
			pattern:       "J",
			cfg:           config{},
			expectedOut:   "Judas Priest\nBon Jovi\nJunkyard\n",
			expectedCount: 3,
		},
		{
			name:    "Recurse A Directory Tree",
			path:    "test-data/*",
			pattern: "Nirvana",
			cfg: config{
				recurse:   true,
				printPath: true,
			},
			expectedOut:   "test-data/rockbands.txt:Nirvana\ntest-data/test-subdir/BFS1985.txt:Since Bruce Springsteen, Madonna, way before Nirvana\ntest-data/test-subdir/BFS1985.txt:On the radio was Springsteen, Madonna, way before Nirvana\ntest-data/test-subdir/BFS1985.txt:And bring back Springsteen, Madonna, way before Nirvana\ntest-data/test-subdir/BFS1985.txt:Bruce Springsteen, Madonna, way before Nirvana\n",
			expectedCount: 5,
		},
		{
			name:    "A Directory Tree without recurse",
			path:    "test-data/*",
			pattern: "Nirvana",
			cfg: config{
				printPath: true,
			},
			expectedOut:   "test-data/rockbands.txt:Nirvana\ngogrep: test-data/test-subdir: Is a directory\n",
			expectedCount: 1,
		},
		{
			name:    "Digit Pattern",
			path:    "test-data/test-subdir/BFS1985.txt",
			pattern: "\\d",
			cfg:     config{},
			expectedOut: `Her dreams went out the door when she turned 24
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
`,
			expectedCount: 9,
		},
		{
			name:          "Word Pattern",
			path:          "test-data/symbols.txt",
			pattern:       "\\w",
			cfg:           config{},
			expectedOut:   "pound\ndollar\n",
			expectedCount: 2,
		},
		{
			name:          "Beginning Of The Line Pattern",
			path:          "test-data/rockbands.txt",
			pattern:       "^A",
			cfg:           config{},
			expectedOut:   "AC/DC\nAerosmith\nAccept\nApril Wine\nAutograph\n",
			expectedCount: 5,
		},
		{
			name:          "End Of The Line Pattern",
			path:          "test-data/rockbands.txt",
			pattern:       "na$",
			cfg:           config{},
			expectedOut:   "Nirvana\n",
			expectedCount: 1,
		},
		{
			name:    "Case Insensitive Pattern",
			path:    "test-data/rockbands.txt",
			pattern: "b",
			cfg: config{
				caseInsensitive: true,
			},
			expectedOut:   "Black Sabbath\nMr. Big\nBon Jovi\nBad English\nBoston\nBad Company\nRainbow\nVandenberg\nLoverboy\nBaton Rogue\nBulletBoys\nLynch Mob\nBlue Murder\nKing Cobra\nBabylon A.D.\nRoxy Blue\nBang Tango\n",
			expectedCount: 17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApplication()
			app.config = tt.cfg

			r, w, err := os.Pipe()
			require.NoError(t, err, "Error creating pipe")

			origStdout := os.Stdout
			os.Stdout = w

			paths, err := filepath.Glob(tt.path)
			require.NoError(t, err, "Error expanding path")

			if app.config.caseInsensitive {
				tt.pattern = "(?i)" + tt.pattern
			}
			pat := regexp.MustCompile(tt.pattern)
			err = app.ProcessPaths(paths, pat)
			require.NoError(t, err)

			w.Close()
			os.Stdout = origStdout

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			require.NoErrorf(t, err, "Failed to read captured output")
			r.Close()
			actual := buf.String()

			assert.Equal(t, tt.expectedOut, actual, "Actual does not match expected")

			actualCount := app.config.matchCount
			assert.Equal(t, tt.expectedCount, actualCount, "Actual count does not match expected")
		})
	}
}

func TestProcessStdin(t *testing.T) {
	tests := []struct {
		name          string
		inFile        string
		cfg           config
		pattern       string
		expectedOut   string
		expectedCount int
	}{
		{
			name:          "Word Pattern",
			inFile:        "test-data/rockbands.txt",
			cfg:           config{},
			pattern:       "Jovi",
			expectedOut:   "Bon Jovi\n",
			expectedCount: 1,
		},
		{
			name:   "Invert Word Pattern",
			inFile: "test-data/symbols.txt",
			cfg: config{
				invertMatch: true,
			},
			pattern:       "\\w",
			expectedOut:   "!\n@\n£\n$\n%\n^\n&\n*\n(\n)\n",
			expectedCount: 10,
		},
		{
			name:    "Digit Pattern",
			inFile:  "test-data/test-subdir/BFS1985.txt",
			cfg:     config{},
			pattern: "\\d",
			expectedOut: `Her dreams went out the door when she turned 24
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 1985
There was U2 and Blondie, and music still on MTV
'Cause she's still preoccupied with 19, 19, 1985
`,
			expectedCount: 9,
		},
		{
			name:          "Beginning Of The Line Pattern",
			inFile:        "test-data/rockbands.txt",
			cfg:           config{},
			pattern:       "^Ap",
			expectedOut:   "April Wine\n",
			expectedCount: 1,
		},
		{
			name:          "End Of The Line Pattern",
			inFile:        "test-data/rockbands.txt",
			cfg:           config{},
			pattern:       "ph$",
			expectedOut:   "Triumph\nAutograph\n",
			expectedCount: 2,
		},
		{
			name:   "Case Insensitive Pattern",
			inFile: "test-data/rockbands.txt",
			cfg: config{
				caseInsensitive: true,
			},
			pattern:       "b",
			expectedOut:   "Black Sabbath\nMr. Big\nBon Jovi\nBad English\nBoston\nBad Company\nRainbow\nVandenberg\nLoverboy\nBaton Rogue\nBulletBoys\nLynch Mob\nBlue Murder\nKing Cobra\nBabylon A.D.\nRoxy Blue\nBang Tango\n",
			expectedCount: 17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdinR, stdinW, err := os.Pipe()
			require.NoError(t, err, "Error creating stdin pipe")

			f, err := os.ReadFile(tt.inFile)
			require.NoError(t, err, "Error reading input file")

			_, err = stdinW.Write(f)
			require.NoError(t, err, "Error writing to pipe")
			stdinW.Close()

			origStdin := os.Stdin
			defer func() { os.Stdin = origStdin }()
			os.Stdin = stdinR

			stdoutR, stdoutW, err := os.Pipe()
			require.NoError(t, err, "Error creating stdout pipe")

			origStdout := os.Stdout
			defer func() { os.Stdout = origStdout }()
			os.Stdout = stdoutW

			app := NewApplication()
			app.config = tt.cfg

			if app.config.caseInsensitive {
				tt.pattern = "(?i)" + tt.pattern
			}
			pat := regexp.MustCompile(tt.pattern)

			err = app.ProcessStdin(pat)
			require.NoError(t, err)

			stdoutW.Close()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, stdoutR)
			require.NoErrorf(t, err, "Failed to read captured output")
			stdoutR.Close()

			actualOut := buf.String()
			assert.Equal(t, tt.expectedOut, actualOut, "Actual does not match expected")

			actualCount := app.config.matchCount
			assert.Equal(t, tt.expectedCount, actualCount, "Actual count does not match expected")
		})
	}
}
