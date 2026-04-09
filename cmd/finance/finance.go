package finance

import (
	"fmt"
	"strconv"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdFinance(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finance",
		Short: "财务服务",
		Long:  "财务对账、支付账单、营业统计、结账方式与商品销售汇总",
	}
	cmd.AddCommand(newCmdStatement(f))
	cmd.AddCommand(newCmdStats(f))
	return cmd
}

func newFinanceAPI(f *cmdutil.Factory) (*client.FinanceAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewFinanceAPI(apiClient), nil
}

func newCmdStatement(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "statement", Short: "对账与账单"}
	cmd.AddCommand(newCmdSplitFlows(f))
	cmd.AddCommand(newCmdWechatBills(f))
	cmd.AddCommand(newCmdAlipayBills(f))
	cmd.AddCommand(newCmdWechatBillURL(f))
	return cmd
}

func newCmdStats(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "stats", Short: "财务统计与字典"}
	cmd.AddCommand(newCmdBusinessSummary(f))
	cmd.AddCommand(newCmdOrderTypes(f))
	cmd.AddCommand(newCmdSettleScenes(f))
	cmd.AddCommand(newCmdItemTurnover(f))
	return cmd
}

func newCmdSplitFlows(f *cmdutil.Factory) *cobra.Command {
	var createdAtStart, createdAtEnd string
	var pageNum, pageSize int
	var shopCodes []string
	cmd := &cobra.Command{
		Use:   "split-flows",
		Short: "查询分账明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createdAtStart == "" || createdAtEnd == "" || pageNum == 0 || pageSize == 0 {
				return fmt.Errorf("--created-at-start、--created-at-end、--page、--page-size 必填")
			}
			params := map[string]interface{}{
				"createdAtStart": createdAtStart,
				"createdAtEnd":   createdAtEnd,
				"pageNum":        pageNum,
				"pageSize":       pageSize,
			}
			if len(shopCodes) > 0 {
				params["shopCodes"] = shopCodes
			}
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetSplitOrderFlows(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.OriginOrderNo, item.RefundOrderNo, item.ShopCode, item.ShopName,
					item.SplitType, item.SplitStatus, item.SplitAmount, item.CreatedAt,
				})
			}
			if err := writeRows(f, result, []string{"订单号", "退款单号", "门店编码", "门店名称", "分账类型", "分账状态", "预计分账金额", "创建时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&createdAtStart, "created-at-start", "", "订单创建开始时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().StringVar(&createdAtEnd, "created-at-end", "", "订单创建结束时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	cmd.Flags().StringSliceVar(&shopCodes, "shop-codes", nil, "门店编码列表")
	return cmd
}

func newCmdWechatBills(f *cmdutil.Factory) *cobra.Command {
	var billDate string
	var pageNum, pageSize int
	var shopCodes []string
	cmd := &cobra.Command{
		Use:   "wechat-bills",
		Short: "查询微信支付账单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if billDate == "" || pageNum == 0 || pageSize == 0 {
				return fmt.Errorf("--bill-date、--page、--page-size 必填")
			}
			params := map[string]interface{}{
				"billDate": billDate,
				"pageNum":  pageNum,
				"pageSize": pageSize,
			}
			if len(shopCodes) > 0 {
				params["shopCodes"] = shopCodes
			}
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetWechatBills(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.TradeNo, item.ChannelNo, item.ShopCode, item.ShopName,
					fmt.Sprintf("%.2f", item.TradeAmount), fmt.Sprintf("%.2f", item.MerchantReceiptAmount),
					fmt.Sprintf("%.2f", item.RefundAmount), item.TradeTime,
				})
			}
			if err := writeRows(f, result, []string{"内部支付单号", "渠道交易号", "门店编码", "门店名称", "交易金额(元)", "商户实收(元)", "退款金额(元)", "交易时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&billDate, "bill-date", "", "账单日期 yyyy-MM-dd")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	cmd.Flags().StringSliceVar(&shopCodes, "shop-codes", nil, "门店编码列表")
	return cmd
}

func newCmdAlipayBills(f *cmdutil.Factory) *cobra.Command {
	var billDate, multiMark string
	var billType, pageNum, pageSize int
	var multiMarks, sellerIDs, subMchIDs []string
	var regionID int64
	cmd := &cobra.Command{
		Use:   "alipay-bills",
		Short: "查询支付宝支付账单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if billDate == "" || pageNum == 0 || pageSize == 0 {
				return fmt.Errorf("--bill-date、--page、--page-size 必填")
			}
			params := map[string]interface{}{
				"billDate": billDate,
				"pageNum":  pageNum,
				"pageSize": pageSize,
			}
			if billType != 0 {
				params["billType"] = billType
			}
			if multiMark != "" {
				params["multiMark"] = multiMark
			}
			if len(multiMarks) > 0 {
				params["multiMarks"] = multiMarks
			}
			if regionID != 0 {
				params["regionId"] = regionID
			}
			if len(sellerIDs) > 0 {
				params["sellerIds"] = sellerIDs
			}
			if len(subMchIDs) > 0 {
				params["subMchIds"] = subMchIDs
			}
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetAlipayBills(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for idx, item := range result.Data {
				rows = append(rows, []string{
					strconv.Itoa(idx + 1),
					string(item),
				})
			}
			if err := writeRows(f, result, []string{"序号", "原始记录"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条，总页数 %d\n", result.Total, result.TotalPage)
			return nil
		},
	}
	cmd.Flags().StringVar(&billDate, "bill-date", "", "交易日期")
	cmd.Flags().IntVar(&billType, "bill-type", 0, "账单类型，0=微信 1=支付宝")
	cmd.Flags().StringVar(&multiMark, "multi-mark", "", "单个门店编码")
	cmd.Flags().StringSliceVar(&multiMarks, "multi-marks", nil, "门店编码列表")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	cmd.Flags().Int64Var(&regionID, "region-id", 0, "机构ID")
	cmd.Flags().StringSliceVar(&sellerIDs, "seller-ids", nil, "支付宝商家ID列表")
	cmd.Flags().StringSliceVar(&subMchIDs, "sub-mch-ids", nil, "微信商户号列表")
	return cmd
}

func newCmdWechatBillURL(f *cmdutil.Factory) *cobra.Command {
	var billDate, shopCode string
	cmd := &cobra.Command{
		Use:   "wechat-bill-url",
		Short: "查询微信支付账单URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if billDate == "" {
				return fmt.Errorf("--bill-date 必填")
			}
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetWechatBillURL(cmd.Context(), billDate, shopCode)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"文件地址"}, [][]string{{result.FileURL}})
		},
	}
	cmd.Flags().StringVar(&billDate, "bill-date", "", "账单日期 yyyy-MM-dd")
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	return cmd
}

