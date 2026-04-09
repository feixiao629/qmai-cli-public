package order

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdOrder(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order",
		Short: "订单服务",
		Long:  "订单查询、订单上报、评价回复、制作单查询",
	}
	cmd.AddCommand(newCmdQuery(f))
	cmd.AddCommand(newCmdReport(f))
	cmd.AddCommand(newCmdReview(f))
	return cmd
}

func newOrderAPI(f *cmdutil.Factory) (*client.OrderAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewOrderAPI(apiClient), nil
}

func newCmdQuery(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "query", Short: "订单查询"}
	cmd.AddCommand(newCmdUserOrders(f))
	cmd.AddCommand(newCmdDetail(f))
	cmd.AddCommand(newCmdRechargeOrders(f))
	cmd.AddCommand(newCmdRechargeRefundOrders(f))
	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdMemberOrdered(f))
	cmd.AddCommand(newCmdProductionRecords(f))
	return cmd
}

func newCmdReport(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "report", Short: "订单上报"}
	cmd.AddCommand(newCmdUpload(f))
	cmd.AddCommand(newCmdRefundUp(f))
	cmd.AddCommand(newCmdCompletedBatchUpload(f))
	cmd.AddCommand(newCmdRefundedBatchUpload(f))
	return cmd
}

func newCmdReview(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "review", Short: "订单评价"}
	cmd.AddCommand(newCmdReply(f))
	return cmd
}

func newCmdUserOrders(f *cmdutil.Factory) *cobra.Command {
	var orderAtStart, orderAtEnd, shopCode string
	var size int
	var sourceList, statusList []int
	var totalAmountStart int
	var userID int64
	cmd := &cobra.Command{
		Use:   "user-orders",
		Short: "查询用户订单列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderAtStart == "" || orderAtEnd == "" || size == 0 || userID == 0 {
				return fmt.Errorf("--order-at-start、--order-at-end、--size、--user-id 必填")
			}
			params := map[string]interface{}{
				"orderAtStart": orderAtStart,
				"orderAtEnd":   orderAtEnd,
				"size":         size,
				"userId":       userID,
			}
			if shopCode != "" {
				params["shopCode"] = shopCode
			}
			if len(sourceList) > 0 {
				params["sourceList"] = sourceList
			}
			if len(statusList) > 0 {
				params["statusList"] = statusList
			}
			if totalAmountStart != 0 {
				params["totalAmountStart"] = totalAmountStart
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetUserOrderList(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.OrderNo, item.ShopCode, strconv.Itoa(item.Status), strconv.Itoa(item.Source),
					fmt.Sprintf("%.2f", float64(item.TotalAmount)/100), fmt.Sprintf("%.2f", float64(item.ActualAmount)/100), item.CreatedAt,
				})
			}
			if err := writeRows(f, result, []string{"订单号", "门店编码", "状态", "来源", "总金额(元)", "实付(元)", "下单时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&orderAtStart, "order-at-start", "", "下单开始时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().StringVar(&orderAtEnd, "order-at-end", "", "下单结束时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	cmd.Flags().IntVar(&size, "size", 0, "订单数量阈值，最大100")
	cmd.Flags().IntSliceVar(&sourceList, "source-list", nil, "订单来源列表")
	cmd.Flags().IntSliceVar(&statusList, "status-list", nil, "订单状态列表")
	cmd.Flags().IntVar(&totalAmountStart, "total-amount-start", 0, "订单总金额起始，单位分")
	cmd.Flags().Int64Var(&userID, "user-id", 0, "用户ID")
	return cmd
}

func newCmdDetail(f *cmdutil.Factory) *cobra.Command {
	var bizType int
	var orderNo string
	var userID int64
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "查询订单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if bizType == 0 || orderNo == "" {
				return fmt.Errorf("--biz-type 和 --order-no 必填")
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetOrderDetail(cmd.Context(), bizType, orderNo, userID)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"订单号", "门店", "总金额(元)", "实付(元)", "优惠(元)", "下单时间", "完成时间"}, [][]string{{
				result.OrderNo, result.ShopName, fmt.Sprintf("%.2f", float64(result.TotalAmount)/100), fmt.Sprintf("%.2f", float64(result.ActualAmount)/100), fmt.Sprintf("%.2f", float64(result.DiscountAmount)/100), result.CreatedAt, result.CompletedAt,
			}})
		},
	}
	cmd.Flags().IntVar(&bizType, "biz-type", 0, "业务类型，5=新饮食 4=新休闲")
	cmd.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	cmd.Flags().Int64Var(&userID, "user-id", 0, "会员ID")
	return cmd
}

