package config

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/core"
	"github.com/spf13/cobra"
)

func newCmdInit(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "交互式初始化配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(f)
		},
	}
}

func runInit(f *cmdutil.Factory) error {
	out := f.IOStreams.Out
	reader := bufio.NewReader(f.IOStreams.In)

	fmt.Fprintln(out, "qmai CLI 配置初始化")
	fmt.Fprintln(out, strings.Repeat("-", 30))

	// Profile name
	fmt.Fprint(out, "Profile 名称 (default): ")
	profileName, _ := reader.ReadString('\n')
	profileName = strings.TrimSpace(profileName)
	if profileName == "" {
		profileName = "default"
	}

	// Shop Code (optional)
	fmt.Fprint(out, "门店编码 shopCode (可选): ")
	shopCode, _ := reader.ReadString('\n')
	shopCode = strings.TrimSpace(shopCode)

	// API Base URL
	fmt.Fprintf(out, "API Base URL (%s): ", core.DefaultBaseURL)
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)

	// Default format
	fmt.Fprint(out, "默认输出格式 [table/json/csv] (table): ")
	format, _ := reader.ReadString('\n')
	format = strings.TrimSpace(format)
	if format == "" {
		format = "table"
	}

	// Build config
	cfg := &core.CliConfig{
		ActiveProfile: profileName,
		DefaultFormat: format,
		Profiles:      make(map[string]*core.Profile),
	}

	profile := &core.Profile{
		Name:     profileName,
		ShopCode: shopCode,
	}
	if baseURL != "" {
		profile.BaseURL = baseURL
	}
	cfg.Profiles[profileName] = profile

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	path, _ := core.ConfigPath()
	fmt.Fprintf(out, "\n✓ 配置已保存到 %s\n", path)
	fmt.Fprintf(out, "  Active profile: %s\n", profileName)
	fmt.Fprintln(out, "\n下一步: 运行 'qmai auth login' 配置开放平台凭证")

	return nil
}
