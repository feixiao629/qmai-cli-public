package marketing

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

func NewCmdMarketing(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketing",
		Short: "营销服务",
		Long:  "券管理、活动管理、礼品卡管理、促销查询与结算页算价",
	}
	cmd.AddCommand(newCmdCoupon(f))
	cmd.AddCommand(newCmdCampaign(f))
	cmd.AddCommand(newCmdGiftCard(f))
	cmd.AddCommand(newCmdPricing(f))
	return cmd
}

func newMarketingAPI(f *cmdutil.Factory) (*client.MarketingAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewMarketingAPI(apiClient), nil
}

func newCmdCoupon(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "coupon", Short: "券管理"}
	cmd.AddCommand(newCmdCouponStatus(f))
	cmd.AddCommand(newCmdCouponDetail(f))
	cmd.AddCommand(newCmdCouponTemplate(f))
	cmd.AddCommand(newCmdCouponTemplateEnable(f))
	cmd.AddCommand(newCmdCouponChoose(f))
	cmd.AddCommand(newCmdAnonymousCoupon(f))
	cmd.AddCommand(newCmdCouponTemplateBatch(f))
	cmd.AddCommand(newCmdCouponTemplateByThirdCode(f))
	cmd.AddCommand(newCmdGrantByActivity(f))
	cmd.AddCommand(newCmdGrantByActivityAsync(f))
	cmd.AddCommand(newCmdGrantByTemplate(f))
	return cmd
}

func newCmdCampaign(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "campaign", Short: "活动管理"}
	cmd.AddCommand(newCmdCampaignList(f))
	cmd.AddCommand(newCmdExchangeDispatch(f))
	cmd.AddCommand(newCmdExchangeDisable(f))
	cmd.AddCommand(newCmdExchangeStatus(f))
	cmd.AddCommand(newCmdTaskRecords(f))
	cmd.AddCommand(newCmdTaskClaim(f))
	cmd.AddCommand(newCmdRecycleCoupons(f))
	cmd.AddCommand(newCmdRecoverGrant(f))
	return cmd
}

func newCmdGiftCard(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "gift-card", Short: "卡管理"}
	cmd.AddCommand(newCmdGiftCardConsume(f))
	cmd.AddCommand(newCmdGiftCardConsumeBatch(f))
	cmd.AddCommand(newCmdGiftCardReverse(f))
	cmd.AddCommand(newCmdGiftCardReverseBatch(f))
	cmd.AddCommand(newCmdGiftCardPartReverse(f))
	cmd.AddCommand(newCmdGiftCardIssue(f))
	cmd.AddCommand(newCmdGiftCardReportLoss(f))
	cmd.AddCommand(newCmdGiftCardRelieveLoss(f))
	cmd.AddCommand(newCmdGiftCardRecycle(f))
	cmd.AddCommand(newCmdGiftCardExchange(f))
	cmd.AddCommand(newCmdGiftCardGrantByTemplate(f))
	cmd.AddCommand(newCmdGiftCardInfo(f))
	cmd.AddCommand(newCmdGiftCardFlow(f))
	cmd.AddCommand(newCmdGiftCardTemplate(f))
	cmd.AddCommand(newCmdGiftCardList(f))
	cmd.AddCommand(newCmdGiftCardTemplateBatch(f))
	return cmd
}

func newCmdPricing(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{Use: "pricing", Short: "营销管理与计算"}
	cmd.AddCommand(newCmdPromotionActivities(f))
	cmd.AddCommand(newCmdConfirm(f))
	return cmd
}

func newCmdCouponStatus(f *cmdutil.Factory) *cobra.Command {
	var id string
	var codeType int
	cmd := &cobra.Command{
		Use:   "status",
		Short: "查询券状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("--id 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCouponStatus(cmd.Context(), id, codeType)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"券标识", "状态", "订单号", "门店名称", "门店编码", "使用时间"}, [][]string{{
				id, strconv.Itoa(result.Status), result.UseOrderCode, result.UseStoreName, result.UseStoreNo, result.UseTime,
			}})
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "券标识")
	cmd.Flags().IntVar(&codeType, "type", 0, "类型，1=券码 2=兑换码")
	return cmd
}

