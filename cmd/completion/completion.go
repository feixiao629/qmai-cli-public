package completion

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCmdCompletion creates the completion command
func NewCmdCompletion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "生成 shell 自动补全脚本",
		Long: `Generate shell completion scripts for qmai.

Examples:
  # Bash
  source <(qmai completion bash)

  # Zsh
  qmai completion zsh > "${fpath[1]}/_qmai"

  # Fish
  qmai completion fish | source

  # PowerShell
  qmai completion powershell | Out-String | Invoke-Expression`,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
