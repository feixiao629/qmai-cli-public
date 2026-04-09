package config

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdSet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "设置配置项",
		Long: `Set a configuration value.

Available keys:
  default_format    Output format (json, table, csv)
  debug             Debug mode (true, false)
  active_profile    Active profile name`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSet(f, args[0], args[1])
		},
	}
}

func runSet(f *cmdutil.Factory, key, value string) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	switch key {
	case "default_format":
		cfg.DefaultFormat = value
	case "debug":
		cfg.Debug = value == "true"
	case "active_profile":
		if _, ok := cfg.Profiles[value]; !ok {
			return fmt.Errorf("profile %q not found", value)
		}
		cfg.ActiveProfile = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Fprintf(f.IOStreams.Out, "✓ %s = %s\n", key, value)
	return nil
}
