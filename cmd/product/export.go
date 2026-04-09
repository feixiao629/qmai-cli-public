package product

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdExport(f *cmdutil.Factory) *cobra.Command {
	var file, name string
	var saleChannel, saleType int

	cmd := &cobra.Command{
		Use:   "export",
		Short: "导出商品",
		Long: `导出商品到 CSV 或 JSON 文件。

示例:
  qmai product export --file products.csv
  qmai product export --file products.json --name 咖啡`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file 必填")
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
			result, err := api.List(cmd.Context(), profile.ShopCode, name, saleChannel, saleType, 1, 50)
			if err != nil {
				return err
			}

			ext := strings.ToLower(filepath.Ext(file))
			switch ext {
			case ".json":
				data, err := json.MarshalIndent(result.Data, "", "  ")
				if err != nil {
					return err
				}
				if err := os.WriteFile(file, data, 0o644); err != nil {
					return err
				}
			case ".csv":
				out, err := os.Create(file)
				if err != nil {
					return err
				}
				defer out.Close()

				w := csv.NewWriter(out)
				w.Write([]string{"id", "name", "price", "status", "inventory", "category"})
				for _, p := range result.Data {
					status := "下架"
					if p.Status == 10 {
						status = "上架"
					}
					inventory := "-"
					if len(p.GoodsSkuList) > 0 {
						inventory = fmt.Sprintf("%d", int(p.GoodsSkuList[0].Inventory))
					}
					w.Write([]string{
						fmt.Sprintf("%d", p.ID),
						p.Name,
						fmt.Sprintf("%.2f", float64(p.ShowPriceLow)/100),
						status,
						inventory,
						strings.Join(p.CategoryNameList, "/"),
					})
				}
				w.Flush()
				if err := w.Error(); err != nil {
					return err
				}
			default:
				return fmt.Errorf("不支持的格式: %s (使用 .csv 或 .json)", ext)
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 已导出 %d 个商品到 %s\n", len(result.Data), file)
			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "导出文件路径")
	cmd.Flags().StringVar(&name, "name", "", "按名称过滤")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "销售渠道")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "商品类型")

	return cmd
}
