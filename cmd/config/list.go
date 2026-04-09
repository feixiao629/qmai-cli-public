package config

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "查看所有配置",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f)
		},
	}
}

func runList(f *cmdutil.Factory) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	out := f.IOStreams.Out

	fmt.Fprintf(out, "active_profile: %s\n", cfg.ActiveProfile)
	fmt.Fprintf(out, "default_format: %s\n", cfg.Format())
	fmt.Fprintf(out, "debug:          %v\n", cfg.Debug)
	fmt.Fprintf(out, "config_path:    %s\n", cfg.Path())

	if len(cfg.Profiles) > 0 {
		fmt.Fprintln(out, "\nProfiles:")
		for name, p := range cfg.Profiles {
			marker := "  "
			if name == cfg.ActiveProfile {
				marker = "* "
			}
			fmt.Fprintf(out, "  %s%s", marker, name)
			if p.ShopCode != "" {
				fmt.Fprintf(out, " (shop: %s)", p.ShopCode)
			}
			fmt.Fprintln(out)
		}
	}

	return nil
}
