package queue

import (
	"fmt"
	"strconv"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdQueue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "排队服务",
		Long:  "查询门店排队进度、订单排队进度和门店叫号列表",
	}
	cmd.AddCommand(newCmdShopProgress(f))
	cmd.AddCommand(newCmdOrderProgress(f))
	cmd.AddCommand(newCmdShopQueueNos(f))
	return cmd
}

func newQueueAPI(f *cmdutil.Factory) (*client.QueueAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewQueueAPI(apiClient), nil
}

func newCmdShopProgress(f *cmdutil.Factory) *cobra.Command {
	var shopType int
	var shopIDs []int64
	cmd := &cobra.Command{
		Use:   "shop-progress",
		Short: "查询门店排队进度",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopType == 0 || len(shopIDs) == 0 {
				return fmt.Errorf("--shop-type、--shop-ids 必填")
			}
			api, err := newQueueAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryShopQueueProgress(cmd.Context(), shopIDs, shopType)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.ShopID, strconv.Itoa(item.OrderNum), strconv.Itoa(item.CupTotal), strconv.Itoa(item.MakeTime),
				})
			}
			return writeRows(f, result, []string{"门店ID", "排队订单数", "排队杯数", "等待时间(秒)"}, rows)
		},
	}
	cmd.Flags().IntVar(&shopType, "shop-type", 0, "门店类型 1=企迈 4=美团外卖 5=饿了么 16=京东到家 25=抖音团购")
	cmd.Flags().Int64SliceVar(&shopIDs, "shop-ids", nil, "门店ID列表")
	return cmd
}

func newCmdOrderProgress(f *cmdutil.Factory) *cobra.Command {
	var orderNo, sourceNo string
	cmd := &cobra.Command{
		Use:   "order-progress",
		Short: "查询订单排队进度",
		RunE: func(cmd *cobra.Command, args []string) error {
			if orderNo == "" && sourceNo == "" {
				return fmt.Errorf("--order-no 和 --source-no 至少一个必填")
			}
			api, err := newQueueAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryOrderQueueProgress(cmd.Context(), orderNo, sourceNo)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"排队订单数", "排队杯数", "等待时间(秒)", "订单状态"}, [][]string{{
				strconv.Itoa(result.OrderNum), strconv.Itoa(result.CupTotal), strconv.Itoa(result.MakeTime), strconv.Itoa(result.OrderStatus),
			}})
		},
	}
	cmd.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	cmd.Flags().StringVar(&sourceNo, "source-no", "", "渠道订单号")
	return cmd
}

func newCmdShopQueueNos(f *cmdutil.Factory) *cobra.Command {
	var page, size int
	var shopCode string
	var statusList []int
	cmd := &cobra.Command{
		Use:   "shop-queue-nos",
		Short: "查询门店排队叫号列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopCode == "" || page == 0 || size == 0 {
				return fmt.Errorf("--shop-code、--page、--size 必填")
			}
			api, err := newQueueAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryShopQueueNoList(cmd.Context(), shopCode, page, size, statusList)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.QueueNo, item.OrderNo, strconv.Itoa(item.QueueNoStatus), strconv.Itoa(item.QueueNoOrderSource), strconv.FormatInt(item.UserID, 10),
				})
			}
			if err := writeRows(f, result, []string{"取单号", "订单号", "取单号状态", "订单来源", "用户ID"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	cmd.Flags().IntVar(&page, "page", 1, "页码")
	cmd.Flags().IntVar(&size, "size", 10, "每页条数，1-20")
	cmd.Flags().IntSliceVar(&statusList, "status-list", nil, "取单号状态列表")
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
