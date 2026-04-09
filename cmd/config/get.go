package config

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "读取配置项",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, args[0])
		},
	}
}

func runGet(f *cmdutil.Factory, key string) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	var value string
	switch key {
	case "default_format":
		value = cfg.Format()
	case "debug":
		value = fmt.Sprintf("%v", cfg.Debug)
	case "active_profile":
		value = cfg.ActiveProfile
	case "base_url":
		value = cfg.BaseURL()
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	fmt.Fprintln(f.IOStreams.Out, value)
	return nil
}
