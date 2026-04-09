package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdEstimateClear(f *cmdutil.Factory) *cobra.Command {
	var tradeMark string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "estimate-clear",
		Short: "门店商品估清",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tradeMark == "" {
				return fmt.Errorf("--trade-mark 必填")
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将对商品 %s 执行估清\n", tradeMark)
				return nil
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			_, err = api.SellOut(cmd.Context(), profile.ShopCode, tradeMark, 1)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已估清商品 %s\n", tradeMark)
			return nil
		},
	}

	cmd.Flags().StringVar(&tradeMark, "trade-mark", "", "商品标识")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCancelEstimateClear(f *cmdutil.Factory) *cobra.Command {
	var tradeMark string
	var tradeMarks []string
	var saleChannel int
	var saleType int
	var isAllEmpty int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "cancel-estimate-clear",
		Short: "取消门店商品估清",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			if tradeMark == "" && len(tradeMarks) == 0 {
				return fmt.Errorf("--trade-mark 或 --trade-marks 至少一个必填")
			}
			params := client.FillUpParams{
				IsAllEmpty:    isAllEmpty,
				SaleChannel:   saleChannel,
				SaleType:      saleType,
				StoreCode:     profile.ShopCode,
				TradeMark:     tradeMark,
				TradeMarkList: tradeMarks,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将取消估清: tradeMark=%s tradeMarks=%v\n", tradeMark, tradeMarks)
				return nil
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			resp, err := api.FillUp(cmd.Context(), params)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交取消估清，成功 %d 条，失败 %d 条\n", len(resp.SuccessList), len(resp.FailList))
			return nil
		},
	}

	cmd.Flags().StringVar(&tradeMark, "trade-mark", "", "单个商品标识")
	cmd.Flags().StringSliceVar(&tradeMarks, "trade-marks", nil, "商品标识列表")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "售卖渠道")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "售卖类型")
	cmd.Flags().IntVar(&isAllEmpty, "is-all-empty", 0, "估清标识，0=全部渠道和类型 1=指定渠道和类型")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}
