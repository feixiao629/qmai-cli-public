package auth

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/core"
	"github.com/madaima/qmai-cli/internal/keychain"
	"github.com/spf13/cobra"
)

func newCmdLogin(f *cmdutil.Factory) *cobra.Command {
	var profileName string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "配置开放平台凭证",
		Long:  "输入企迈开放平台的 openId、grantCode、openKey 和门店编码",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(f, profileName)
		},
	}

	cmd.Flags().StringVar(&profileName, "profile", "", "Profile name (default: \"default\")")

	return cmd
}

func runLogin(f *cmdutil.Factory, profileName string) error {
	if profileName == "" {
		profileName = f.EffectiveProfile()
		if profileName == "" {
			profileName = "default"
		}
	}

	out := f.IOStreams.Out
	reader := bufio.NewReader(f.IOStreams.In)

	fmt.Fprint(out, "请输入 openId: ")
	openId, _ := reader.ReadString('\n')
	openId = strings.TrimSpace(openId)
	if openId == "" {
		return fmt.Errorf("openId 不能为空")
	}

	fmt.Fprint(out, "请输入 grantCode: ")
	grantCode, _ := reader.ReadString('\n')
	grantCode = strings.TrimSpace(grantCode)
	if grantCode == "" {
		return fmt.Errorf("grantCode 不能为空")
	}

	fmt.Fprint(out, "请输入 openKey: ")
	openKey, _ := reader.ReadString('\n')
	openKey = strings.TrimSpace(openKey)
	if openKey == "" {
		return fmt.Errorf("openKey 不能为空")
	}

	fmt.Fprint(out, "请输入门店编码 (shopCode): ")
	shopCode, _ := reader.ReadString('\n')
	shopCode = strings.TrimSpace(shopCode)

	// Store openKey in keychain (sensitive)
	kc := f.Keychain
	if kc == nil {
		kc = keychain.NewOSKeychain()
	}
	if err := kc.Set(profileName, openKey); err != nil {
		return fmt.Errorf("保存 openKey 失败: %w", err)
	}

	// Store openId/grantCode/shopCode in config (non-sensitive)
	cfg, err := f.Config()
	if err != nil {
		return err
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]*core.Profile)
	}

	cfg.SetProfile(profileName, &core.Profile{
		OpenId:    openId,
		GrantCode: grantCode,
		ShopCode:  shopCode,
	})
	cfg.ActiveProfile = profileName

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Fprintf(out, "✓ 凭证已保存 (profile: %s)\n", profileName)
	return nil
}
