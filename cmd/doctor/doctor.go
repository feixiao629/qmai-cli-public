package doctor

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/core"
	"github.com/madaima/qmai-cli/internal/keychain"
	"github.com/spf13/cobra"
)

// NewCmdDoctor creates the doctor command
func NewCmdDoctor(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "诊断检查",
		Long:  "检查配置、连接和认证状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(f)
		},
	}
}

func runDoctor(f *cmdutil.Factory) error {
	out := f.IOStreams.Out
	allOK := true

	// 1. Config check
	fmt.Fprintln(out, "检查配置...")
	cfg, err := f.Config()
	if err != nil {
		fmt.Fprintf(out, "  ✗ 配置加载失败: %v\n", err)
		allOK = false
	} else {
		fmt.Fprintf(out, "  ✓ 配置文件: %s\n", cfg.Path())
		fmt.Fprintf(out, "  ✓ Active profile: %s\n", cfg.ActiveProfile)
		fmt.Fprintf(out, "  ✓ API Base URL: %s\n", cfg.BaseURL())
	}

	if cfg == nil {
		cfg = &core.CliConfig{}
	}

	// 2. Auth check
	fmt.Fprintln(out, "\n检查认证...")
	profile := f.EffectiveProfile()

	kc := f.Keychain
	if kc == nil {
		kc = keychain.NewOSKeychain()
	}

	openKey, err := kc.Get(profile)
	if err != nil {
		fmt.Fprintf(out, "  ✗ Keychain 访问失败: %v\n", err)
		allOK = false
	} else if openKey == "" {
		fmt.Fprintf(out, "  ✗ openKey 未存储 (profile: %s)\n", profile)
		allOK = false
	} else {
		fmt.Fprintf(out, "  ✓ openKey 已存储 (profile: %s)\n", profile)
	}

	if p := cfg.Profiles[profile]; p != nil {
		if p.OpenId != "" {
			fmt.Fprintf(out, "  ✓ openId: %s\n", p.OpenId)
		} else {
			fmt.Fprintf(out, "  ✗ openId 未配置\n")
			allOK = false
		}
		if p.ShopCode != "" {
			fmt.Fprintf(out, "  ✓ shopCode: %s\n", p.ShopCode)
		} else {
			fmt.Fprintf(out, "  ! shopCode 未配置\n")
		}
	}

	// 3. Connectivity check
	fmt.Fprintln(out, "\n检查连接...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.BaseURL(), nil)
	if err != nil {
		fmt.Fprintf(out, "  ✗ 请求创建失败: %v\n", err)
		allOK = false
	} else {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(out, "  ✗ 连接失败: %v\n", err)
			allOK = false
		} else {
			resp.Body.Close()
			fmt.Fprintf(out, "  ✓ API 可达 (HTTP %d)\n", resp.StatusCode)
		}
	}

	// Summary
	fmt.Fprintln(out)
	if allOK {
		fmt.Fprintln(out, "所有检查通过 ✓")
	} else {
		fmt.Fprintln(out, "存在问题，请检查上述输出")
	}

	return nil
}
