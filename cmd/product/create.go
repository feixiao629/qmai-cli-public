package product

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var name, tradeNo, className string
	var price float64
	var stock, saleType int
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建商品",
		Long: `通过开放平台 Sync API 创建商品。

示例:
  qmai product create --name "美式咖啡" --price 28.00 --trade-no P001
  qmai product create --from-json goods.json
  qmai product create --name "拿铁" --price 32.00 --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var goodsList []client.ShopGoods

			if fromJSON != "" {
				data, err := os.ReadFile(fromJSON)
				if err != nil {
					return fmt.Errorf("读取 JSON 文件失败: %w", err)
				}
				if err := json.Unmarshal(data, &goodsList); err != nil {
					return fmt.Errorf("解析 JSON 失败: %w", err)
				}
			} else {
				if name == "" {
					return fmt.Errorf("--name 必填")
				}
				if price <= 0 {
					return fmt.Errorf("--price 必须大于 0")
				}
				goodsList = []client.ShopGoods{{
					TradeNo:    tradeNo,
					TradeName:  name,
					TradePrice: price,
					ClassName:  className,
					Stock:      stock,
				}}
			}

			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将要同步 %d 个商品:\n", len(goodsList))
				for _, g := range goodsList {
					fmt.Fprintf(f.IOStreams.Out, "  %s (%.2f)\n", g.TradeName, g.TradePrice)
				}
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
				saleType = 1 // default: 普通商品
			}

			api := client.NewProductAPI(apiClient)
			if err := api.Sync(cmd.Context(), profile.ShopCode, saleType, goodsList); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 已同步 %d 个商品\n", len(goodsList))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "商品名称")
	cmd.Flags().Float64Var(&price, "price", 0, "价格")
	cmd.Flags().StringVar(&tradeNo, "trade-no", "", "商品编号")
	cmd.Flags().StringVar(&className, "class", "", "分类名称")
	cmd.Flags().IntVar(&stock, "stock", 0, "库存")
	cmd.Flags().IntVar(&saleType, "sale-type", 1, "商品类型 (1=普通, 2=套餐)")
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "从 JSON 文件批量创建")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际创建")

	return cmd
}
