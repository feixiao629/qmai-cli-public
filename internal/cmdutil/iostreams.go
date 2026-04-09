package cmdutil

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/term"
)

// IOStreams provides the standard streams for command I/O
type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer
}

// DefaultIOStreams returns IOStreams connected to os stdin/stdout/stderr
func DefaultIOStreams() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// IsTerminal returns true if stdout is a terminal
func (s *IOStreams) IsTerminal() bool {
	if f, ok := s.Out.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

// TestIOStreams returns IOStreams suitable for testing
func TestIOStreams() (*IOStreams, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	return &IOStreams{
		In:     io.NopCloser(&bytes.Buffer{}),
		Out:    stdout,
		ErrOut: stderr,
	}, stdout, stderr
}
