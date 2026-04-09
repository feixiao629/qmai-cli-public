package config

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/core"
	"github.com/spf13/cobra"
)

func newCmdProfile(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "多门店 profile 管理",
	}

	cmd.AddCommand(newCmdProfileAdd(f))
	cmd.AddCommand(newCmdProfileRemove(f))
	cmd.AddCommand(newCmdProfileList(f))

	return cmd
}

func newCmdProfileAdd(f *cmdutil.Factory) *cobra.Command {
	var shopCode, baseURL string

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "添加新 profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfileAdd(f, args[0], shopCode, baseURL)
		},
	}

	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "API Base URL")

	return cmd
}

func runProfileAdd(f *cmdutil.Factory, name, shopCode, baseURL string) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[name]; exists {
		return fmt.Errorf("profile %q already exists", name)
	}

	p := &core.Profile{
		Name:     name,
		ShopCode: shopCode,
		BaseURL:  baseURL,
	}
	cfg.SetProfile(name, p)

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Fprintf(f.IOStreams.Out, "✓ Profile %q added\n", name)
	return nil
}

func newCmdProfileRemove(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <name>",
		Short:   "删除 profile",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfileRemove(f, args[0])
		},
	}
}

func runProfileRemove(f *cmdutil.Factory, name string) error {
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[name]; !exists {
		return fmt.Errorf("%w: %s", core.ErrProfileNotFound, name)
	}

	delete(cfg.Profiles, name)
	if cfg.ActiveProfile == name {
		cfg.ActiveProfile = ""
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Fprintf(f.IOStreams.Out, "✓ Profile %q removed\n", name)
	return nil
}

func newCmdProfileList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有 profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			out := f.IOStreams.Out
			if len(cfg.Profiles) == 0 {
				fmt.Fprintln(out, "No profiles configured. Run 'qmai config init' to create one.")
				return nil
			}

			for name, p := range cfg.Profiles {
				marker := "  "
				if name == cfg.ActiveProfile {
					marker = "* "
				}
				fmt.Fprintf(out, "%s%s", marker, name)
				if p.ShopCode != "" {
					fmt.Fprintf(out, " (shop: %s)", p.ShopCode)
				}
				if p.BaseURL != "" {
					fmt.Fprintf(out, " [%s]", p.BaseURL)
				}
				fmt.Fprintln(out)
			}
			return nil
		},
	}
}