func newRechargeBaseCmd(f *cmdutil.Factory, use, short string, refund bool) *cobra.Command {
	var shopCode, startTime, endTime string
	var pageNo, pageSize int
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopCode == "" || startTime == "" || endTime == "" || pageNo == 0 || pageSize == 0 {
				return fmt.Errorf("--shop-code、--start-time、--end-time、--page、--page-size 必填")
			}
			params := map[string]interface{}{"shopCode": shopCode, "startTime": startTime, "endTime": endTime, "pageNo": pageNo, "pageSize": pageSize}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			if refund {
				result, err := api.ListRechargeRefundOrders(cmd.Context(), params)
				if err != nil {
					return err
				}
				rows := make([][]string, 0, len(*result))
				for _, item := range *result {
					rows = append(rows, []string{item.OrderNo, item.RefundNo, item.ShopCode, strconv.Itoa(item.RefundStatus), fmt.Sprintf("%.2f", float64(item.RefundAmount)/100), item.DealAt})
				}
				return writeRows(f, result, []string{"订单号", "退款单号", "门店编码", "退款状态", "退款金额(元)", "处理时间"}, rows)
			}
			result, err := api.ListRechargeOrders(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.OrderNo, item.ShopCode, strconv.Itoa(item.Status), fmt.Sprintf("%.2f", float64(item.TotalAmount)/100), fmt.Sprintf("%.2f", float64(item.ActualAmount)/100), item.CreatedAt})
			}
			return writeRows(f, result, []string{"订单号", "门店编码", "状态", "总金额(元)", "实付(元)", "下单时间"}, rows)
		},
	}
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	cmd.Flags().StringVar(&startTime, "start-time", "", "开始时间 yyyyMMdd")
	cmd.Flags().StringVar(&endTime, "end-time", "", "结束时间 yyyyMMdd")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	return cmd
}

func newCmdRechargeOrders(f *cmdutil.Factory) *cobra.Command {
	return newRechargeBaseCmd(f, "recharge-orders", "查询门店储值消费订单", false)
}

func newCmdRechargeRefundOrders(f *cmdutil.Factory) *cobra.Command {
	return newRechargeBaseCmd(f, "recharge-refunds", "查询门店储值消费退款单", true)
}

func newCmdStatus(f *cmdutil.Factory) *cobra.Command {
	var orderNo string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "查询订单状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderNo == "" {
				return fmt.Errorf("--order-no 必填")
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetOrderStatus(cmd.Context(), orderNo)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"订单号", "状态", "取餐号", "用户ID"}, [][]string{{orderNo, result.Status, result.StoreOrderNo, strconv.FormatInt(result.UserID, 10)}})
		},
	}
	cmd.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	return cmd
}

func newCmdMemberOrdered(f *cmdutil.Factory) *cobra.Command {
	var bizType int
	var excludeStatusList, payStatusList, sourceList []int
	var userID string
	cmd := &cobra.Command{
		Use:   "member-ordered",
		Short: "查询用户是否下过订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if bizType == 0 || userID == "" {
				return fmt.Errorf("--biz-type 和 --user-id 必填")
			}
			params := map[string]interface{}{"bizType": bizType, "userId": userID}
			if len(excludeStatusList) > 0 {
				params["excludeStatusList"] = excludeStatusList
			}
			if len(payStatusList) > 0 {
				params["payStatusList"] = payStatusList
			}
			if len(sourceList) > 0 {
				params["sourceList"] = sourceList
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CheckMemberOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			label := "非新用户"
			if int(*result) == 0 {
				label = "新用户"
			}
			return writeRows(f, result, []string{"用户ID", "结果值", "说明"}, [][]string{{userID, strconv.Itoa(int(*result)), label}})
		},
	}
	cmd.Flags().IntVar(&bizType, "biz-type", 0, "业务类型，5=新饮食")
	cmd.Flags().IntSliceVar(&excludeStatusList, "exclude-status-list", nil, "排除订单状态列表")
	cmd.Flags().IntSliceVar(&payStatusList, "pay-status-list", nil, "支付状态列表")
	cmd.Flags().IntSliceVar(&sourceList, "source-list", nil, "来源列表")
	cmd.Flags().StringVar(&userID, "user-id", "", "用户ID")
	return cmd
}

