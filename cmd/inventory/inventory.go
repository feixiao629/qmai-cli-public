package inventory

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

func NewCmdInventory(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inventory",
		Short: "进销存",
		Long:  "订货、调拨、退货、入库、出库和库存查询",
	}
	cmd.AddCommand(newCmdSales(f))
	cmd.AddCommand(newCmdStock(f))
	cmd.AddCommand(newCmdMaster(f))
	return cmd
}

func newInventoryAPI(f *cmdutil.Factory) (*client.InventoryAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewInventoryAPI(apiClient), nil
}

func newCmdSales(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "sales", Short: "销售管理"}
	cmd.AddCommand(newCmdDeclareReceive(f))
	cmd.AddCommand(newCmdDeclareDetail(f))
	cmd.AddCommand(newCmdRequireUpdate(f))
	cmd.AddCommand(newCmdRequireDeliver(f))
	cmd.AddCommand(newCmdRequireDetailUpdate(f))
	cmd.AddCommand(newCmdReturnCancel(f))
	cmd.AddCommand(newCmdReturnExamine(f))
	cmd.AddCommand(newCmdReturnReceipt(f))
	cmd.AddCommand(newCmdDeliveryArrive(f))
	cmd.AddCommand(newCmdRequireList(f))
	cmd.AddCommand(newCmdRequireCreate(f))
	cmd.AddCommand(newCmdTransferList(f))
	cmd.AddCommand(newCmdTransferDetail(f))
	cmd.AddCommand(newCmdReturnList(f))
	return cmd
}

func newCmdStock(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "stock", Short: "库存管理"}
	cmd.AddCommand(newCmdInboundCreate(f))
	cmd.AddCommand(newCmdInboundList(f))
	cmd.AddCommand(newCmdInboundUpdate(f))
	cmd.AddCommand(newCmdInboundFinish(f))
	cmd.AddCommand(newCmdOutboundCreate(f))
	cmd.AddCommand(newCmdOutboundList(f))
	cmd.AddCommand(newCmdOutboundUpdate(f))
	cmd.AddCommand(newCmdStockOccupy(f))
	cmd.AddCommand(newCmdStockRelease(f))
	cmd.AddCommand(newCmdStoreLedger(f))
	cmd.AddCommand(newCmdRealtimeList(f))
	return cmd
}

func newCmdMaster(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "master", Short: "基础资料管理"}
	cmd.AddCommand(newCmdInventoryAdjust(f))
	cmd.AddCommand(newCmdTransferAudit(f))
	cmd.AddCommand(newCmdProductCreate(f))
	cmd.AddCommand(newCmdCategoryCreate(f))
	cmd.AddCommand(newCmdUnitCreate(f))
	cmd.AddCommand(newCmdProductUpdate(f))
	cmd.AddCommand(newCmdSupplierCreate(f))
	cmd.AddCommand(newCmdProductDistribute(f))
	cmd.AddCommand(newCmdSupplierUpdate(f))
	cmd.AddCommand(newCmdMachiningCardBatchCreate(f))
	return cmd
}

func newCmdDeclareReceive(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "declare-receive",
		Short: "报货单接单/拒单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"报货单接单/拒单", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ReceiveDeclareOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdDeclareDetail(f *cmdutil.Factory) *cobra.Command {
	var declareNo, shopCode string
	cmd := &cobra.Command{
		Use:   "declare-detail",
		Short: "报货单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if declareNo == "" {
				return fmt.Errorf("--declare-no 必填")
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetDeclareOrderDetail(cmd.Context(), declareNo, shopCode)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"报货单号", "状态", "金额", "实际金额", "审核金额", "品项数", "总数", "创建时间"}, [][]string{{
				result.DeclareNo, strconv.Itoa(result.OrderStatus), fmt.Sprintf("%.2f", result.Amount), fmt.Sprintf("%.2f", result.ActualAmount), fmt.Sprintf("%.2f", result.AuditAmount), strconv.Itoa(result.ProductCateNum), strconv.Itoa(result.ProductNum), result.CreatedAt,
			}})
		},
	}
	cmd.Flags().StringVar(&declareNo, "declare-no", "", "报货单号")
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	return cmd
}

func newCmdRequireUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "require-update",
		Short: "订货单取消/完成",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"订货单取消/完成", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CompleteRequireOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdRequireDeliver(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "require-deliver",
		Short: "订货单发货",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"订货单发货", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.DeliverRequireOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdReturnCancel(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "return-cancel",
		Short: "退货单取消/驳回",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"退货单取消/驳回", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CancelReturnOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdReturnExamine(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "return-examine",
		Short: "退货单审核",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"退货单审核", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ExamineReturnOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdReturnReceipt(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "return-receipt",
		Short: "退货单收货/拒收",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"退货单收货/拒收", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ReceiptReturnOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdDeliveryArrive(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "delivery-arrive",
		Short: "配送单确认送达",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"配送单确认送达", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ConfirmDeliveryArrive(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdRequireList(f *cmdutil.Factory) *cobra.Command {
	var requireNo, startTime, endTime string
	var pageNo, pageSize int
	var storeIDs, distributeIDList, requireProIDList, statusList []int64
	cmd := &cobra.Command{
		Use:   "require-list",
		Short: "查询订货单列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize}
			if requireNo != "" {
				params["requireNo"] = requireNo
			}
			if startTime != "" {
				params["startTime"] = startTime
			}
			if endTime != "" {
				params["endTime"] = endTime
			}
			if len(storeIDs) > 0 {
				params["storeIds"] = storeIDs
			}
			if len(distributeIDList) > 0 {
				params["distributeIdList"] = distributeIDList
			}
			if len(requireProIDList) > 0 {
				params["requireProIdList"] = requireProIDList
			}
			if len(statusList) > 0 {
				params["statusList"] = statusList
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListRequireOrders(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.RequireNo, item.DeclareNo, item.WarehouseName, item.StoreName,
					strconv.Itoa(item.OrderStatus), fmt.Sprintf("%.2f", item.Amount), fmt.Sprintf("%.2f", item.ProductNum), item.OrderAt,
				})
			}
			if err := writeRows(f, result, []string{"订货单号", "来源单号", "配送中心", "门店", "状态", "金额", "数量", "订货时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&requireNo, "require-no", "", "订货单号")
	cmd.Flags().StringVar(&startTime, "start-time", "", "创建开始时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().StringVar(&endTime, "end-time", "", "创建结束时间 yyyy-MM-dd HH:mm:ss")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	cmd.Flags().Int64SliceVar(&storeIDs, "store-ids", nil, "门店ID列表")
	cmd.Flags().Int64SliceVar(&distributeIDList, "distribute-ids", nil, "配送中心ID列表")
	cmd.Flags().Int64SliceVar(&requireProIDList, "require-pro-ids", nil, "订货单品项ID列表")
	cmd.Flags().Int64SliceVar(&statusList, "status-list", nil, "状态列表")
	return cmd
}

func newCmdRequireCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "require-create",
		Short: "订货单创建",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"订货单创建", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.CreateRequireOrder(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交订货单创建请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdRequireDetailUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "require-detail-update",
		Short: "订货单品项批量修改",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"订货单品项批量修改", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.UpdateRequireOrderDetails(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交订货单品项批量修改请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdTransferList(f *cmdutil.Factory) *cobra.Command {
	var createdStartAt, createdEndAt, creator, inWareNo, outWareNo, transferStartAt, transferEndAt string
	var pageNo, pageSize, transferType int
	var productCodeList, transferNoList []string
	var statusList []int64
	cmd := &cobra.Command{
		Use:   "transfer-list",
		Short: "查询调拨单列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize}
			if createdStartAt != "" {
				params["createdStartAt"] = createdStartAt
			}
			if createdEndAt != "" {
				params["createdEndAt"] = createdEndAt
			}
			if creator != "" {
				params["creator"] = creator
			}
			if inWareNo != "" {
				params["inWareNo"] = inWareNo
			}
			if outWareNo != "" {
				params["outWareNo"] = outWareNo
			}
			if transferStartAt != "" {
				params["transferStartAt"] = transferStartAt
			}
			if transferEndAt != "" {
				params["transferEndAt"] = transferEndAt
			}
			if transferType != 0 {
				params["transferType"] = transferType
			}
			if len(productCodeList) > 0 {
				params["productCodeList"] = productCodeList
			}
			if len(transferNoList) > 0 {
				params["transferNoList"] = transferNoList
			}
			if len(statusList) > 0 {
				params["statusList"] = statusList
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListTransferOrders(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.TransferNo, item.OutWareName, item.InWareName, strconv.Itoa(item.Status),
					fmt.Sprintf("%.2f", item.PriceAmount), fmt.Sprintf("%.2f", item.ProductNum), item.TransferAt,
				})
			}
			if err := writeRows(f, result, []string{"调拨单号", "调出仓", "调入仓", "状态", "金额", "数量", "调拨日期"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&createdStartAt, "created-start-at", "", "创建开始时间")
	cmd.Flags().StringVar(&createdEndAt, "created-end-at", "", "创建结束时间")
	cmd.Flags().StringVar(&creator, "creator", "", "创建人")
	cmd.Flags().StringVar(&inWareNo, "in-ware-no", "", "调入仓库编码")
	cmd.Flags().StringVar(&outWareNo, "out-ware-no", "", "调出仓库编码")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	cmd.Flags().StringSliceVar(&productCodeList, "product-codes", nil, "品项编码列表")
	cmd.Flags().Int64SliceVar(&statusList, "status-list", nil, "状态列表")
	cmd.Flags().StringVar(&transferStartAt, "transfer-start-at", "", "调拨开始时间")
	cmd.Flags().StringVar(&transferEndAt, "transfer-end-at", "", "调拨结束时间")
	cmd.Flags().StringSliceVar(&transferNoList, "transfer-nos", nil, "调拨单号列表")
	cmd.Flags().IntVar(&transferType, "transfer-type", 0, "调拨类型")
	return cmd
}

func newCmdTransferDetail(f *cmdutil.Factory) *cobra.Command {
	var transferNo string
	cmd := &cobra.Command{
		Use:   "transfer-detail",
		Short: "查询调拨单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if transferNo == "" {
				return fmt.Errorf("--transfer-no 必填")
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetTransferOrderDetail(cmd.Context(), transferNo)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"调拨单号", "调出仓", "调入仓", "状态", "金额", "数量", "创建人", "创建时间"}, [][]string{{
				result.TransferNo, result.OutWareName, result.InWareName, strconv.Itoa(result.Status), fmt.Sprintf("%.2f", result.PriceAmount), fmt.Sprintf("%.2f", result.ProductNum), result.Creator, result.CreatedAt,
			}})
		},
	}
	cmd.Flags().StringVar(&transferNo, "transfer-no", "", "调拨单号")
	return cmd
}

func newCmdReturnList(f *cmdutil.Factory) *cobra.Command {
	var requireNo, returnNo, startDate, endDate string
	var pageNo, pageSize, refundStatus, returnStatus int
	var storeIDs, distributeIDList, statusList []int64
	cmd := &cobra.Command{
		Use:   "return-list",
		Short: "查询退货单列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize}
			if requireNo != "" {
				params["requireNo"] = requireNo
			}
			if returnNo != "" {
				params["returnNo"] = returnNo
			}
			if startDate != "" {
				params["startDate"] = startDate
			}
			if endDate != "" {
				params["endDate"] = endDate
			}
			if refundStatus != 0 {
				params["refundStatus"] = refundStatus
			}
			if returnStatus != 0 {
				params["returnStatus"] = returnStatus
			}
			if len(storeIDs) > 0 {
				params["storeIds"] = storeIDs
			}
			if len(distributeIDList) > 0 {
				params["distributeIdList"] = distributeIDList
			}
			if len(statusList) > 0 {
				params["statusList"] = statusList
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListReturnOrders(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.ReturnNo, item.RequireNo, item.StoreName, item.WarehouseName,
					strconv.Itoa(item.ReturnStatus), strconv.Itoa(item.RefundStatus),
					fmt.Sprintf("%.2f", item.ReturnAmount), item.CreatedAt,
				})
			}
			if err := writeRows(f, result, []string{"退货单号", "订货单号", "门店", "配送中心", "单据状态", "退款状态", "退货金额", "创建时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&requireNo, "require-no", "", "来源订货单号")
	cmd.Flags().StringVar(&returnNo, "return-no", "", "退货单号")
	cmd.Flags().StringVar(&startDate, "start-date", "", "退货开始日期")
	cmd.Flags().StringVar(&endDate, "end-date", "", "退货结束日期")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	cmd.Flags().IntVar(&refundStatus, "refund-status", 0, "退款状态")
	cmd.Flags().IntVar(&returnStatus, "return-status", 0, "单据状态")
	cmd.Flags().Int64SliceVar(&storeIDs, "store-ids", nil, "门店ID列表")
	cmd.Flags().Int64SliceVar(&distributeIDList, "distribute-ids", nil, "配送中心ID列表")
	cmd.Flags().Int64SliceVar(&statusList, "status-list", nil, "状态列表")
	return cmd
}

func newCmdInboundList(f *cmdutil.Factory) *cobra.Command {
	var createdStartAt, createdEndAt string
	var pageNo, pageSize int
	cmd := &cobra.Command{
		Use:   "inbound-list",
		Short: "批量查询入库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createdStartAt == "" || createdEndAt == "" {
				return fmt.Errorf("--created-start-at、--created-end-at 必填")
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListInboundOrders(cmd.Context(), map[string]interface{}{
				"createdStartAt": createdStartAt,
				"createdEndAt":   createdEndAt,
				"pageNo":         pageNo,
				"pageSize":       pageSize,
			})
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.InboundNo, item.BizNo, item.WarehouseName, item.SupplierName,
					strconv.Itoa(item.Status), fmt.Sprintf("%.2f", item.Amount), item.InboundAt,
				})
			}
			return writeRows(f, result, []string{"入库单号", "业务单据", "仓库", "供应商", "状态", "金额", "入库时间"}, rows)
		},
	}
	cmd.Flags().StringVar(&createdStartAt, "created-start-at", "", "创建开始时间")
	cmd.Flags().StringVar(&createdEndAt, "created-end-at", "", "创建结束时间")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	return cmd
}

func newCmdInboundCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "inbound-create",
		Short: "创建入库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建入库单", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CreateInboundOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdInboundUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "inbound-update",
		Short: "入库单物品数量更新",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"入库单物品数量更新", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.UpdateInboundOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdInboundFinish(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "inbound-finish",
		Short: "关闭入库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"关闭入库单", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.FinishInboundOrders(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交关闭入库单请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdOutboundList(f *cmdutil.Factory) *cobra.Command {
	var createdStartAt, createdEndAt string
	var pageNo, pageSize int
	cmd := &cobra.Command{
		Use:   "outbound-list",
		Short: "批量查询出库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createdStartAt == "" || createdEndAt == "" {
				return fmt.Errorf("--created-start-at、--created-end-at 必填")
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListOutboundOrders(cmd.Context(), map[string]interface{}{
				"createdStartAt": createdStartAt,
				"createdEndAt":   createdEndAt,
				"pageNo":         pageNo,
				"pageSize":       pageSize,
			})
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.OutboundNo, item.BizNo, item.WarehouseName, item.ReceiptName,
					strconv.Itoa(item.Status), fmt.Sprintf("%.2f", item.Amount), item.OutboundAt,
				})
			}
			return writeRows(f, result, []string{"出库单号", "业务单据", "仓库", "收货方", "状态", "金额", "出库时间"}, rows)
		},
	}
	cmd.Flags().StringVar(&createdStartAt, "created-start-at", "", "创建开始时间")
	cmd.Flags().StringVar(&createdEndAt, "created-end-at", "", "创建结束时间")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	return cmd
}

func newCmdOutboundCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "outbound-create",
		Short: "创建出库单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建出库单", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CreateOutboundOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdOutboundUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "outbound-update",
		Short: "出库单物品数量更新",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"出库单物品数量更新", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.UpdateOutboundOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdStockOccupy(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "occupy",
		Short: "物品库存锁定",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"物品库存锁定", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.OccupyProductStock(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"是否成功"}, [][]string{{strconv.FormatBool(bool(*result))}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdStockRelease(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "release",
		Short: "物品库存释放",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"物品库存释放", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ReleaseProductStock(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"是否成功"}, [][]string{{strconv.FormatBool(bool(*result))}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdInventoryAdjust(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "inventory-adjust",
		Short: "创建盘点调整单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建盘点调整单", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.CreateInventoryAdjust(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交盘点调整单创建请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdTransferAudit(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "transfer-audit",
		Short: "调拨申请单审核",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"调拨申请单审核", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.AuditTransferOrder(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交调拨申请单审核请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdProductCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "product-create",
		Short: "新增物品信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"新增物品信息", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CreateProduct(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCategoryCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "category-create",
		Short: "创建物品分类",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建物品分类", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CreateCategory(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdUnitCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "unit-create",
		Short: "创建物品单位",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建物品单位", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CreateUnit(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"单位编码"}, [][]string{{result.UnitCode}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdProductUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "product-update",
		Short: "修改物品信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"修改物品信息", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.UpdateProduct(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"业务ID"}, [][]string{{strconv.FormatInt(int64(*result), 10)}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdSupplierCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "supplier-create",
		Short: "新增供应商",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"新增供应商", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.CreateSupplier(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交供应商创建请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdProductDistribute(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "product-distribute",
		Short: "品项下发",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"品项下发", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.BatchGroupProduct(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交品项下发请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdSupplierUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "supplier-update",
		Short: "修改供应商基础资料",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"修改供应商基础资料", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.UpdateSupplier(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交供应商更新请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdMachiningCardBatchCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "machining-card-batch-create",
		Short: "批量创建半成品成本卡",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"批量创建半成品成本卡", fromJSON}})
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			if err := api.BatchCreateMachiningCards(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交半成品成本卡批量创建请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdStoreLedger(f *cmdutil.Factory) *cobra.Command {
	var startDate, endDate, storeID string
	var pageNo, pageSize int
	var categoryIDs, productCodes, subjectCodes []string
	cmd := &cobra.Command{
		Use:   "store-ledger",
		Short: "门店库存日志查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if startDate == "" || endDate == "" {
				return fmt.Errorf("--start-date、--end-date 必填")
			}
			params := map[string]interface{}{
				"start_date": startDate,
				"end_date":   endDate,
				"pageNo":     pageNo,
				"pageSize":   pageSize,
			}
			if storeID != "" {
				params["store_id"] = storeID
			}
			if len(categoryIDs) > 0 {
				params["category_id"] = categoryIDs
			}
			if len(productCodes) > 0 {
				params["product_code"] = productCodes
			}
			if len(subjectCodes) > 0 {
				params["financial_subject_code"] = subjectCodes
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetStoreInventorySummary(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.ResultList))
			for _, item := range result.ResultList {
				rows = append(rows, []string{
					item.TheDate, item.StoreName, item.ProductCode, item.ProductName,
					fmt.Sprintf("%.2f", item.BeforeNum), fmt.Sprintf("%.2f", item.AfterNum), item.UnitName,
				})
			}
			if err := writeRows(f, result, []string{"日期", "门店", "品项编码", "品项名称", "期初数量", "期末数量", "单位"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.TotalCount)
			return nil
		},
	}
	cmd.Flags().StringVar(&startDate, "start-date", "", "开始日期")
	cmd.Flags().StringVar(&endDate, "end-date", "", "结束日期")
	cmd.Flags().StringVar(&storeID, "store-id", "", "门店ID")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	cmd.Flags().StringSliceVar(&categoryIDs, "category-ids", nil, "品项分类ID列表")
	cmd.Flags().StringSliceVar(&productCodes, "product-codes", nil, "品项编码列表")
	cmd.Flags().StringSliceVar(&subjectCodes, "subject-codes", nil, "统计科目编码列表")
	return cmd
}

func newCmdRealtimeList(f *cmdutil.Factory) *cobra.Command {
	var pageNo, pageSize, isEmpty int
	var warehouseNos, productCodes []string
	cmd := &cobra.Command{
		Use:   "realtime-list",
		Short: "查询实时库存列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize}
			if len(warehouseNos) > 0 {
				params["warehouseNoList"] = warehouseNos
			}
			if len(productCodes) > 0 {
				params["productCodeList"] = productCodes
			}
			if isEmpty != 0 {
				params["isEmpty"] = isEmpty
			}
			api, err := newInventoryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListWarehouseProducts(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.Data))
			for _, item := range result.Data {
				rows = append(rows, []string{
					item.WarehouseNo, item.WarehouseName, item.ProductCode, item.ProductName,
					fmt.Sprintf("%.2f", item.Quantity), fmt.Sprintf("%.2f", item.AvailableQuantity), fmt.Sprintf("%.2f", item.CurrentAmount),
				})
			}
			if err := writeRows(f, result, []string{"仓库编码", "仓库名称", "品项编码", "品项名称", "库存量", "可用库存量", "现存金额"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 10, "每页条数")
	cmd.Flags().StringSliceVar(&warehouseNos, "warehouse-nos", nil, "仓库编码列表")
	cmd.Flags().StringSliceVar(&productCodes, "product-codes", nil, "品项编码列表")
	cmd.Flags().IntVar(&isEmpty, "is-empty", 0, "是否包含 0 库存，1=包含")
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