func newCmdCouponDetail(f *cmdutil.Factory) *cobra.Command {
	var code string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "查询券详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if code == "" {
				return fmt.Errorf("--user-coupon-code 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCouponDetail(cmd.Context(), code)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"券码", "模板ID", "模板名称", "状态", "归属会员", "门店", "使用时间"}, [][]string{{
				strconv.FormatInt(result.ID, 10),
				strconv.FormatInt(result.TemplateID, 10),
				result.TemplateName,
				strconv.Itoa(result.Status),
				result.CustomerID,
				result.UseStoreName,
				result.UseTime,
			}})
		},
	}
	cmd.Flags().StringVar(&code, "user-coupon-code", "", "券码")
	return cmd
}

func newCmdCouponTemplate(f *cmdutil.Factory) *cobra.Command {
	var id int64
	cmd := &cobra.Command{
		Use:   "template",
		Short: "查询券模板详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return fmt.Errorf("--id 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCouponTemplateDetail(cmd.Context(), id)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"模板ID", "模板名称", "品牌ID", "状态", "券类型", "三方编码"}, [][]string{{
				strconv.FormatInt(result.ID, 10), result.Name, strconv.FormatInt(result.SellerID, 10), strconv.Itoa(result.Status), strconv.Itoa(result.Type), result.ThirdBizCode,
			}})
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "券模板ID")
	return cmd
}

func newCmdCouponTemplateEnable(f *cmdutil.Factory) *cobra.Command {
	var id int64
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "template-enable",
		Short: "启用券模板",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return fmt.Errorf("--id 必填")
			}
			if dryRun {
				return writeRows(f, map[string]interface{}{"id": id}, []string{"动作", "模板ID"}, [][]string{{"启用券模板", strconv.FormatInt(id, 10)}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.EnableCouponTemplate(cmd.Context(), client.CouponTemplateEnableParams{"id": id})
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"模板ID", "启用结果"}, [][]string{{strconv.FormatInt(id, 10), strconv.FormatBool(result.Data)}})
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "券模板ID")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCouponChoose(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	cmd := &cobra.Command{
		Use:   "choose",
		Short: "筛选可用优惠券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ChooseCouponsAsOne(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					strconv.FormatInt(item.CardID, 10),
					item.Title,
					strconv.Itoa(item.UseStatus),
					strconv.Itoa(item.UnUseStatus),
					fmt.Sprintf("%.2f", item.UserCouponCanDisAmount),
					item.EndAt,
				})
			}
			return writeRows(f, result, []string{"券码", "名称", "状态", "可用状态", "可优惠金额", "到期时间"}, rows)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	return cmd
}

func newCmdAnonymousCoupon(f *cmdutil.Factory) *cobra.Command {
	var code string
	cmd := &cobra.Command{
		Use:   "anonymous",
		Short: "查询不记名券详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if code == "" {
				return fmt.Errorf("--code 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetAnonymousCoupon(cmd.Context(), code)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"券码", "模板ID", "模板名称", "状态", "模板面额"}, [][]string{{
				result.Code, strconv.FormatInt(result.TemplateID, 10), result.TemplateName, strconv.Itoa(result.Status), fmt.Sprintf("%.2f", result.Amount),
			}})
		},
	}
	cmd.Flags().StringVar(&code, "code", "", "券码")
	return cmd
}

