package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var name, className, tradeNo string
	var price float64
	var stock, saleType int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "update <tradeMark>",
		Short: "更新商品",
		Long:  "通过 Sync API 更新商品信息",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tradeMark := args[0]

			goods := client.ShopGoods{
				TradeNo:    tradeNo,
				TradeName:  name,
				TradePrice: price,
				ClassName:  className,
				Stock:      stock,
			}

			// tradeNo defaults to tradeMark if not specified
			if goods.TradeNo == "" {
				goods.TradeNo = tradeMark
			}

			if goods.TradeName == "" {
				return fmt.Errorf("--name 必填")
			}

			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将要更新商品 %s:\n", tradeMark)
				fmt.Fprintf(f.IOStreams.Out, "  名称: %s\n", goods.TradeName)
				fmt.Fprintf(f.IOStreams.Out, "  价格: %.2f\n", goods.TradePrice)
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

			if saleType == 0 {
				saleType = 1
			}

			api := client.NewProductAPI(apiClient)
			if err := api.Sync(cmd.Context(), profile.ShopCode, saleType, []client.ShopGoods{goods}); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 商品已更新 (%s)\n", tradeMark)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "商品名称")
	cmd.Flags().Float64Var(&price, "price", 0, "价格")
	cmd.Flags().StringVar(&tradeNo, "trade-no", "", "商品编号")
	cmd.Flags().StringVar(&className, "class", "", "分类名称")
	cmd.Flags().IntVar(&stock, "stock", 0, "库存")
	cmd.Flags().IntVar(&saleType, "sale-type", 1, "商品类型 (1=普通, 2=套餐)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览")

	return cmd
}
