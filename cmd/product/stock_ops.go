package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdSoldOut(f *cmdutil.Factory) *cobra.Command {
	return newExternalStockCmd(f, "sold-out", "门店商品批量售罄", true)
}

func newCmdFillFull(f *cmdutil.Factory) *cobra.Command {
	return newExternalStockCmd(f, "fill-full", "门店商品批量置满", false)
}

func newExternalStockCmd(f *cmdutil.Factory, use, short string, empty bool) *cobra.Command {
	var tradeMarks []string
	var specCodes []string
	var saleChannels []int
	var saleTypes []int
	var isAll bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(tradeMarks) == 0 {
				return fmt.Errorf("--trade-marks 必填")
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			params := client.ExternalStockParams{
				IsAll:           isAll,
				SaleChannelList: saleChannels,
				SaleTypeList:    saleTypes,
				MultiMark:       profile.ShopCode,
				TradeMarkList:   tradeMarks,
				SpecCodeList:    specCodes,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将执行 %s: %v\n", use, tradeMarks)
				return nil
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			if empty {
				err = api.ExternalEmpty(cmd.Context(), params)
			} else {
				err = api.ExternalFull(cmd.Context(), params)
			}
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已执行 %s，共 %d 个商品\n", use, len(tradeMarks))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tradeMarks, "trade-marks", nil, "商品标识列表")
	cmd.Flags().StringSliceVar(&specCodes, "spec-codes", nil, "规格码列表")
	cmd.Flags().IntSliceVar(&saleChannels, "sale-channels", nil, "售卖渠道列表")
	cmd.Flags().IntSliceVar(&saleTypes, "sale-types", nil, "售卖类型列表")
	cmd.Flags().BoolVar(&isAll, "is-all", false, "是否全部")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdPracticeStatus(f *cmdutil.Factory) *cobra.Command {
	var practiceValues []string
	var status int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "practice-status",
		Short: "门店商品做法启用停用",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(practiceValues) == 0 {
				return fmt.Errorf("--practice-values 必填")
			}
			if status != 0 && status != 1 {
				return fmt.Errorf("--status 必须为 0(停用) 或 1(启用)")
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将更新做法状态: %v -> %d\n", practiceValues, status)
				return nil
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			err = api.PracticeOnOff(cmd.Context(), client.PracticeStatusParams{
				PracticeValues: practiceValues,
				ShopCode:       profile.ShopCode,
				Status:         status,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已更新做法状态，共 %d 个做法值\n", len(practiceValues))
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&practiceValues, "practice-values", nil, "做法值名称列表")
	cmd.Flags().IntVar(&status, "status", 0, "状态，0=停用 1=启用")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdDeleteTask(f *cmdutil.Factory) *cobra.Command {
	var saleChannel int
	var saleType int
	var specCodes []string
	var tradeMarks []string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "delete-task",
		Short: "提交商品删除任务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if saleChannel == 0 {
				return fmt.Errorf("--sale-channel 必填")
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			profile := cfg.Profiles[f.EffectiveProfile()]
			if profile == nil || profile.ShopCode == "" {
				return fmt.Errorf("未配置门店编码，运行 'qmai auth login' 设置 shopCode")
			}
			if len(specCodes) == 0 && len(tradeMarks) == 0 {
				return fmt.Errorf("--spec-codes 或 --trade-marks 至少一个必填")
			}
			params := client.DeleteTaskParams{
				SaleChannel:   saleChannel,
				SaleType:      saleType,
				SpecCodeList:  specCodes,
				TradeMarkList: tradeMarks,
				StoreCode:     profile.ShopCode,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将提交商品删除任务: tradeMarks=%v specCodes=%v\n", tradeMarks, specCodes)
				return nil
			}
			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}
			api := client.NewProductAPI(apiClient)
			resp, err := api.SubmitDeleteTask(cmd.Context(), params)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交删除任务: taskId=%s\n", resp.TaskID)
			return nil
		},
	}

	cmd.Flags().IntVar(&saleChannel, "sale-channel", 0, "售卖渠道 1-美团 2-饿了么 3-小程序 4-POS")
	cmd.Flags().IntVar(&saleType, "sale-type", 0, "售卖类型")
	cmd.Flags().StringSliceVar(&specCodes, "spec-codes", nil, "规格码列表")
	cmd.Flags().StringSliceVar(&tradeMarks, "trade-marks", nil, "商品标识列表")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}