func newCmdCouponTemplateBatch(f *cmdutil.Factory) *cobra.Command {
	var ids []int64
	cmd := &cobra.Command{
		Use:   "template-batch",
		Short: "批量查询券模板详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(ids) == 0 {
				return fmt.Errorf("--ids 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCouponTemplateListByIDs(cmd.Context(), ids)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					strconv.FormatInt(item.ID, 10), item.Name, strconv.Itoa(item.Status), strconv.Itoa(item.Type), item.ThirdBizCode,
				})
			}
			return writeRows(f, result, []string{"模板ID", "模板名称", "状态", "券类型", "三方编码"}, rows)
		},
	}
	cmd.Flags().Int64SliceVar(&ids, "ids", nil, "券模板 ID 列表")
	return cmd
}

func newCmdCouponTemplateByThirdCode(f *cmdutil.Factory) *cobra.Command {
	var sellerType int
	var thirdBizCode string
	cmd := &cobra.Command{
		Use:   "template-by-third-code",
		Short: "根据三方业务券编码查询券模板",
		RunE: func(cmd *cobra.Command, args []string) error {
			if sellerType == 0 || thirdBizCode == "" {
				return fmt.Errorf("--seller-type 和 --third-biz-code 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCouponTemplateByThirdCode(cmd.Context(), sellerType, thirdBizCode)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"三方业务券编码", "模板ID"}, [][]string{{thirdBizCode, strconv.FormatInt(result.ID, 10)}})
		},
	}
	cmd.Flags().IntVar(&sellerType, "seller-type", 0, "业务类型，2=新饮食")
	cmd.Flags().StringVar(&thirdBizCode, "third-biz-code", "", "三方业务券编码")
	return cmd
}

func newCmdGrantByActivity(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "grant-activity",
		Short: "根据活动实时发券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"根据活动实时发券", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GrantCouponsByActivity(cmd.Context(), client.ActivityRealtimeGrantParams(params))
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{strconv.FormatInt(item.ID, 10), strconv.FormatInt(item.TemplateID, 10)})
			}
			return writeRows(f, result, []string{"券ID", "模板ID"}, rows)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGrantByActivityAsync(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "grant-activity-async",
		Short: "根据活动异步发券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"根据活动异步发券", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			if err := api.GrantCouponsByActivityAsync(cmd.Context(), client.ActivityAsyncGrantParams(params)); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交异步发券: %s\n", fromJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGrantByTemplate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "grant-template",
		Short: "根据券模板实时发券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"根据券模板实时发券", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GrantCouponsByTemplate(cmd.Context(), client.TemplateRealtimeGrantParams(params))
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"券ID", "发券数量", "券ID列表"}, [][]string{{
				result.CardID,
				strconv.Itoa(result.Num),
				fmt.Sprint(result.CardIDList),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCampaignList(f *cmdutil.Factory) *cobra.Command {
	var channelID, pageNo, pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "查询发券活动",
		RunE: func(cmd *cobra.Command, args []string) error {
			if channelID == 0 {
				return fmt.Errorf("--channel-id 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListCouponActivities(cmd.Context(), channelID, pageNo, pageSize)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					strconv.FormatInt(item.ActivityID, 10), item.ActivityCode, item.ActivityName, item.StartDate, item.EndDate, strconv.Itoa(item.RemainderNum),
				})
			}
			if err := writeRows(f, result, []string{"活动ID", "活动编码", "活动名称", "开始时间", "结束时间", "剩余库存"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d/%d 页)\n", result.Total, result.CurrentPage, result.TotalPage)
			return nil
		},
	}
	cmd.Flags().IntVar(&channelID, "channel-id", 0, "渠道 ID")
	cmd.Flags().IntVar(&pageNo, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdExchangeStatus(f *cmdutil.Factory) *cobra.Command {
	var code string
	cmd := &cobra.Command{
		Use:   "exchange-status",
		Short: "查询兑换码状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if code == "" {
				return fmt.Errorf("--code 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryExchangeCodeStatus(cmd.Context(), code)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"兑换码", "活动编码", "活动名称", "状态"}, [][]string{{code, result.ActivityCode, result.ActivityName, strconv.Itoa(result.ExchangeStatus)}})
		},
	}
	cmd.Flags().StringVar(&code, "code", "", "兑换码")
	return cmd
}

func newCmdExchangeDispatch(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "exchange-dispatch",
		Short: "兑换码下发",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"兑换码下发", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.DispatchExchangeCodes(cmd.Context(), client.ExchangeCodeDispatchParams(params))
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.Code})
			}
			return writeRows(f, result, []string{"兑换码"}, rows)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdExchangeDisable(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "exchange-disable",
		Short: "撤销兑换码下发",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"撤销兑换码下发", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.DisableExchangeCodes(cmd.Context(), client.ExchangeCodeDisableParams(params))
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{item.Code, strconv.FormatBool(item.Disabled), item.FailReason})
			}
			return writeRows(f, result, []string{"兑换码", "是否作废成功", "失败原因"}, rows)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdTaskRecords(f *cmdutil.Factory) *cobra.Command {
	var activityIDs []int64
	var customerID int64
	var mobilePhone string
	cmd := &cobra.Command{
		Use:   "task-records",
		Short: "查询指定任务活动用户参与记录",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(activityIDs) == 0 || customerID == 0 {
				return fmt.Errorf("--activity-ids 和 --customer-id 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetActivityTasks(cmd.Context(), activityIDs, customerID, mobilePhone)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.ActivityID, item.ActivityName, strconv.Itoa(item.ActivityStatus), strconv.Itoa(item.CurrentJoinNum), strconv.Itoa(item.ResidueCanJoinNum),
				})
			}
			return writeRows(f, result, []string{"活动ID", "活动名称", "状态", "当前参与次数", "剩余可参与次数"}, rows)
		},
	}
	cmd.Flags().Int64SliceVar(&activityIDs, "activity-ids", nil, "活动 ID 列表")
	cmd.Flags().Int64Var(&customerID, "customer-id", 0, "会员 ID")
	cmd.Flags().StringVar(&mobilePhone, "mobile-phone", "", "手机号")
	return cmd
}

