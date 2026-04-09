package delivery

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

func NewCmdDelivery(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delivery",
		Short: "聚合配送",
		Long:  "聚合配送订单创建、取消、详情、骑手位置与状态同步",
	}
	cmd.AddCommand(newCmdOrder(f))
	cmd.AddCommand(newCmdStatus(f))
	return cmd
}

func newDeliveryAPI(f *cmdutil.Factory) (*client.DeliveryAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewDeliveryAPI(apiClient), nil
}

func newCmdOrder(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "order", Short: "配送订单管理"}
	cmd.AddCommand(newCmdCreate(f))
	cmd.AddCommand(newCmdCancel(f))
	cmd.AddCommand(newCmdDetail(f))
	cmd.AddCommand(newCmdRiderLocation(f))
	cmd.AddCommand(newCmdCancelAll(f))
	cmd.AddCommand(newCmdCreateAndUpdate(f))
	return cmd
}

func newCmdStatus(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "status", Short: "配送状态同步"}
	cmd.AddCommand(newCmdUpdate(f))
	return cmd
}

func newCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建配送订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建配送订单", fromJSON}})
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			if err := api.CreateDeliveryOrder(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交配送订单创建请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCancel(f *cmdutil.Factory) *cobra.Command {
	var cancelCode, cancelReason, multiMark, originOrderNo string
	var orderSource int
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "取消配送订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if multiMark == "" || orderSource == 0 || originOrderNo == "" {
				return fmt.Errorf("--multi-mark、--order-source、--origin-order-no 必填")
			}
			params := map[string]interface{}{
				"multiMark":     multiMark,
				"orderSource":   orderSource,
				"originOrderNo": originOrderNo,
			}
			if cancelCode != "" {
				params["cancelCode"] = cancelCode
			}
			if cancelReason != "" {
				params["cancelReason"] = cancelReason
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "门店编码", "订单来源", "订单号"}, [][]string{{"取消配送订单", multiMark, strconv.Itoa(orderSource), originOrderNo}})
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			if err := api.CancelDeliveryOrder(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交取消请求: %s\n", originOrderNo)
			return nil
		},
	}
	cmd.Flags().StringVar(&cancelCode, "cancel-code", "", "取消代号")
	cmd.Flags().StringVar(&cancelReason, "cancel-reason", "", "取消原因")
	cmd.Flags().StringVar(&multiMark, "multi-mark", "", "门店编码")
	cmd.Flags().IntVar(&orderSource, "order-source", 0, "订单来源")
	cmd.Flags().StringVar(&originOrderNo, "origin-order-no", "", "订单号")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdDetail(f *cmdutil.Factory) *cobra.Command {
	var orderFlag, originOrderNo string
	var orderSource int
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "查询配送单详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderSource == 0 || originOrderNo == "" {
				return fmt.Errorf("--order-source、--origin-order-no 必填")
			}
			params := map[string]interface{}{
				"orderSource":   orderSource,
				"originOrderNo": originOrderNo,
			}
			if orderFlag != "" {
				params["orderFlag"] = orderFlag
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetDeliveryOrderInfo(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := [][]string{{
				result.BizID, result.OrderStatus, result.DriverName, result.DriverPhone,
				result.DeliveryCost, result.PushTime, result.CancelDesc, result.ErrorDesc,
			}}
			if err := writeRows(f, result, []string{"业务标识", "配送状态", "骑手姓名", "骑手电话", "配送费", "推送时间", "取消原因", "错误描述"}, rows); err != nil {
				return err
			}
			if len(result.DeliveryOrderLogs) > 0 {
				logRows := make([][]string, 0, len(result.DeliveryOrderLogs))
				for _, item := range result.DeliveryOrderLogs {
					logRows = append(logRows, []string{item.LogTime, item.LogStr})
				}
				fmt.Fprintln(f.IOStreams.Out, "")
				return writeRows(f, result.DeliveryOrderLogs, []string{"日志时间", "日志描述"}, logRows)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&orderFlag, "order-flag", "", "呼叫标识 0/1/2/9")
	cmd.Flags().IntVar(&orderSource, "order-source", 0, "订单来源")
	cmd.Flags().StringVar(&originOrderNo, "origin-order-no", "", "订单号")
	return cmd
}

func newCmdRiderLocation(f *cmdutil.Factory) *cobra.Command {
	var orderFlag, originOrderNo string
	var orderSource int
	cmd := &cobra.Command{
		Use:   "rider-location",
		Short: "查询骑手当前位置",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderSource == 0 || originOrderNo == "" {
				return fmt.Errorf("--order-source、--origin-order-no 必填")
			}
			params := map[string]interface{}{
				"orderSource":   orderSource,
				"originOrderNo": originOrderNo,
			}
			if orderFlag != "" {
				params["orderFlag"] = orderFlag
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetRiderLocation(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"骑手姓名", "骑手电话", "三方配送单号", "运单状态", "纬度", "经度", "三方平台编码"}, [][]string{{
				result.CarrierName, result.CarrierPhone, result.DeliveryOrderNo, result.DeliveryStatus, result.Latitude, result.Longitude, result.UploadDeliveryType,
			}})
		},
	}
	cmd.Flags().StringVar(&orderFlag, "order-flag", "", "呼叫标识 0/1/2/9")
	cmd.Flags().IntVar(&orderSource, "order-source", 0, "订单来源")
	cmd.Flags().StringVar(&originOrderNo, "origin-order-no", "", "订单号")
	return cmd
}

