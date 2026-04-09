package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdBatchStatus(f *cmdutil.Factory) *cobra.Command {
	var action string
	var tradeMarks []string
	var saleChannel int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "batch-status",
		Short: "批量上架/下架",
		Long: `批量修改商品上架或下架状态。

示例:
  qmai product batch-status --action up --trade-marks TM001,TM002
  qmai product batch-status --action down --trade-marks TM001 --sale-channel 1
  qmai product batch-status --action up --trade-marks TM001 --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if action != "up" && action != "down" {
				return fmt.Errorf("--action 必须为 'up' 或 'down'")
			}
			if len(tradeMarks) == 0 {
				return fmt.Errorf("--trade-marks 必填")
			}

			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将要%s %d 个商品\n",
					map[string]string{"up": "上架", "down": "下架"}[action],
					len(tradeMarks))
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

			switch action {
			case "up":
				if err := api.BatchUp(cmd.Context(), profile.ShopCode, tradeMarks, saleChannel); err != nil {
					return err
				}
				fmt.Fprintf(f.IOStreams.Out, "✓ 已上架 %d 个商品\n", len(tradeMarks))
			case "down":
				if err := api.BatchDown(cmd.Context(), profile.ShopCode, tradeMarks, saleChannel); err != nil {
					return err
				}
				fmt.Fprintf(f.IOStreams.Out, "✓ 已下架 %d 个商品\n", len(tradeMarks))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&action, "action", "", "操作: up (上架) 或 down (下架)")
	cmd.Flags().StringSliceVar(&tradeMarks, "trade-marks", nil, "商品编码列表")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "销售渠道")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览")

	return cmd
}