func newCmdRecoverGrant(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "revoke-grant",
		Short: "撤销发券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"撤销发券", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.RecoverGrantedCoupons(cmd.Context(), client.CouponRecoveryParams(params))
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"回收状态", "订单发放券ID", "成功回收券ID"}, [][]string{{
				strconv.Itoa(result.RecoveryStatus),
				fmt.Sprint(result.CouponIDs),
				fmt.Sprint(result.RecoveryIDs),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdTaskClaim(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "task-claim",
		Short: "领取活动任务",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"领取活动任务", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			if err := api.ClaimActivityTask(cmd.Context(), client.ActivityTaskClaimParams(params)); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交任务领取: %s\n", fromJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdRecycleCoupons(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "recycle-coupons",
		Short: "回收券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"回收券", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			if err := api.RecycleCoupons(cmd.Context(), client.CouponRecycleParams(params)); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交回收券: %s\n", fromJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardInfo(f *cmdutil.Factory) *cobra.Command {
	var cardNo string
	var takeTemplate int
	cmd := &cobra.Command{
		Use:   "info",
		Short: "查询卡信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCardInfo(cmd.Context(), cardNo, takeTemplate)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"卡号", "卡名称", "状态", "余额", "会员ID", "模板ID"}, [][]string{{
				result.CardNo, result.Name, strconv.Itoa(result.Status), fmt.Sprintf("%.2f", result.RemainingAmount), strconv.FormatInt(result.CustomerID, 10), strconv.FormatInt(result.TemplateID, 10),
			}})
		},
	}
	cmd.Flags().StringVar(&cardNo, "card-no", "", "卡号")
	cmd.Flags().IntVar(&takeTemplate, "take-card-template", 0, "是否返回卡模板，1=是")
	return cmd
}

