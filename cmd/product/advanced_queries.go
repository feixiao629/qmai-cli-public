package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdAttachList(f *cmdutil.Factory) *cobra.Command {
	var saleChannel int
	var saleType int
	var stockStatus int

	cmd := &cobra.Command{
		Use:   "attach-list",
		Short: "查询门店加料列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			result, err := api.ListAttachGoods(cmd.Context(), client.GoodsAttachListParams{
				SaleChannel: saleChannel,
				SaleType:    saleType,
				ShopCode:    profile.ShopCode,
				StockStatus: stockStatus,
			})
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result))
			for _, item := range result {
				rows = append(rows, []string{
					item.ID,
					item.AttachCode,
					item.Name,
					item.ShowPrice,
					item.Inventory,
					fmt.Sprintf("%d", item.Status),
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"ID", "加料编码", "名称", "展示价格", "库存", "状态"}, rows)
		},
	}

	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "售卖渠道，不传默认小程序 3")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "售卖类型 1=堂食 2=外卖")
	cmd.Flags().IntVar(&stockStatus, "stock-status", 0, "库存状态 0=无库存 1=有库存")
	return cmd
}

func newCmdListWithPractice(f *cmdutil.Factory) *cobra.Command {
	var name string
	var pageNo int
	var pageSize int
	var saleChannel int
	var saleType int
	var saleMethod int
	var status string
	var filterDownAttach bool
	var includeProps []string

	cmd := &cobra.Command{
		Use:   "list-with-practice",
		Short: "查询门店商品列表(支持做法)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			if saleChannel == 0 || saleType == 0 || status == "" {
				return fmt.Errorf("--sale-channel、--sale-type、--status 必填")
			}
			if pageNo == 0 {
				pageNo = 1
			}
			if pageSize == 0 {
				pageSize = 10
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			result, err := api.ListShopGoodsWithPractice(cmd.Context(), client.ShopGoodsListParams{
				FilterDownAttach: filterDownAttach,
				IncludeProps:     includeProps,
				Name:             name,
				PageNo:           pageNo,
				PageSize:         pageSize,
				SaleChannel:      saleChannel,
				SaleMethod:       saleMethod,
				SaleType:         saleType,
				ShopCode:         profile.ShopCode,
				Status:           status,
			})
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					fmt.Sprintf("%d", item.ID),
					fmt.Sprintf("%d", item.GoodsID),
					item.Name,
					fmt.Sprintf("%d", item.Status),
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			if err := fmtr.Write(result, []string{"门店商品ID", "商品库ID", "名称", "状态"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, pageNo, pageSize)
			return nil
		},
	}

	cmd.Flags().BoolVar(&filterDownAttach, "filter-down-attach", false, "是否过滤下架加料商品")
	cmd.Flags().StringVar(&name, "name", "", "商品名称")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "售卖渠道")
	cmd.Flags().IntVar(&saleMethod, "sale-method", 0, "售卖方式 0=正常售卖 1=仅套餐中售卖")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "售卖类型")
	cmd.Flags().StringVar(&status, "status", "", "上下架状态 10=上架 20=下架")
	cmd.Flags().StringSliceVar(&includeProps, "include-properties", nil, "控制返回属性，如 ATTACH,PRACTICE,SKU")
	return cmd
}

func newCmdRealtime(f *cmdutil.Factory) *cobra.Command {
	var storeID int64
	var goodsIDs []int64
	var saleType int
	var saleChannel int

	cmd := &cobra.Command{
		Use:   "realtime",
		Short: "门店商品实时数据查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if storeID == 0 || len(goodsIDs) == 0 || saleType == 0 || saleChannel == 0 {
				return fmt.Errorf("--store-id、--goods-ids、--sale-type、--sale-channel 必填")
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			result, err := api.ListRealtimeGoods(cmd.Context(), client.GoodsRealtimeParams{
				StoreID:     storeID,
				GoodsIDList: goodsIDs,
				SaleType:    saleType,
				SaleChannel: saleChannel,
			})
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result))
			for _, item := range result {
				rows = append(rows, []string{
					fmt.Sprintf("%d", item.ID),
					fmt.Sprintf("%d", item.GoodsID),
					fmt.Sprintf("%d", item.SaleChannel),
					fmt.Sprintf("%d", item.SaleType),
					fmt.Sprintf("%d", item.Status),
					fmt.Sprintf("%d", item.PackingFee),
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"门店商品ID", "商品库ID", "售卖渠道", "售卖类型", "状态", "包装费(分)"}, rows)
		},
	}

	cmd.Flags().Int64Var(&storeID, "store-id", 0, "门店 ID")
	cmd.Flags().Int64SliceVar(&goodsIDs, "goods-ids", nil, "商品库商品 ID 列表")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "售卖类型")
	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "售卖渠道")
	return cmd
}

func newCmdEnergy(f *cmdutil.Factory) *cobra.Command {
	var goodsID int64
	var storeID int64

	cmd := &cobra.Command{
		Use:   "energy",
		Short: "查询商品能量值",
		RunE: func(cmd *cobra.Command, args []string) error {
			if goodsID == 0 || storeID == 0 {
				return fmt.Errorf("--goods-id、--store-id 必填")
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			result, err := api.GetGoodsEnergy(cmd.Context(), goodsID, storeID)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result))
			for _, item := range result {
				rows = append(rows, []string{
					fmt.Sprintf("%d", item.GoodsID),
					fmt.Sprintf("%d", item.SkuID),
					fmt.Sprintf("%.2f", item.EnergyValue),
					item.Unit,
					item.GradingIdentification,
					item.GradingInstructions,
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"商品库ID", "SKU ID", "能量值", "单位", "分级标识", "分级说明"}, rows)
		},
	}

	cmd.Flags().Int64Var(&goodsID, "goods-id", 0, "商品库商品 ID")
	cmd.Flags().Int64Var(&storeID, "store-id", 0, "门店 ID")
	return cmd
}