func newCmdBusinessSummary(f *cmdutil.Factory) *cobra.Command {
	var shopCode, startDate, endDate string
	var pageNo, pageSize int
	cmd := &cobra.Command{
		Use:   "business-summary",
		Short: "查询营业统计数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopCode == "" || startDate == "" || endDate == "" {
				return fmt.Errorf("--shop-code、--start-date、--end-date 必填")
			}
			params := map[string]interface{}{
				"shopCode":   shopCode,
				"start_date": startDate,
				"end_date":   endDate,
				"pageNo":     pageNo,
				"pageSize":   pageSize,
			}
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetBusinessSummary(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.ResultList))
			for _, item := range result.ResultList {
				rows = append(rows, []string{
					item.RecordTime, item.ShopCode, item.ShopName,
					fmt.Sprintf("%.2f", item.BusinessAmt), fmt.Sprintf("%.2f", item.IncomeAmt),
					fmt.Sprintf("%.2f", item.RefundAmt), strconv.Itoa(item.RefundNum),
				})
			}
			if err := writeRows(f, result, []string{"日期", "门店编码", "门店名称", "营业额(元)", "实收(元)", "退款金额(元)", "退款单数"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.TotalCount)
			return nil
		},
	}
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	cmd.Flags().StringVar(&startDate, "start-date", "", "开始日期 yyyy-MM-dd")
	cmd.Flags().StringVar(&endDate, "end-date", "", "结束日期 yyyy-MM-dd")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	return cmd
}

func newCmdOrderTypes(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "order-types",
		Short: "查询订单渠道",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListOrderTypes(cmd.Context())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.OrderType, item.OrderTypeName})
			}
			return writeRows(f, result, []string{"渠道类型", "渠道名称"}, rows)
		},
	}
}

func newCmdSettleScenes(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "settle-scenes",
		Short: "查询订单结账方式",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListSettleScenes(cmd.Context())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.Code, item.Name, item.SellerID})
			}
			return writeRows(f, result, []string{"结账方式编码", "结账方式名称", "店铺ID"}, rows)
		},
	}
}

func newCmdItemTurnover(f *cmdutil.Factory) *cobra.Command {
	var shopCode, startDate string
	var shopID int64
	var pageNo, pageSize int
	cmd := &cobra.Command{
		Use:   "item-turnover",
		Short: "查询门店商品销售汇总",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopCode == "" && shopID == 0 {
				return fmt.Errorf("--shop-code 和 --shop-id 至少一个必填")
			}
			if pageNo == 0 || pageSize == 0 {
				return fmt.Errorf("--page 和 --page-size 必填")
			}
			params := map[string]interface{}{
				"pageNo":   pageNo,
				"pageSize": pageSize,
			}
			if shopCode != "" {
				params["shopCode"] = shopCode
			}
			if shopID != 0 {
				params["shopId"] = shopID
			}
			if startDate != "" {
				params["start_date"] = startDate
			}
			api, err := newFinanceAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetStoreTurnover(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.ResultList))
			for _, item := range result.ResultList {
				rows = append(rows, []string{
					item.RecordDate, item.StoreCode, item.StoreName, item.Name,
					fmt.Sprintf("%.2f", item.Num), fmt.Sprintf("%.2f", item.ReceivableAmount),
					fmt.Sprintf("%.2f", item.ReceivedAmount), fmt.Sprintf("%.2f", item.RefundAmount),
				})
			}
			if err := writeRows(f, result, []string{"营业日期", "门店编码", "门店名称", "商品名称", "销量", "销售额(元)", "实收(元)", "退款金额(元)"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.TotalCount)
			return nil
		},
	}
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	cmd.Flags().Int64Var(&shopID, "shop-id", 0, "门店ID")
	cmd.Flags().StringVar(&startDate, "start-date", "", "开始日期 yyyy-MM-dd")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数，最大 200")
	return cmd
}

func writeRows(f *cmdutil.Factory, data interface{}, headers []string, rows [][]string) error {
	format, err := output.ParseFormat(f.EffectiveFormat())
	if err != nil {
		return err
	}
	fmtr := output.NewFormatter(f.IOStreams.Out, format)
	return fmtr.Write(data, headers, rows)
}