func newCmdGiftCardConsume(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "consume",
		Short: "礼品卡余额扣减",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"礼品卡余额扣减", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ConsumeGiftCard(cmd.Context(), client.GiftCardConsumeParams(params))
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"卡号", "成功", "卡类型", "会员ID", "扣减总额"}, [][]string{{
				result.CardConsumeInfo.CardNo,
				strconv.FormatBool(result.Success),
				strconv.Itoa(result.CardType),
				strconv.FormatInt(result.CustomerID, 10),
				fmt.Sprintf("%.2f", result.CardConsumeInfo.DecreaseTotalBalance),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardConsumeBatch(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "consume-batch",
		Short: "批量扣减礼品卡余额",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"批量扣减礼品卡余额", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ConsumeMultiGiftCards(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.CardConsumeInfo))
			for _, item := range result.CardConsumeInfo {
				rows = append(rows, []string{
					item.CardConsumeInfo.CardNo,
					strconv.FormatBool(item.Success),
					fmt.Sprintf("%.2f", item.CardConsumeInfo.DecreaseTotalBalance),
					strconv.FormatInt(item.CustomerID, 10),
				})
			}
			return writeRows(f, result, []string{"卡号", "成功", "扣减总额", "会员ID"}, rows)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardReverse(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "reverse",
		Short: "礼品卡余额冲正",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"礼品卡余额冲正", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ReverseGiftCardConsume(cmd.Context(), client.GiftCardConsumeReverseParams(params))
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"卡号", "成功", "卡类型", "会员ID", "冲正后充值余额"}, [][]string{{
				result.CardConsumeInfo.CardNo,
				strconv.FormatBool(result.Success),
				strconv.Itoa(result.CardType),
				strconv.FormatInt(result.CustomerID, 10),
				fmt.Sprintf("%.2f", result.CardConsumeInfo.AfterRechargeBalance),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardReverseBatch(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "reverse-batch",
		Short: "批量冲正礼品卡余额",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"批量冲正礼品卡余额", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			if err := api.ReverseMultiGiftCards(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交批量礼品卡冲正: %s\n", fromJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardPartReverse(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "part-reverse",
		Short: "礼品卡余额部分冲正",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"礼品卡余额部分冲正", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.PartReverseGiftCardConsume(cmd.Context(), client.GiftCardConsumePartReverseParams(params))
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"卡号", "消费业务ID", "冲正业务ID", "冲正总额", "冲正后总余额"}, [][]string{{
				result.CardNo,
				result.BizID,
				result.SubBizID,
				fmt.Sprintf("%.2f", result.ReverseTotalBalance),
				fmt.Sprintf("%.2f", result.AfterTotalBalance),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardIssue(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "发放礼品卡",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"发放礼品卡", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			if err := api.BindGiftCard(cmd.Context(), client.GiftCardBindParams(params)); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交礼品卡发放: %s\n", fromJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardReportLoss(f *cmdutil.Factory) *cobra.Command {
	return newGiftCardActionCmd(f, "report-loss", "礼品卡挂失", "礼品卡挂失", func(api *client.MarketingAPI, cmd *cobra.Command, params map[string]interface{}) error {
		return api.ReportGiftCardLoss(cmd.Context(), client.GiftCardReportLossParams(params))
	})
}

func newCmdGiftCardRelieveLoss(f *cmdutil.Factory) *cobra.Command {
	return newGiftCardActionCmd(f, "relieve-loss", "取消礼品卡挂失", "取消礼品卡挂失", func(api *client.MarketingAPI, cmd *cobra.Command, params map[string]interface{}) error {
		return api.RelieveGiftCardLoss(cmd.Context(), client.GiftCardRelieveLossParams(params))
	})
}

func newCmdGiftCardRecycle(f *cmdutil.Factory) *cobra.Command {
	return newGiftCardActionCmd(f, "recycle", "回收礼品卡", "回收礼品卡", func(api *client.MarketingAPI, cmd *cobra.Command, params map[string]interface{}) error {
		return api.RecycleGiftCard(cmd.Context(), client.GiftCardRecycleParams(params))
	})
}

func newCmdGiftCardExchange(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "exchange",
		Short: "兑换礼品卡",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"兑换礼品卡", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ExchangeGiftCard(cmd.Context(), client.GiftCardExchangeParams(params))
			if err != nil {
				return err
			}
			return writeGiftCardGrantRows(f, result)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardGrantByTemplate(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "grant-template",
		Short: "根据模板发放礼品卡",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{"根据模板发放礼品卡", fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.BatchGrantGiftCard(cmd.Context(), client.GiftCardBatchGrantParams(params))
			if err != nil {
				return err
			}
			return writeGiftCardGrantRows(f, result)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdGiftCardFlow(f *cmdutil.Factory) *cobra.Command {
	var cardNo string
	var changeTypes []int
	var isDescFlag int
	var pageNum int
	var pageSize int
	cmd := &cobra.Command{
		Use:   "flow",
		Short: "查询礼品卡消费明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cardNo == "" {
				return fmt.Errorf("--card-no 必填")
			}
			params := map[string]interface{}{"cardNo": cardNo, "pageNum": pageNum, "pageSize": pageSize}
			if len(changeTypes) > 0 {
				params["changeTypes"] = changeTypes
			}
			if isDescFlag != 0 {
				params["isDescFlag"] = isDescFlag
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListGiftCardFlows(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.CreatedAt, strconv.Itoa(item.ActionType), fmt.Sprintf("%.2f", item.UseAmount), fmt.Sprintf("%.2f", item.Balance), item.UseOrderCode, item.UseStoreName,
				})
			}
			if err := writeRows(f, result, []string{"时间", "类型", "核销金额", "余额", "订单号", "门店"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d/%d 页)\n", result.TotalCount, result.PageNum, result.PageCount)
			return nil
		},
	}
	cmd.Flags().StringVar(&cardNo, "card-no", "", "卡号")
	cmd.Flags().IntSliceVar(&changeTypes, "change-types", nil, "变更类型列表")
	cmd.Flags().IntVar(&isDescFlag, "is-desc", 0, "1=倒序")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdGiftCardTemplate(f *cmdutil.Factory) *cobra.Command {
	var id int64
	cmd := &cobra.Command{
		Use:   "template",
		Short: "查询礼品卡模板详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return fmt.Errorf("--id 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetGiftCardTemplateDetail(cmd.Context(), id)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"模板ID", "名称", "品牌ID", "状态", "卡类型", "三方编码"}, [][]string{{
				result.ID, result.Name, strconv.FormatInt(result.SellerID, 10), strconv.Itoa(result.Status), strconv.Itoa(result.Type), result.ThirdBizCode,
			}})
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "卡模板 ID")
	return cmd
}

func newCmdGiftCardList(f *cmdutil.Factory) *cobra.Command {
	var cardNos []string
	var cardStatus int
	var customerID int64
	var customerCode string
	var customerCodeType int
	var shopCode string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "查询会员礼品卡/实体卡列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{}
			if len(cardNos) > 0 {
				params["cardNos"] = cardNos
			}
			if cardStatus != 0 {
				params["cardStatus"] = cardStatus
			}
			if customerID != 0 {
				params["customerId"] = customerID
			} else if customerCode != "" && customerCodeType != 0 {
				params["customerCodeType"] = map[string]interface{}{"customerCode": customerCode, "type": customerCodeType}
			} else {
				return fmt.Errorf("--customer-id 或 (--customer-code 与 --customer-code-type) 必填")
			}
			if shopCode != "" {
				params["shopCode"] = shopCode
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListCustomerGiftCards(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.CardNo, item.Name, strconv.Itoa(item.CardStatus), fmt.Sprintf("%.2f", item.TotalBalance), strconv.Itoa(item.CardType), item.EndAt,
				})
			}
			return writeRows(f, result, []string{"卡号", "卡名称", "状态", "总余额", "卡类型", "到期时间"}, rows)
		},
	}
	cmd.Flags().StringSliceVar(&cardNos, "card-nos", nil, "卡号列表")
	cmd.Flags().IntVar(&cardStatus, "card-status", 0, "卡状态")
	cmd.Flags().Int64Var(&customerID, "customer-id", 0, "会员 ID")
	cmd.Flags().StringVar(&customerCode, "customer-code", "", "会员标识")
	cmd.Flags().IntVar(&customerCodeType, "customer-code-type", 0, "会员标识类型")
	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码")
	return cmd
}

func newCmdGiftCardTemplateBatch(f *cmdutil.Factory) *cobra.Command {
	var ids []int64
	cmd := &cobra.Command{
		Use:   "template-batch",
		Short: "批量查询礼品卡模板",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(ids) == 0 {
				return fmt.Errorf("--ids 必填")
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListGiftCardTemplates(cmd.Context(), ids)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.ID, item.Name, strconv.Itoa(item.Status), strconv.Itoa(item.Type), fmt.Sprintf("%.2f", item.SalePrice), item.ThirdBizCode,
				})
			}
			return writeRows(f, result, []string{"模板ID", "名称", "状态", "卡类型", "售价", "三方编码"}, rows)
		},
	}
	cmd.Flags().Int64SliceVar(&ids, "ids", nil, "卡模板 ID 列表")
	return cmd
}

func newCmdPromotionActivities(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	cmd := &cobra.Command{
		Use:   "activities",
		Short: "查询门店促销活动列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListStorePromotionActivities(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.ActivityCode, item.ActivityName, item.ActivityTitle, item.ActivityType, strconv.Itoa(item.ActivityStatus),
				})
			}
			return writeRows(f, result, []string{"活动编码", "活动名称", "活动标题", "活动类型", "活动状态"}, rows)
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	return cmd
}

