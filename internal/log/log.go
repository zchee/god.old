// Copyright 2017 The log Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var debug = len(os.Getenv("GOD_DEBUG")) > 0

// Use golang's standard logger by default.
//
// Access is not mutex-protected, do not modify except in init()
// functions.
var logger = log.New(os.Stderr, "", log.Lshortfile)

func init() {
	// enable spew ContinueOnMethod
	spew.Config.ContinueOnMethod = true
}

// New creates a new std log packages Logger.
// The out variable sets the destination to which log data will be written.
//
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *log.Logger {
	return log.New(out, prefix, flag)
}

// Fatal is equivalent to Print() followed by a call to os.Exit() with a non-zero exit code.
func Fatal(a ...interface{}) {
	logger.Output(2, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit() with a non-zero exit code.
func Fatalf(format string, a ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit()) with a non-zero exit code.
func Fatalln(a ...interface{}) {
	logger.Output(2, fmt.Sprintln(a...))
	os.Exit(1)
}

// Print prints to the logger. Arguments are handled in the manner of fmt.Print.
func Print(a ...interface{}) {
	logger.Output(2, fmt.Sprint(a...))
}

// Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.
func Printf(format string, a ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, a...))
}

// Println prints to the logger. Arguments are handled in the manner of fmt.Println.
func Println(a ...interface{}) {
	logger.Output(2, fmt.Sprintln(a...))
}

// Debug prints to the logger if debug is true. Arguments are handled in the manner of fmt.Print.
func Debug(a ...interface{}) {
	if debug {
		logger.Output(2, fmt.Sprint(a...))
	}
}

// Debugf prints to the logger if debug is true. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, a ...interface{}) {
	if debug {
		logger.Output(2, fmt.Sprintf(format, a...))
	}
}

// Debugln prints to the logger if debug is true. Arguments are handled in the manner of fmt.Println.
func Debugln(a ...interface{}) {
	if debug {
		logger.Output(2, fmt.Sprintln(a...))
	}
}

// Dump spew package wrapper that displays the passed parameters to standard out with newlines,
// customizable indentation, and additional debug information such as complete types and
// all pointer addresses used to indirect to the final value.
func Dump(a ...interface{}) {
	logger.Output(2, spew.Sdump(a...))
}
