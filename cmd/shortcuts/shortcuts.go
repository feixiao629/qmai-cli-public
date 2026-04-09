package shortcuts

import (
	"fmt"
	"strconv"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// RegisterShortcuts adds shortcut subcommands under the appropriate parent commands
func RegisterShortcuts(productCmd *cobra.Command, f *cmdutil.Factory) {
	productCmd.AddCommand(newProductQuickAdd(f))
	productCmd.AddCommand(newProductOnSale(f))
	productCmd.AddCommand(newProductOffSale(f))
	productCmd.AddCommand(newProductPriceAdjust(f))
}

// --- Product Shortcuts ---

func newProductQuickAdd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "+quick-add <name> <price>",
		Short: "快速添加商品",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			price, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("价格格式错误: %s", args[1])
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
			goods := []client.ShopGoods{{
				TradeName:  name,
				TradePrice: price,
			}}
			if err := api.Sync(cmd.Context(), profile.ShopCode, 1, goods); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 商品已创建: %s (%.2f)\n", name, price)
			return nil
		},
	}
}

func newProductOnSale(f *cmdutil.Factory) *cobra.Command {
	var tradeMarks []string
	var saleChannel int

	cmd := &cobra.Command{
		Use:   "+on-sale",
		Short: "批量上架商品",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(tradeMarks) == 0 {
				return fmt.Errorf("--trade-marks 必填")
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
				return fmt.Errorf("未配置门店编码")
			}

			api := client.NewProductAPI(apiClient)
			if err := api.BatchUp(cmd.Context(), profile.ShopCode, tradeMarks, saleChannel); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已上架 %d 个商品\n", len(tradeMarks))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tradeMarks, "trade-marks", nil, "商品编码列表")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "销售渠道")
	return cmd
}

func newProductOffSale(f *cmdutil.Factory) *cobra.Command {
	var tradeMarks []string
	var saleChannel int

	cmd := &cobra.Command{
		Use:   "+off-sale",
		Short: "批量下架商品",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(tradeMarks) == 0 {
				return fmt.Errorf("--trade-marks 必填")
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
				return fmt.Errorf("未配置门店编码")
			}

			api := client.NewProductAPI(apiClient)
			if err := api.BatchDown(cmd.Context(), profile.ShopCode, tradeMarks, saleChannel); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已下架 %d 个商品\n", len(tradeMarks))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tradeMarks, "trade-marks", nil, "商品编码列表")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "销售渠道")
	return cmd
}

func newProductPriceAdjust(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "+price-adjust <adjustment>",
		Short: "快速调价",
		Long:  "快速调整商品价格，如: +price-adjust +10%",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(f.IOStreams.Out, "请先预览: qmai product batch-price --adjust %s --dry-run\n", args[0])
			return nil
		},
	}

	return cmd
}