func newCmdProductionRecords(f *cmdutil.Factory) *cobra.Command {
	var orderNo string
	var orderNoList []string
	var pdcType int
	var storeID int64
	var userID int64
	var userIDList []int64
	cmd := &cobra.Command{
		Use:   "production-records",
		Short: "查询商品制作单列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if userID == 0 {
				return fmt.Errorf("--user-id 必填")
			}
			params := map[string]interface{}{"userId": userID}
			if orderNo != "" {
				params["orderNo"] = orderNo
			}
			if len(orderNoList) > 0 {
				params["orderNoList"] = orderNoList
			}
			if pdcType != 0 {
				params["pdcType"] = pdcType
			}
			if storeID != 0 {
				params["storeId"] = storeID
			}
			if len(userIDList) > 0 {
				params["userIdList"] = userIDList
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListProductionRecords(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.OrderNo, strconv.Itoa(item.PdcType), strconv.Itoa(item.Status), item.ProductionSign, item.ProductionAtStart, item.ProductionAtEnd})
			}
			return writeRows(f, result, []string{"订单号", "制作单类型", "状态", "制作单标识", "开始时间", "结束时间"}, rows)
		},
	}
	cmd.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	cmd.Flags().StringSliceVar(&orderNoList, "order-no-list", nil, "订单号列表")
	cmd.Flags().IntVar(&pdcType, "pdc-type", 0, "制作单类型，0=裱花 1=现烤")
	cmd.Flags().Int64Var(&storeID, "store-id", 0, "门店ID")
	cmd.Flags().Int64Var(&userID, "user-id", 0, "用户ID")
	cmd.Flags().Int64SliceVar(&userIDList, "user-id-list", nil, "用户ID列表")
	return cmd
}

func newCmdUpload(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "订单上传",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.OrderUpload(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"订单号"}, [][]string{{result.OrderNo}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	return cmd
}

func newCmdRefundUp(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	cmd := &cobra.Command{
		Use:   "refund-up",
		Short: "退款订单上报",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.RefundOrderUp(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"订单ID", "订单号", "退款单号"}, [][]string{{strconv.FormatInt(result.OrderID, 10), result.OrderNo, result.RefundOrderNo}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	return cmd
}

func newCmdCompletedBatchUpload(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "completed-batch",
		Short: "提交已完成订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"提交已完成订单", fromJSON}})
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.OrderBatchUpload(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"成功数", "失败数", "失败订单号"}, [][]string{{
				strconv.Itoa(result.Success),
				strconv.Itoa(result.Error),
				fmt.Sprint(result.ErrorOrderNo),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdRefundedBatchUpload(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "refunded-batch",
		Short: "提交已退款订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"提交已退款订单", fromJSON}})
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			result, err := api.RefundOrderBatchUpload(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"成功数", "失败数", "失败订单号"}, [][]string{{
				strconv.Itoa(result.Success),
				strconv.Itoa(result.Error),
				fmt.Sprint(result.ErrorOrderNo),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdReply(f *cmdutil.Factory) *cobra.Command {
	var orderNo, replyAt, sellerReplyInfo string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "reply",
		Short: "回复订单评价",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderNo == "" || replyAt == "" || sellerReplyInfo == "" {
				return fmt.Errorf("--order-no、--reply-at、--seller-reply-info 必填")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将回复订单 %s 评价\n", orderNo)
				return nil
			}
			api, err := newOrderAPI(f)
			if err != nil {
				return err
			}
			if err := api.ReplyUserComment(cmd.Context(), orderNo, replyAt, sellerReplyInfo); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已回复订单 %s 评价\n", orderNo)
			return nil
		},
	}
	cmd.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	cmd.Flags().StringVar(&replyAt, "reply-at", "", "回复时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().StringVar(&sellerReplyInfo, "seller-reply-info", "", "回复内容")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func loadJSONFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read json file: %w", err)
	}
	var params map[string]interface{}
	if err := json.Unmarshal(data, &params); err != nil {
		return nil, fmt.Errorf("parse json file: %w", err)
	}
	return params, nil
}

func writeRows(f *cmdutil.Factory, data interface{}, headers []string, rows [][]string) error {
	format, err := output.ParseFormat(f.EffectiveFormat())
	if err != nil {
		return err
	}
	fmtr := output.NewFormatter(f.IOStreams.Out, format)
	return fmtr.Write(data, headers, rows)
}
