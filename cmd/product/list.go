package product

import (
	"fmt"
	"strings"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	var name string
	var saleChannel, saleType, page, pageSize int

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "商品列表",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
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
			result, err := api.List(cmd.Context(), profile.ShopCode, name, saleChannel, saleType, page, pageSize)
			if err != nil {
				return err
			}

			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}

			headers := []string{"ID", "名称", "价格", "状态", "库存", "分类"}
			rows := make([][]string, len(result.Data))
			for i, p := range result.Data {
				status := "下架"
				if p.Status == 10 {
					status = "上架"
				}
				price := fmt.Sprintf("%.2f", float64(p.ShowPriceLow)/100)
				inventory := "-"
				if len(p.GoodsSkuList) > 0 {
					inventory = fmt.Sprintf("%d", int(p.GoodsSkuList[0].Inventory))
				}
				category := strings.Join(p.CategoryNameList, "/")
				rows[i] = []string{
					fmt.Sprintf("%d", p.ID),
					p.Name,
					price,
					status,
					inventory,
					category,
				}
			}

			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			if err := fmtr.Write(result.Data, headers, rows); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, page, pageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "按名称搜索")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "销售渠道 (1=堂食, 2=外卖)")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "商品类型 (1=普通, 2=套餐)")
	cmd.Flags().IntVar(&page, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页数量")

	return cmd
}