func newCmdCancelAll(f *cmdutil.Factory) *cobra.Command {
	var cancelCode, cancelReason, multiMark, originOrderNo string
	var handCancel bool
	var orderSource int
	var shopID int64
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "cancel-all",
		Short: "取消全部配送单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderSource == 0 || originOrderNo == "" {
				return fmt.Errorf("--order-source、--origin-order-no 必填")
			}
			params := map[string]interface{}{
				"orderSource":   orderSource,
				"originOrderNo": originOrderNo,
			}
			if cancelCode != "" {
				params["cancelCode"] = cancelCode
			}
			if cancelReason != "" {
				params["cancelReason"] = cancelReason
			}
			if handCancel {
				params["handCancel"] = handCancel
			}
			if multiMark != "" {
				params["multiMark"] = multiMark
			}
			if shopID != 0 {
				params["shopId"] = shopID
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "订单来源", "订单号"}, [][]string{{"取消全部配送单", strconv.Itoa(orderSource), originOrderNo}})
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CancelAllDeliveryOrder(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.DeliveryID, item.OriginOrderNo, strconv.FormatBool(item.Result)})
			}
			return writeRows(f, result, []string{"聚合运单号", "订单号", "结果"}, rows)
		},
	}
	cmd.Flags().StringVar(&cancelCode, "cancel-code", "", "取消代号")
	cmd.Flags().StringVar(&cancelReason, "cancel-reason", "", "取消原因")
	cmd.Flags().BoolVar(&handCancel, "hand-cancel", false, "手动取消")
	cmd.Flags().StringVar(&multiMark, "multi-mark", "", "门店编码")
	cmd.Flags().IntVar(&orderSource, "order-source", 0, "订单来源")
	cmd.Flags().StringVar(&originOrderNo, "origin-order-no", "", "订单号")
	cmd.Flags().Int64Var(&shopID, "shop-id", 0, "门店ID")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var cancelDesc, driverName, driverPhone, originOrderNo string
	var orderStatus, platformType int
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "更新订单配送状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderStatus == 0 || originOrderNo == "" {
				return fmt.Errorf("--order-status、--origin-order-no 必填")
			}
			params := map[string]interface{}{
				"orderStatus":   orderStatus,
				"originOrderNo": originOrderNo,
			}
			if cancelDesc != "" {
				params["cancelDesc"] = cancelDesc
			}
			if driverName != "" {
				params["driverName"] = driverName
			}
			if driverPhone != "" {
				params["driverPhone"] = driverPhone
			}
			if platformType != 0 {
				params["platformType"] = platformType
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "订单号", "状态"}, [][]string{{"更新订单配送状态", originOrderNo, strconv.Itoa(orderStatus)}})
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			if err := api.UpdateSelfOrderStatus(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交配送状态更新: %s -> %d\n", originOrderNo, orderStatus)
			return nil
		},
	}
	cmd.Flags().StringVar(&cancelDesc, "cancel-desc", "", "取消原因")
	cmd.Flags().StringVar(&driverName, "driver-name", "", "配送员姓名")
	cmd.Flags().StringVar(&driverPhone, "driver-phone", "", "配送员手机号")
	cmd.Flags().IntVar(&orderStatus, "order-status", 0, "配送状态 1-8")
	cmd.Flags().StringVar(&originOrderNo, "origin-order-no", "", "订单号")
	cmd.Flags().IntVar(&platformType, "platform-type", 0, "配送平台")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCreateAndUpdate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "create-and-update",
		Short: "创建配送单并更新配送状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"创建配送单并更新配送状态", fromJSON}})
			}
			api, err := newDeliveryAPI(f)
			if err != nil {
				return err
			}
			if err := api.CreateAndUpdateDeliveryStatus(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintln(f.IOStreams.Out, "✓ 已提交配送单创建与状态更新请求")
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
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
