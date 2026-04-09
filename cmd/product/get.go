package product

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id|name>",
		Short: "查看商品详情",
		Long:  "通过商品列表接口查询单个商品详情（按 ID 或名称匹配）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyword := args[0]

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
			// Try name search first for better matching
			result, err := api.List(cmd.Context(), profile.ShopCode, keyword, 0, 0, 1, 50)
			if err != nil {
				return err
			}

			var found *client.OpenProduct
			for _, p := range result.Data {
				idStr := strconv.FormatInt(p.ID, 10)
				if idStr == keyword || p.Name == keyword {
					p := p
					found = &p
					break
				}
			}

			// If exact match not found, use first result from name search
			if found == nil && len(result.Data) > 0 {
				found = &result.Data[0]
			}

			if found == nil {
				return fmt.Errorf("未找到商品: %s", keyword)
			}

			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}

			status := "下架"
			if found.Status == 10 {
				status = "上架"
			}
			price := fmt.Sprintf("%.2f", float64(found.ShowPriceLow)/100)
			inventory := "-"
			if len(found.GoodsSkuList) > 0 {
				inventory = fmt.Sprintf("%d", int(found.GoodsSkuList[0].Inventory))
			}
			category := strings.Join(found.CategoryNameList, "/")

			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(found,
				[]string{"ID", "名称", "价格", "状态", "库存", "分类"},
				[][]string{{
					fmt.Sprintf("%d", found.ID),
					found.Name,
					price,
					status,
					inventory,
					category,
				}},
			)
		},
	}
}
