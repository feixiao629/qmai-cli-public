package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var force bool
	var saleChannel int

	cmd := &cobra.Command{
		Use:   "delete <tradeMark>",
		Short: "下架商品",
		Long:  "通过 BatchDown API 下架指定商品（开放平台无直接删除接口）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tradeMark := args[0]

			if !force {
				fmt.Fprintf(f.IOStreams.Out, "确定要下架商品 %s 吗？使用 --force 确认\n", tradeMark)
				return nil
			}

			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}

			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}

			api := client.NewProductAPI(apiClient)
			if err := api.BatchDown(cmd.Context(), profile.ShopCode, []string{tradeMark}, saleChannel); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 商品已下架 (%s)\n", tradeMark)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "确认下架")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "销售渠道")

	return cmd
}
