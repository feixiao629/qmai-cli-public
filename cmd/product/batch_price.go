package product

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdBatchPrice(f *cmdutil.Factory) *cobra.Command {
	var adjust string
	var saleType int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "batch-price",
		Short: "批量调价",
		Long: `批量调整商品价格（先查询当前价格，计算新价格后通过 Sync API 更新）。

示例:
  qmai product batch-price --adjust +10% --dry-run
  qmai product batch-price --adjust -5%
  qmai product batch-price --adjust +2.00`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if adjust == "" {
				return fmt.Errorf("--adjust 必填 (如 +10%%, -5%%, +2.00)")
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
			result, err := api.List(cmd.Context(), profile.ShopCode, "", 0, 0, 1, 50)
			if err != nil {
				return err
			}

			if len(result.Data) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "没有匹配的商品")
				return nil
			}

			type priceUpdate struct {
				Product  client.OpenProduct
				OldPrice float64
				NewPrice float64
			}
			var updates []priceUpdate
			for _, p := range result.Data {
				oldPrice := float64(p.ShowPriceLow) / 100
				newPrice, err := applyPriceAdjust(oldPrice, adjust)
				if err != nil {
					return err
				}
				updates = append(updates, priceUpdate{Product: p, OldPrice: oldPrice, NewPrice: newPrice})
			}

			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将要调整 %d 个商品价格:\n", len(updates))
				for _, u := range updates {
					fmt.Fprintf(f.IOStreams.Out, "  %s: %.2f → %.2f\n", u.Product.Name, u.OldPrice, u.NewPrice)
				}
				return nil
			}

			// Build goods list with new prices
			goodsList := make([]client.ShopGoods, len(updates))
			for i, u := range updates {
				goodsList[i] = client.ShopGoods{
					TradeName:  u.Product.Name,
					TradePrice: u.NewPrice,
				}
			}

			if saleType == 0 {
				saleType = 1
			}

			if err := api.Sync(cmd.Context(), profile.ShopCode, saleType, goodsList); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 已调整 %d 个商品价格\n", len(updates))
			return nil
		},
	}

	cmd.Flags().StringVar(&adjust, "adjust", "", "调价幅度 (+10%, -5%, +2.00)")
	cmd.Flags().IntVar(&saleType, "sale-type", 1, "商品类型 (1=普通, 2=套餐)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览调价")

	return cmd
}

func applyPriceAdjust(price float64, adjust string) (float64, error) {
	adjust = strings.TrimSpace(adjust)
	if strings.HasSuffix(adjust, "%") {
		pctStr := strings.TrimSuffix(adjust, "%")
		pct, err := strconv.ParseFloat(pctStr, 64)
		if err != nil {
			return 0, fmt.Errorf("无效百分比: %s", adjust)
		}
		result := price * (1 + pct/100)
		return math.Round(result*100) / 100, nil
	}
	delta, err := strconv.ParseFloat(adjust, 64)
	if err != nil {
		return 0, fmt.Errorf("无效调价: %s", adjust)
	}
	result := price + delta
	if result < 0 {
		result = 0
	}
	return math.Round(result*100) / 100, nil
}
