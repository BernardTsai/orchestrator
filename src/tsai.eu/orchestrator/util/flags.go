package util

import (
	"flag"
)

var debug *bool

//------------------------------------------------------------------------------

// ParseCommandLineOptions parses the options of the CLI
func ParseCommandLineOptions() {
	debug = flag.Bool("debug", false, "turns on debug logging")

	flag.Parse()
}

//------------------------------------------------------------------------------

// Debug indicates if debug mode has been requested
func Debug() bool {
	return *debug
}

//------------------------------------------------------------------------------
