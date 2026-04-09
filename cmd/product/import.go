package product

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdImport(f *cmdutil.Factory) *cobra.Command {
	var file string
	var saleType int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "import",
		Short: "批量导入商品",
		Long: `从 CSV 或 JSON 文件批量导入商品（通过 Sync API）。

CSV 格式 (首行为表头):
  trade_no,trade_name,trade_price,class_name,stock

示例:
  qmai product import --file products.csv --dry-run
  qmai product import --file products.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file 必填")
			}

			goodsList, err := parseGoodsFile(file)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将要导入 %d 个商品:\n", len(goodsList))
				for i, g := range goodsList {
					if i >= 10 {
						fmt.Fprintf(f.IOStreams.Out, "  ... 还有 %d 个\n", len(goodsList)-10)
						break
					}
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
				saleType = 1
			}

			api := client.NewProductAPI(apiClient)
			if err := api.Sync(cmd.Context(), profile.ShopCode, saleType, goodsList); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "✓ 已导入 %d 个商品\n", len(goodsList))
			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "导入文件路径 (CSV/JSON)")
	cmd.Flags().IntVar(&saleType, "sale-type", 1, "商品类型 (1=普通, 2=套餐)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览导入")

	return cmd
}

func parseGoodsFile(path string) ([]client.ShopGoods, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return parseGoodsJSON(path)
	case ".csv":
		return parseGoodsCSV(path)
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s (使用 .csv 或 .json)", ext)
	}
}

func parseGoodsJSON(path string) ([]client.ShopGoods, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var goods []client.ShopGoods
	if err := json.Unmarshal(data, &goods); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}
	return goods, nil
}

func parseGoodsCSV(path string) ([]client.ShopGoods, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSV 解析失败: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV 文件为空或只有表头")
	}

	header := records[0]
	idx := make(map[string]int)
	for i, h := range header {
		idx[strings.TrimSpace(strings.ToLower(h))] = i
	}

	var goods []client.ShopGoods
	for i, row := range records[1:] {
		g := client.ShopGoods{}
		if col, ok := idx["trade_name"]; ok && col < len(row) {
			g.TradeName = row[col]
		}
		if col, ok := idx["trade_no"]; ok && col < len(row) {
			g.TradeNo = row[col]
		}
		if col, ok := idx["trade_price"]; ok && col < len(row) {
			price, err := strconv.ParseFloat(row[col], 64)
			if err != nil {
				return nil, fmt.Errorf("第 %d 行价格格式错误: %s", i+2, row[col])
			}
			g.TradePrice = price
		}
		if col, ok := idx["class_name"]; ok && col < len(row) {
			g.ClassName = row[col]
		}
		if col, ok := idx["stock"]; ok && col < len(row) {
			stock, _ := strconv.Atoi(row[col])
			g.Stock = stock
		}
		if g.TradeName == "" {
			return nil, fmt.Errorf("第 %d 行缺少 trade_name", i+2)
		}
		goods = append(goods, g)
	}

	return goods, nil
}