func newCmdConfirm(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	cmd := &cobra.Command{
		Use:   "confirm",
		Short: "结算页算价",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ConfirmOrderPricing(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeRows(f, result, []string{"总价", "实付", "优惠", "商品金额", "商品小计", "餐盒费"}, [][]string{{
				fmt.Sprintf("%.2f", result.OrderFee.TotalAmount),
				fmt.Sprintf("%.2f", result.OrderFee.ActualAmount),
				fmt.Sprintf("%.2f", result.OrderFee.DiscountAmount),
				fmt.Sprintf("%.2f", result.OrderFee.GoodsAmount),
				fmt.Sprintf("%.2f", result.OrderFee.GoodsSubTotal),
				fmt.Sprintf("%.2f", result.OrderFee.MealBoxAmount),
			}})
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
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

func newGiftCardActionCmd(f *cmdutil.Factory, use, short, action string, invoke func(api *client.MarketingAPI, cmd *cobra.Command, params map[string]interface{}) error) *cobra.Command {
	var fromJSON string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeRows(f, params, []string{"动作", "文件"}, [][]string{{action, fromJSON}})
			}
			api, err := newMarketingAPI(f)
			if err != nil {
				return err
			}
			if err := invoke(api, cmd, params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交%s: %s\n", action, fromJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func writeGiftCardGrantRows(f *cmdutil.Factory, result *[]client.GiftCardGrantItem) error {
	rows := make([][]string, 0, len(*result))
	for _, item := range *result {
		rows = append(rows, []string{
			item.CardNo,
			strconv.FormatInt(item.TemplateID, 10),
			item.Name,
			fmt.Sprintf("%.2f", item.Amount),
			fmt.Sprintf("%.2f", item.SalePrice),
			strconv.Itoa(item.Type),
		})
	}
	return writeRows(f, result, []string{"卡号", "模板ID", "名称", "面值", "售价", "卡类型"}, rows)
}

func writeRows(f *cmdutil.Factory, data interface{}, headers []string, rows [][]string) error {
	format, err := output.ParseFormat(f.EffectiveFormat())
	if err != nil {
		return err
	}
	fmtr := output.NewFormatter(f.IOStreams.Out, format)
	return fmtr.Write(data, headers, rows)
}
