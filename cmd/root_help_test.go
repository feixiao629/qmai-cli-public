package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestRootCommandTreeHelp walks the full command tree and ensures each node's
// Help() succeeds (no panics, no init errors from conflicting flags, etc.).
func TestRootCommandTreeHelp(t *testing.T) {
	root := NewRootCmd("test")
	root.SetOut(ioDiscard{})
	root.SetErr(ioDiscard{})
	root.InitDefaultHelpCmd()

	var visit func(*cobra.Command, []string)
	visit = func(c *cobra.Command, path []string) {
		t.Helper()
		name := c.Name()
		p := append(append([]string(nil), path...), name)

		buf := new(bytes.Buffer)
		c.SetOut(buf)
		c.SetErr(buf)
		if err := c.Help(); err != nil {
			t.Fatalf("Help() failed for %s: %v", strings.Join(p, " "), err)
		}
		if buf.Len() == 0 {
			t.Fatalf("Help() wrote nothing for %s", strings.Join(p, " "))
		}

		for _, sub := range c.Commands() {
			if sub == nil || sub.Hidden {
				continue
			}
			switch sub.Name() {
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
				continue
			}
			visit(sub, p)
		}
	}

	visit(root, nil)
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
