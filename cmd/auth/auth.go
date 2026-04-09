package auth

import (
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the auth command group
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "认证管理",
		Long:  "管理开放平台认证凭证",
	}

	cmd.AddCommand(newCmdLogin(f))
	cmd.AddCommand(newCmdLogout(f))

	return cmd
}
