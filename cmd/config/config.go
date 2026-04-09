package config

import (
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdConfig creates the config command group
func NewCmdConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置管理",
		Long:  "管理 qmai CLI 配置，包括初始化、读写配置项、多门店 profile",
	}

	cmd.AddCommand(newCmdInit(f))
	cmd.AddCommand(newCmdSet(f))
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdProfile(f))

	return cmd
}
