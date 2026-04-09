package auth

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/keychain"
	"github.com/spf13/cobra"
)

func newCmdLogout(f *cmdutil.Factory) *cobra.Command {
	var profileName string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "清除凭证",
		Long:  "清除当前 profile 的开放平台凭证",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogout(f, profileName)
		},
	}

	cmd.Flags().StringVar(&profileName, "profile", "", "Profile to logout from (default: active profile)")

	return cmd
}

func runLogout(f *cmdutil.Factory, profileName string) error {
	if profileName == "" {
		profileName = f.EffectiveProfile()
	}

	// Delete openKey from keychain
	kc := f.Keychain
	if kc == nil {
		kc = keychain.NewOSKeychain()
	}
	_ = kc.Delete(profileName) // ignore error if not found

	// Clear profile credentials from config
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	if p, ok := cfg.Profiles[profileName]; ok {
		p.OpenId = ""
		p.GrantCode = ""
		p.ShopCode = ""
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("保存配置失败: %w", err)
		}
	}

	fmt.Fprintf(f.IOStreams.Out, "✓ 已登出 (profile: %s)\n", profileName)
	return nil
}
