package member

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

func NewCmdMember(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "会员服务",
		Long:  "会员资产、会员信息、会员标签查询与变更",
	}

	cmd.AddCommand(newCmdAsset(f))
	cmd.AddCommand(newCmdProfile(f))
	cmd.AddCommand(newCmdTag(f))

	return cmd
}

func newMemberAPI(f *cmdutil.Factory) (*client.MemberAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewMemberAPI(apiClient), nil
}

func newCmdAsset(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset",
		Short: "会员资产查询与变更",
	}
	cmd.AddCommand(newCmdPoints(f))
	cmd.AddCommand(newCmdCoupons(f))
	cmd.AddCommand(newCmdCouponList(f))
	cmd.AddCommand(newCmdBalance(f))
	cmd.AddCommand(newCmdBalanceFlow(f))
	cmd.AddCommand(newCmdBalanceDecrease(f))
	cmd.AddCommand(newCmdDiscount(f))
	cmd.AddCommand(newCmdLevelExperience(f))
	cmd.AddCommand(newCmdPointsFlow(f))
	cmd.AddCommand(newCmdPersonalAsset(f))
	cmd.AddCommand(newCmdDepositRules(f))
	cmd.AddCommand(newCmdInflateStatus(f))
	cmd.AddCommand(newCmdCouponDetailList(f))
	cmd.AddCommand(newCmdRecharge(f))
	cmd.AddCommand(newCmdBalanceDebit(f))
	cmd.AddCommand(newCmdBalanceReverse(f))
	cmd.AddCommand(newCmdPointsDebit(f))
	cmd.AddCommand(newCmdPointsReverse(f))
	cmd.AddCommand(newCmdPointsAdd(f))
	cmd.AddCommand(newCmdConsume(f))
	cmd.AddCommand(newCmdConsumeReverse(f))
	cmd.AddCommand(newCmdCouponWriteOff(f))
	cmd.AddCommand(newCmdCouponReverse(f))
	cmd.AddCommand(newCmdOfflineBalanceOperation(f))
	cmd.AddCommand(newCmdRechargeReverseOp(f))
	return cmd
}

func newCmdProfile(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "会员信息查询",
	}
	cmd.AddCommand(newCmdInfo(f))
	cmd.AddCommand(newCmdIDByCode(f))
	cmd.AddCommand(newCmdIDsByPhone(f))
	cmd.AddCommand(newCmdBaseInfo(f))
	cmd.AddCommand(newCmdSearch(f))
	cmd.AddCommand(newCmdRiskLevel(f))
	cmd.AddCommand(newCmdLevel(f))
	cmd.AddCommand(newCmdOpenID(f))
	cmd.AddCommand(newCmdDynamicCode(f))
	cmd.AddCommand(newCmdRegisterPhone(f))
	cmd.AddCommand(newCmdRegisterThirdID(f))
	cmd.AddCommand(newCmdSendCaptcha(f))
	cmd.AddCommand(newCmdCheckCaptcha(f))
	cmd.AddCommand(newCmdBlacklistAdd(f))
	cmd.AddCommand(newCmdBlacklistRemove(f))
	cmd.AddCommand(newCmdBlacklistStatus(f))
	cmd.AddCommand(newCmdFreeze(f))
	cmd.AddCommand(newCmdUnfreeze(f))
	cmd.AddCommand(newCmdLogoff(f))
	cmd.AddCommand(newCmdUpdatePhone(f))
	cmd.AddCommand(newCmdFreezeRecord(f))
	cmd.AddCommand(newCmdHasOpenOrder(f))
	cmd.AddCommand(newCmdWeComInfo(f))
	cmd.AddCommand(newCmdSignStatus(f))
	cmd.AddCommand(newCmdAccountLevel(f))
	cmd.AddCommand(newCmdConditionQuery(f))
	cmd.AddCommand(newCmdUpdate(f))
	return cmd
}

func newCmdTag(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "会员标签查询",
	}
	cmd.AddCommand(newCmdTags(f))
	cmd.AddCommand(newCmdTagDetail(f))
	cmd.AddCommand(newCmdBrandTags(f))
	cmd.AddCommand(newCmdTagIDs(f))
	cmd.AddCommand(newCmdGroupList(f))
	cmd.AddCommand(newCmdGroupLabels(f))
	cmd.AddCommand(newCmdCreateTag(f))
	cmd.AddCommand(newCmdDeleteTag(f))
	cmd.AddCommand(newCmdClearTagMembers(f))
	cmd.AddCommand(newCmdDeleteCustomerTag(f))
	cmd.AddCommand(newCmdTagSettings(f))
	cmd.AddCommand(newCmdMark(f))
	return cmd
}

func newCmdPoints(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var orderTime string

	cmd := &cobra.Command{
		Use:   "points",
		Short: "查询会员积分",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerPoints(cmd.Context(), customerID, orderTime)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "总积分", "永久积分", "即将过期积分", "过期时间"}, [][]string{{
				customerID,
				fmt.Sprintf("%.2f", result.TotalPoints),
				fmt.Sprintf("%.2f", result.ForeverPoints),
				fmt.Sprintf("%.2f", result.SoonExpiredPoints),
				result.ExpiredTime,
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&orderTime, "order-time", "", "下单时间，格式参考企迈文档")
	return cmd
}

func newCmdCoupons(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var cardID string
	var templateID string
	var useStatus int

	cmd := &cobra.Command{
		Use:   "coupons",
		Short: "查询会员优惠券",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerCoupons(cmd.Context(), customerID, cardID, templateID, useStatus)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.CardID,
					item.TemplateID,
					item.TemplateName,
					strconv.Itoa(item.CouponType),
					strconv.Itoa(item.UseStatus),
					item.EndTime,
				})
			}
			if err := writeSingleRows(f, result, []string{"卡券ID", "模板ID", "模板名称", "券类型", "状态", "结束时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条\n", result.Total)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&cardID, "card-id", "", "卡券 ID")
	cmd.Flags().StringVar(&templateID, "template-id", "", "模板 ID")
	cmd.Flags().IntVar(&useStatus, "use-status", 0, "使用状态")
	return cmd
}

func newCmdCouponList(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var couponType int
	var useStatus int
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "coupon-list",
		Short: "分页查询会员优惠券列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			if pageNum == 0 {
				pageNum = 1
			}
			if pageSize == 0 {
				pageSize = 20
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerCouponList(cmd.Context(), customerID, couponType, pageNum, pageSize, useStatus)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.CardID,
					item.TemplateID,
					item.TemplateName,
					strconv.Itoa(item.CouponType),
					strconv.Itoa(item.UseStatus),
					item.EndTime,
				})
			}
			if err := writeSingleRows(f, result, []string{"卡券ID", "模板ID", "模板名称", "券类型", "状态", "结束时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, pageNum, pageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().IntVar(&couponType, "coupon-type", 0, "券类型")
	cmd.Flags().IntVar(&useStatus, "use-status", 0, "使用状态")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdBalance(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "查询会员储值余额",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerBalance(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "总余额", "充值余额", "赠送余额", "冻结余额"}, [][]string{{
				result.CustomerID,
				fmt.Sprintf("%.2f", result.TotalBalance),
				fmt.Sprintf("%.2f", result.RechargeBalance),
				fmt.Sprintf("%.2f", result.GiftBalance),
				fmt.Sprintf("%.2f", result.FreezeBalance),
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdBalanceFlow(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var orderNo string
	var startTime string
	var endTime string
	var changeType string
	var changeTypes []string
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "balance-flow",
		Short: "查询会员储值余额明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			params := map[string]interface{}{
				"customerId": customerID,
				"pageNumber": pageNum,
				"pageSize":   pageSize,
			}
			if orderNo != "" {
				params["orderNo"] = orderNo
			}
			if startTime != "" {
				params["startTime"] = startTime
			}
			if endTime != "" {
				params["endTime"] = endTime
			}
			if changeType != "" {
				params["changeType"] = changeType
			}
			if len(changeTypes) > 0 {
				params["changeTypes"] = changeTypes
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryAccountFlow(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.ChangeTime,
					item.ChangeType,
					fmt.Sprintf("%.2f", item.Amount),
					fmt.Sprintf("%.2f", item.AfterAmount),
					item.OrderNo,
					item.Remark,
				})
			}
			if err := writeSingleRows(f, result, []string{"变更时间", "变更类型", "变更金额", "变更后余额", "订单号", "备注"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, pageNum, pageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&orderNo, "order-no", "", "订单号")
	cmd.Flags().StringVar(&startTime, "start-time", "", "开始时间")
	cmd.Flags().StringVar(&endTime, "end-time", "", "结束时间")
	cmd.Flags().StringVar(&changeType, "change-type", "", "变更类型")
	cmd.Flags().StringSliceVar(&changeTypes, "change-types", nil, "变更类型列表")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdBalanceDecrease(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var bizID string

	cmd := &cobra.Command{
		Use:   "balance-decrease",
		Short: "查询会员储值余额扣减明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" || bizID == "" {
				return fmt.Errorf("--customer-id 和 --biz-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetDecreaseBalanceBiz(cmd.Context(), customerID, bizID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "业务ID", "金额", "时间", "订单号"}, [][]string{{
				result.CustomerID,
				result.BizID,
				fmt.Sprintf("%.2f", result.Amount),
				result.ChangeTime,
				result.OrderNo,
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&bizID, "biz-id", "", "业务 ID")
	return cmd
}

func newCmdDiscount(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "discount",
		Short: "查询会员实时折扣",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetLevelConsumeDiscount(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "折扣"}, [][]string{{result.CustomerID, fmt.Sprintf("%.2f", result.Discount)}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdLevelExperience(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "level-experience",
		Short: "查询会员等级和经验值",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetLevelExperience(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "等级", "等级名称", "经验值"}, [][]string{{
				result.CustomerID,
				strconv.Itoa(result.Level),
				result.LevelName,
				fmt.Sprintf("%.2f", result.TotalExperienceNum),
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdPointsFlow(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var orgType string
	var pointsChannel string
	var sourceType string
	var changeTypes []string
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "points-flow",
		Short: "查询会员积分明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			params := map[string]interface{}{
				"customerId": customerID,
				"pageNum":    pageNum,
				"pageSize":   pageSize,
			}
			if orgType != "" {
				params["orgType"] = orgType
			}
			if pointsChannel != "" {
				params["pointsChannel"] = pointsChannel
			}
			if sourceType != "" {
				params["sourceType"] = sourceType
			}
			if len(changeTypes) > 0 {
				params["changeTypeList"] = changeTypes
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCrmPointsFlow(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.ChangeTime,
					item.ChangeType,
					fmt.Sprintf("%.2f", item.ChangePoints),
					fmt.Sprintf("%.2f", item.AfterPoints),
					item.OrgType,
					item.SourceType,
				})
			}
			if err := writeSingleRows(f, result, []string{"变更时间", "变更类型", "变更积分", "变更后积分", "机构类型", "来源类型"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, pageNum, pageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&orgType, "org-type", "", "机构类型")
	cmd.Flags().StringVar(&pointsChannel, "points-channel", "", "积分渠道")
	cmd.Flags().StringVar(&sourceType, "source-type", "", "来源类型")
	cmd.Flags().StringSliceVar(&changeTypes, "change-types", nil, "变更类型列表")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdPersonalAsset(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "personal-asset",
		Short: "查询用户资产",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetPersonalAsset(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "资产余额", "储值余额", "积分", "优惠券数"}, [][]string{{
				result.CustomerID,
				fmt.Sprintf("%.2f", result.Balance),
				fmt.Sprintf("%.2f", result.StoredBalance),
				fmt.Sprintf("%.2f", result.Points),
				strconv.Itoa(result.CouponCount),
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdDepositRules(f *cmdutil.Factory) *cobra.Command {
	var shopCode string
	var depositValue float64

	cmd := &cobra.Command{
		Use:   "deposit-rules",
		Short: "查询门店储值规则",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetStoreDepositRules(cmd.Context(), shopCode, depositValue)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					strconv.FormatInt(item.ID, 10),
					item.ActivityName,
					fmt.Sprintf("%.2f", item.DetailInfo.DepositValue),
					fmt.Sprintf("%.2f", item.DetailInfo.ChargeGiftBalance),
					item.DetailInfo.GiftIntegral,
				})
			}
			return writeSingleRows(f, result, []string{"规则ID", "活动名称", "储值金额", "赠送金额", "赠送积分"}, rows)
		},
	}

	cmd.Flags().StringVar(&shopCode, "shop-code", "", "门店编码，不传则查全部门店")
	cmd.Flags().Float64Var(&depositValue, "deposit-value", 0, "储值金额（元）")
	return cmd
}

func newCmdInflateStatus(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "inflate-status",
		Short: "查询资产膨胀状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetAssetInflateDetail(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "是否已膨胀", "膨胀前等级", "膨胀时间"}, [][]string{{
				result.CustomerID,
				strconv.Itoa(result.IsInflate),
				strconv.Itoa(result.InflateBeforeLevel),
				result.InflateTime,
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdCouponDetailList(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var storeID int64
	var couponType int
	var couponTableType int
	var effectiveType int
	var needQueryStore int
	var sortType int
	var title string
	var useStatus int
	var useScene int
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "coupon-detail-list",
		Short: "查询优惠券列表明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" || storeID == 0 {
				return fmt.Errorf("--customer-id 和 --store-id 必填")
			}
			params := map[string]interface{}{
				"customerId": customerID,
				"store_id":   storeID,
				"pageNum":    pageNum,
				"pageSize":   pageSize,
			}
			if couponType != 0 {
				params["couponType"] = couponType
			}
			if couponTableType != 0 {
				params["couponTableType"] = couponTableType
			}
			if effectiveType != 0 {
				params["effectiveType"] = effectiveType
			}
			if needQueryStore != 0 {
				params["needQueryStore"] = needQueryStore
			}
			if sortType != 0 {
				params["sortType"] = sortType
			}
			if title != "" {
				params["title"] = title
			}
			if useStatus != 0 {
				params["useStatus"] = useStatus
			}
			if useScene != 0 {
				params["useScene"] = useScene
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCouponDetailList(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.CardID,
					item.TemplateID,
					item.Title,
					fmt.Sprintf("%.2f", item.FaceAmount),
					strconv.Itoa(item.UseStatus),
					item.EndAt,
				})
			}
			if err := writeSingleRows(f, result, []string{"卡券ID", "模板ID", "标题", "面额", "状态", "到期时间"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页)\n", result.Total, result.PageNum)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().Int64Var(&storeID, "store-id", 0, "访问门店 ID")
	cmd.Flags().IntVar(&couponType, "coupon-type", 0, "券类型")
	cmd.Flags().IntVar(&couponTableType, "coupon-table-type", 0, "券大类")
	cmd.Flags().IntVar(&effectiveType, "effective-type", 0, "生效状态")
	cmd.Flags().IntVar(&needQueryStore, "need-query-store", 0, "是否查询门店，1=不需要")
	cmd.Flags().IntVar(&sortType, "sort-type", 0, "排序方式")
	cmd.Flags().StringVar(&title, "title", "", "优惠券名称或券号")
	cmd.Flags().IntVar(&useStatus, "use-status", 0, "券状态")
	cmd.Flags().IntVar(&useScene, "use-scene", 0, "使用场景")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdRecharge(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "recharge", "会员储值充值", "会员储值充值", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.RechargeBalance(ctx.Context(), client.MemberRechargeParams(params))
	})
}

func newCmdBalanceDebit(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "balance-debit", "会员储值余额扣减", "会员储值余额扣减", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.DecreaseBalance(ctx.Context(), client.MemberBalanceDecreaseParams(params))
	})
}

func newCmdBalanceReverse(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "balance-reverse", "会员储值余额冲正", "会员储值余额冲正", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.ReverseBalance(ctx.Context(), client.MemberBalanceReverseParams(params))
	})
}

func newCmdPointsDebit(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "points-debit", "会员积分扣减", "会员积分扣减", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.ReducePoints(ctx.Context(), client.MemberPointsDecreaseParams(params))
	})
}

func newCmdPointsReverse(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "points-reverse", "会员积分冲正", "会员积分冲正", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.ReversePoints(ctx.Context(), client.MemberPointsReverseParams(params))
	})
}

func newCmdPointsAdd(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "points-add", "会员积分发放", "会员积分发放", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.AddPoints(ctx.Context(), client.MemberPointsAddParams(params))
	})
}

func newCmdConsume(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "consume", "会员资产消费", "会员资产消费", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.CustomerConsume(ctx.Context(), client.MemberConsumeParams(params))
	})
}

func newCmdConsumeReverse(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "consume-reverse", "会员资产冲正", "会员资产冲正", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.ConsumeReverse(ctx.Context(), client.MemberConsumeReverseParams(params))
	})
}

func newCmdCouponWriteOff(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "coupon-writeoff", "优惠券核销", "优惠券核销", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.WriteOffCoupon(ctx.Context(), client.MemberCouponWriteOffParams(params))
	})
}

func newCmdCouponReverse(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "coupon-reverse", "优惠券冲正", "优惠券冲正", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.ReverseCoupon(ctx.Context(), client.MemberCouponReverseParams(params))
	})
}

func newCmdOfflineBalanceOperation(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "offline-balance-op", "储值账户充值扣减", "储值账户充值扣减", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.OfflineOperateBalance(ctx.Context(), client.OfflineBalanceOperationParams(params))
	})
}

func newCmdRechargeReverseOp(f *cmdutil.Factory) *cobra.Command {
	return newJSONAssetMutationCmd(f, "recharge-reverse", "会员充值冲正", "会员充值冲正", func(api *client.MemberAPI, ctx *cobra.Command, params map[string]interface{}) error {
		return api.RechargeReverse(ctx.Context(), client.MemberRechargeReverseParams(params))
	})
}

func newCmdInfo(f *cmdutil.Factory) *cobra.Command {
	var identifierNumber string
	var queryType int
	var typeCate int
	var discountPrice bool
	var multiMark string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "会员信息查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{}
			if identifierNumber != "" {
				params["identifierNumber"] = identifierNumber
			}
			if queryType != 0 {
				params["type"] = queryType
			}
			if typeCate != 0 {
				params["typeCate"] = typeCate
			}
			if discountPrice {
				params["discountPrice"] = true
			}
			if multiMark != "" {
				params["multiMark"] = multiMark
			}
			if len(params) == 0 {
				return fmt.Errorf("至少提供一个查询条件")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerInfo(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "手机号", "昵称", "姓名", "卡号"}, [][]string{{
				result.CustomerID,
				result.MobilePhone,
				result.NickName,
				result.RealName,
				result.CardNo,
			}})
		},
	}

	cmd.Flags().StringVar(&identifierNumber, "identifier-number", "", "会员标识")
	cmd.Flags().IntVar(&queryType, "type", 0, "查询类型")
	cmd.Flags().IntVar(&typeCate, "type-cate", 0, "标识类别")
	cmd.Flags().BoolVar(&discountPrice, "discount-price", false, "是否返回折扣价信息")
	cmd.Flags().StringVar(&multiMark, "multi-mark", "", "多值标识")
	return cmd
}

func newCmdIDByCode(f *cmdutil.Factory) *cobra.Command {
	var code string
	var codeType int

	cmd := &cobra.Command{
		Use:   "id-by-code",
		Short: "通过会员标识查询会员 ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if code == "" || codeType == 0 {
				return fmt.Errorf("--code 和 --type 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerIDByCode(cmd.Context(), code, codeType)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员标识", "类型", "会员ID"}, [][]string{{code, strconv.Itoa(codeType), result.CustomerID}})
		},
	}

	cmd.Flags().StringVar(&code, "code", "", "会员标识值")
	cmd.Flags().IntVar(&codeType, "type", 0, "会员标识类型")
	return cmd
}

func newCmdIDsByPhone(f *cmdutil.Factory) *cobra.Command {
	var phones []string

	cmd := &cobra.Command{
		Use:   "ids-by-phone",
		Short: "通过手机号批量查询会员 ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(phones) == 0 {
				return fmt.Errorf("--phones 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerIDByPhone(cmd.Context(), phones)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{item.Phone, item.CustomerID})
			}
			return writeSingleRows(f, result, []string{"手机号", "会员ID"}, rows)
		},
	}

	cmd.Flags().StringSliceVar(&phones, "phones", nil, "手机号列表")
	return cmd
}

func newCmdBaseInfo(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var infoType int

	cmd := &cobra.Command{
		Use:   "base-info",
		Short: "会员基础信息查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetBaseInfo(cmd.Context(), customerID, infoType)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "手机号", "昵称", "姓名", "卡号"}, [][]string{{
				result.CustomerID,
				result.MobilePhone,
				result.NickName,
				result.RealName,
				result.CardNo,
			}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().IntVar(&infoType, "type", 0, "查询类型，默认按文档服务端默认值")
	return cmd
}

func newCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var phone string
	var openID string
	var unionID string
	var cardNo string
	var title string
	var userID string
	var queryType int
	var couponType int
	var useStatus int
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "search",
		Short: "会员查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{
				"pageNum":  pageNum,
				"pageSize": pageSize,
			}
			if customerID != "" {
				params["customerId"] = customerID
			}
			if phone != "" {
				params["phone"] = phone
			}
			if openID != "" {
				params["openId"] = openID
			}
			if unionID != "" {
				params["unionId"] = unionID
			}
			if cardNo != "" {
				params["cardNo"] = cardNo
			}
			if title != "" {
				params["title"] = title
			}
			if userID != "" {
				params["userId"] = userID
			}
			if queryType != 0 {
				params["type"] = queryType
			}
			if couponType != 0 {
				params["couponType"] = couponType
			}
			if useStatus != 0 {
				params["useStatus"] = useStatus
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.SearchMembers(cmd.Context(), params)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.CustomerID,
					item.MobilePhone,
					item.NickName,
					item.LevelName,
					item.CardNo,
					strconv.Itoa(item.UseStatus),
				})
			}
			if err := writeSingleRows(f, result, []string{"会员ID", "手机号", "昵称", "等级", "卡号", "状态"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, pageNum, pageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&phone, "phone", "", "手机号")
	cmd.Flags().StringVar(&openID, "open-id", "", "OpenID")
	cmd.Flags().StringVar(&unionID, "union-id", "", "UnionID")
	cmd.Flags().StringVar(&cardNo, "card-no", "", "会员卡号")
	cmd.Flags().StringVar(&title, "title", "", "称谓/搜索标题")
	cmd.Flags().StringVar(&userID, "user-id", "", "用户 ID")
	cmd.Flags().IntVar(&queryType, "type", 0, "查询类型")
	cmd.Flags().IntVar(&couponType, "coupon-type", 0, "券类型")
	cmd.Flags().IntVar(&useStatus, "use-status", 0, "状态")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdRiskLevel(f *cmdutil.Factory) *cobra.Command {
	var code string
	var codeType int

	cmd := &cobra.Command{
		Use:   "risk-level",
		Short: "查询会员风险等级",
		RunE: func(cmd *cobra.Command, args []string) error {
			if code == "" || codeType == 0 {
				return fmt.Errorf("--code 和 --type 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetRiskLevel(cmd.Context(), code, codeType)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"标识", "类型", "风险等级"}, [][]string{{result.Code, strconv.Itoa(result.Type), result.RiskLevel}})
		},
	}

	cmd.Flags().StringVar(&code, "code", "", "会员标识")
	cmd.Flags().IntVar(&codeType, "type", 0, "会员标识类型")
	return cmd
}

func newCmdLevel(f *cmdutil.Factory) *cobra.Command {
	var mobilePhone string

	cmd := &cobra.Command{
		Use:   "level",
		Short: "查询会员等级",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mobilePhone == "" {
				return fmt.Errorf("--mobile-phone 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetMemberLevel(cmd.Context(), mobilePhone)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"手机号", "等级", "等级名称", "卡号"}, [][]string{{
				result.MobilePhone,
				strconv.Itoa(result.Level),
				result.LevelName,
				result.CardNo,
			}})
		},
	}

	cmd.Flags().StringVar(&mobilePhone, "mobile-phone", "", "手机号")
	return cmd
}

func newCmdRegisterPhone(f *cmdutil.Factory) *cobra.Command {
	var mobilePhone string
	var regAppType int
	var username string
	var birthday string
	var gender int
	var remark string
	var multiMark string
	var recruitChannel int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "register-phone",
		Short: "通过手机号注册会员",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mobilePhone == "" || regAppType == 0 {
				return fmt.Errorf("--mobile-phone 和 --reg-app-type 必填")
			}
			params := client.PhoneRegisterParams{
				Birthday:       birthday,
				Gender:         gender,
				MobilePhone:    mobilePhone,
				MultiMark:      multiMark,
				RecruitChannel: recruitChannel,
				RegAppType:     regAppType,
				Remark:         remark,
				Username:       username,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将注册会员手机号 %s\n", mobilePhone)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.PhoneRegister(cmd.Context(), params)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "手机号", "姓名", "创建时间"}, [][]string{{
				result.CustomerID,
				result.Phone,
				result.Username,
				result.CreatedAt,
			}})
		},
	}

	cmd.Flags().StringVar(&mobilePhone, "mobile-phone", "", "手机号")
	cmd.Flags().IntVar(&regAppType, "reg-app-type", 0, "注册来源类型")
	cmd.Flags().StringVar(&username, "username", "", "姓名")
	cmd.Flags().StringVar(&birthday, "birthday", "", "生日，格式 yyyy-MM-dd")
	cmd.Flags().IntVar(&gender, "gender", 0, "性别，0=未知 1=男 2=女")
	cmd.Flags().StringVar(&remark, "remark", "", "备注")
	cmd.Flags().StringVar(&multiMark, "multi-mark", "", "门店编码")
	cmd.Flags().IntVar(&recruitChannel, "recruit-channel", 0, "来源渠道")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var id string
	var username string
	var nickname string
	var birthday string
	var gender int
	var idNumber string
	var province string
	var city string
	var district string
	var country string
	var avatar string
	var qmFrom string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "更新会员基本信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("--id 必填")
			}
			params := client.UpdateMemberParams{
				Avatar:   avatar,
				Birthday: birthday,
				City:     city,
				Country:  country,
				District: district,
				Gender:   gender,
				ID:       id,
				IDNumber: idNumber,
				Nickname: nickname,
				Province: province,
				QMFrom:   qmFrom,
				Username: username,
			}
			if params.Avatar == "" && params.Birthday == "" && params.City == "" && params.Country == "" &&
				params.District == "" && params.Gender == 0 && params.IDNumber == "" && params.Nickname == "" &&
				params.Province == "" && params.QMFrom == "" && params.Username == "" {
				return fmt.Errorf("至少提供一个更新字段")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将更新会员 %s 基本信息\n", id)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.UpdateMember(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已更新会员 %s\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "会员 ID")
	cmd.Flags().StringVar(&username, "username", "", "姓名")
	cmd.Flags().StringVar(&nickname, "nickname", "", "昵称")
	cmd.Flags().StringVar(&birthday, "birthday", "", "生日，格式 yyyy-MM-dd")
	cmd.Flags().IntVar(&gender, "gender", 0, "性别，0=忽略 1=男 2=女")
	cmd.Flags().StringVar(&idNumber, "id-number", "", "身份证号")
	cmd.Flags().StringVar(&province, "province", "", "省份编码")
	cmd.Flags().StringVar(&city, "city", "", "城市编码")
	cmd.Flags().StringVar(&district, "district", "", "地区编码")
	cmd.Flags().StringVar(&country, "country", "", "国家编码")
	cmd.Flags().StringVar(&avatar, "avatar", "", "头像地址")
	cmd.Flags().StringVar(&qmFrom, "qm-from", "", "来源渠道")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdOpenID(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var infoType int

	cmd := &cobra.Command{
		Use:   "open-id",
		Short: "通过会员ID获取微信支付宝用户唯一标识",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerOpenID(cmd.Context(), customerID, infoType)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "OpenID", "UnionID"}, [][]string{{customerID, result.OpenID, result.UnionID}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().IntVar(&infoType, "type", 0, "0=支付宝unionId 1=微信openId 2=微信unionId")
	return cmd
}

func newCmdDynamicCode(f *cmdutil.Factory) *cobra.Command {
	var customerID int64

	cmd := &cobra.Command{
		Use:   "dynamic-code",
		Short: "查询会员的动态码",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == 0 {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetCustomerCode(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "会员码"}, [][]string{{strconv.FormatInt(customerID, 10), result.ClientCode}})
		},
	}

	cmd.Flags().Int64Var(&customerID, "customer-id", 0, "会员 ID")
	return cmd
}

func newCmdRegisterThirdID(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "register-third-id",
		Short: "通过三方ID注册会员",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeSingleRows(f, params, []string{"动作", "文件"}, [][]string{{"通过三方ID注册会员", fromJSON}})
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.RegisterByThirdID(cmd.Context(), client.ThirdIDRegisterParams(params))
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "手机号", "昵称", "姓名"}, [][]string{{result.CustomerID, result.Phone, result.Nickname, result.Username}})
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdSendCaptcha(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "send-captcha",
		Short: "发送短信验证码",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeSingleRows(f, params, []string{"动作", "文件"}, [][]string{{"发送短信验证码", fromJSON}})
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.SendSmsCaptcha(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已发送短信验证码: %s\n", fromJSON)
			return nil
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCheckCaptcha(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "check-captcha",
		Short: "校验短信验证码",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeSingleRows(f, params, []string{"动作", "文件"}, [][]string{{"校验短信验证码", fromJSON}})
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.CheckSmsCaptcha(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已校验短信验证码: %s\n", fromJSON)
			return nil
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdBlacklistAdd(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "blacklist-add",
		Short: "设置黑名单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeSingleRows(f, params, []string{"动作", "文件"}, [][]string{{"设置黑名单", fromJSON}})
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.BatchAddBlackList(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交设置黑名单: %s\n", fromJSON)
			return nil
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdBlacklistRemove(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "blacklist-remove",
		Short: "取消黑名单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			params, err := loadJSONFile(fromJSON)
			if err != nil {
				return err
			}
			if dryRun {
				return writeSingleRows(f, params, []string{"动作", "文件"}, [][]string{{"取消黑名单", fromJSON}})
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.BatchCancelBlackList(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交取消黑名单: %s\n", fromJSON)
			return nil
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "请求 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdBlacklistStatus(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "blacklist-status",
		Short: "黑名单校验",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryMemberBlack(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			label := "正常用户"
			if result.CustomerType == 1 {
				label = "黑名单用户"
			}
			return writeSingleRows(f, result, []string{"会员ID", "用户类型", "说明", "注册时间"}, [][]string{{customerID, strconv.Itoa(result.CustomerType), label, result.RegisterTime}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdFreeze(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var reason string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "freeze",
		Short: "冻结用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" || reason == "" {
				return fmt.Errorf("--customer-id 和 --reason 必填")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将冻结用户 %s\n", customerID)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.FreezeCustomer(cmd.Context(), customerID, reason); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已冻结用户: %s\n", customerID)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&reason, "reason", "", "冻结原因")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdUnfreeze(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "unfreeze",
		Short: "冻结用户-撤销",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将撤销冻结用户 %s\n", customerID)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.UnfreezeCustomer(cmd.Context(), customerID); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已撤销冻结用户: %s\n", customerID)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdLogoff(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var reason string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "logoff",
		Short: "会员注销",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" || reason == "" {
				return fmt.Errorf("--customer-id 和 --reason 必填")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将注销会员 %s\n", customerID)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.LogoffCustomer(cmd.Context(), customerID, reason); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已注销会员: %s\n", customerID)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().StringVar(&reason, "reason", "", "注销原因")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdUpdatePhone(f *cmdutil.Factory) *cobra.Command {
	var customerID int64
	var phone string
	var phoneEncrypt string
	var countryCode string
	var reason string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "update-phone",
		Short: "会员换绑手机号",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == 0 || reason == "" || (phone == "" && phoneEncrypt == "") {
				return fmt.Errorf("--customer-id、--reason 和 --phone/--phone-encrypt 必填")
			}
			params := map[string]interface{}{
				"customerId": customerID,
				"reason":     reason,
			}
			if phone != "" {
				params["phone"] = phone
			}
			if phoneEncrypt != "" {
				params["phoneEncrypt"] = phoneEncrypt
			}
			if countryCode != "" {
				params["countryCode"] = countryCode
			}
			if dryRun {
				return writeSingleRows(f, params, []string{"动作", "会员ID", "手机号", "国家区号"}, [][]string{{"会员换绑手机号", strconv.FormatInt(customerID, 10), phone, countryCode}})
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.UpdateCustomerPhone(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已更新会员 %d 手机号\n", customerID)
			return nil
		},
	}

	cmd.Flags().Int64Var(&customerID, "customer-id", 0, "会员 ID")
	cmd.Flags().StringVar(&phone, "phone", "", "手机号")
	cmd.Flags().StringVar(&phoneEncrypt, "phone-encrypt", "", "AES 加密手机号")
	cmd.Flags().StringVar(&countryCode, "country-code", "", "国家电话代码")
	cmd.Flags().StringVar(&reason, "reason", "", "换绑理由")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdFreezeRecord(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "freeze-record",
		Short: "查询用户冻结记录",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryFreezeRecord(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "冻结状态", "冻结结束时间"}, [][]string{{result.CustomerID, strconv.Itoa(result.FreezeStatus), result.FreezeEndTime}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdHasOpenOrder(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var bizType int

	cmd := &cobra.Command{
		Use:   "has-open-order",
		Short: "查询用户是否存在未完成订单",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" || bizType == 0 {
				return fmt.Errorf("--customer-id 和 --biz-type 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.CanLogoff(cmd.Context(), customerID, bizType)
			if err != nil {
				return err
			}
			label := "不存在未完成订单"
			if result.HasOrder {
				label = "存在未完成订单"
			}
			return writeSingleRows(f, result, []string{"会员ID", "是否有未完成订单", "说明"}, [][]string{{customerID, strconv.FormatBool(result.HasOrder), label}})
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().IntVar(&bizType, "biz-type", 0, "业务类型，5=新饮食")
	return cmd
}

func newCmdWeComInfo(f *cmdutil.Factory) *cobra.Command {
	var customerID int64
	cmd := &cobra.Command{
		Use:   "wecom-info",
		Short: "查询会员企微信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == 0 {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetWeComCustomerInfo(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"会员ID", "昵称", "员工数", "群聊数"}, [][]string{{
				result.CustomerID,
				result.Nickname,
				strconv.Itoa(len(result.BelongEmployees)),
				strconv.Itoa(len(result.BelongCustomerGroupList)),
			}})
		},
	}
	cmd.Flags().Int64Var(&customerID, "customer-id", 0, "会员 ID")
	return cmd
}

func newCmdSignStatus(f *cmdutil.Factory) *cobra.Command {
	var userID int64
	var activityID int64
	cmd := &cobra.Command{
		Use:   "sign-status",
		Short: "查询会员签到状态",
		RunE: func(cmd *cobra.Command, args []string) error {
			if userID == 0 {
				return fmt.Errorf("--user-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryActivitySign(cmd.Context(), userID, activityID)
			if err != nil {
				return err
			}
			label := "今日未签到"
			if result.TodaySign == 1 {
				label = "今日已签到"
			}
			return writeSingleRows(f, result, []string{"会员ID", "今日是否签到", "说明"}, [][]string{{strconv.FormatInt(userID, 10), strconv.Itoa(result.TodaySign), label}})
		},
	}
	cmd.Flags().Int64Var(&userID, "user-id", 0, "会员 ID")
	cmd.Flags().Int64Var(&activityID, "activity-id", 0, "签到活动 ID")
	return cmd
}

func newCmdAccountLevel(f *cmdutil.Factory) *cobra.Command {
	var customerID int64
	cmd := &cobra.Command{
		Use:   "account-level",
		Short: "查询会员储值等级",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == 0 {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryAccountLevel(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{strconv.FormatInt(item.CustomerID, 10), strconv.Itoa(item.Level)})
			}
			return writeSingleRows(f, result, []string{"会员ID", "储值等级"}, rows)
		},
	}
	cmd.Flags().Int64Var(&customerID, "customer-id", 0, "会员 ID")
	return cmd
}

func newCmdConditionQuery(f *cmdutil.Factory) *cobra.Command {
	var customerIDs []int64
	var conditions []int
	cmd := &cobra.Command{
		Use:   "condition-query",
		Short: "按需查询会员信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(customerIDs) == 0 || len(conditions) == 0 {
				return fmt.Errorf("--customer-ids 和 --conditions 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryCustomerCondition(cmd.Context(), customerIDs, conditions)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					item.ID,
					item.MobilePhone,
					item.Nickname,
					item.LevelName,
					fmt.Sprintf("%.2f", item.Balance),
					fmt.Sprintf("%.2f", item.PointNum),
					strconv.Itoa(item.CouponNum),
					strconv.Itoa(item.GiftCardNum),
				})
			}
			return writeSingleRows(f, result, []string{"会员ID", "手机号", "昵称", "等级", "余额", "积分", "优惠券数", "礼品卡数"}, rows)
		},
	}
	cmd.Flags().Int64SliceVar(&customerIDs, "customer-ids", nil, "会员 ID 列表，最多10个")
	cmd.Flags().IntSliceVar(&conditions, "conditions", nil, "查询范围列表，例如 10,20,40,60")
	return cmd
}

func newCmdTags(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "通过会员 ID 查询用户标签",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.PageTagsByCustomerID(cmd.Context(), customerID, pageNum, pageSize)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				name := item.Name
				if name == "" {
					name = item.TagName
				}
				rows = append(rows, []string{
					strconv.FormatInt(item.ID, 10),
					strconv.FormatInt(item.LabelID, 10),
					name,
					item.TypeName,
				})
			}
			if err := writeSingleRows(f, result, []string{"记录ID", "标签ID", "标签名称", "标签类型"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, pageNum, pageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdTagDetail(f *cmdutil.Factory) *cobra.Command {
	var id int64

	cmd := &cobra.Command{
		Use:   "detail",
		Short: "查询标签详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return fmt.Errorf("--id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetTagDetail(cmd.Context(), id)
			if err != nil {
				return err
			}
			return writeSingleRows(f, result, []string{"标签ID", "标签名称", "标签组ID", "标签组名称", "标签类型"}, [][]string{{
				strconv.FormatInt(result.ID, 10),
				result.Name,
				strconv.FormatInt(result.TagGroupID, 10),
				result.TagGroupName,
				result.TypeName,
			}})
		},
	}

	cmd.Flags().Int64Var(&id, "id", 0, "标签 ID")
	return cmd
}

func newCmdBrandTags(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brand-list",
		Short: "查询品牌下标签列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetBrandTags(cmd.Context())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				name := item.Name
				if name == "" {
					name = item.TagName
				}
				rows = append(rows, []string{
					strconv.FormatInt(item.ID, 10),
					strconv.FormatInt(item.LabelID, 10),
					name,
					item.TypeName,
				})
			}
			return writeSingleRows(f, result, []string{"记录ID", "标签ID", "标签名称", "标签类型"}, rows)
		},
	}

	return cmd
}

func newCmdTagIDs(f *cmdutil.Factory) *cobra.Command {
	var customerID string

	cmd := &cobra.Command{
		Use:   "ids",
		Short: "查询会员的标签 ID 列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" {
				return fmt.Errorf("--customer-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetMemberTagIDs(cmd.Context(), customerID)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.LabelIDs))
			for _, id := range result.LabelIDs {
				rows = append(rows, []string{customerID, strconv.FormatInt(id, 10)})
			}
			return writeSingleRows(f, result, []string{"会员ID", "标签ID"}, rows)
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	return cmd
}

func newCmdGroupList(f *cmdutil.Factory) *cobra.Command {
	var labelAttributed int
	var labelGroupName string
	var filterEmpty int

	cmd := &cobra.Command{
		Use:   "groups",
		Short: "根据标签类别查询静态标签组列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if labelAttributed == 0 {
				return fmt.Errorf("--label-attributed 必填，1=系统标签 2=自定义标签")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListStaticLabelGroups(cmd.Context(), labelAttributed, labelGroupName, filterEmpty)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(*result))
			for _, item := range *result {
				rows = append(rows, []string{
					strconv.FormatInt(item.ID, 10),
					strconv.Itoa(item.LabelAttributed),
					item.LabelGroupName,
					strconv.Itoa(item.LabelGroupType),
					item.OperateName,
				})
			}
			return writeSingleRows(f, result, []string{"标签组ID", "标签归属", "标签组名称", "标签组类型", "操作人"}, rows)
		},
	}

	cmd.Flags().IntVar(&labelAttributed, "label-attributed", 0, "标签类别，1=系统标签 2=自定义标签")
	cmd.Flags().StringVar(&labelGroupName, "label-group-name", "", "标签组名称")
	cmd.Flags().IntVar(&filterEmpty, "filter-empty", 0, "过滤空标签组，1=过滤")
	return cmd
}

func newCmdGroupLabels(f *cmdutil.Factory) *cobra.Command {
	var labelGroupID int64
	var labelName string
	var pageNum int
	var pageSize int

	cmd := &cobra.Command{
		Use:   "group-labels",
		Short: "根据标签组查询静态用户标签列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if labelGroupID == 0 {
				return fmt.Errorf("--label-group-id 必填")
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListStaticLabelsByGroup(cmd.Context(), labelGroupID, labelName, pageNum, pageSize)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					strconv.FormatInt(item.ID, 10),
					item.LabelCode,
					item.LabelName,
					strconv.Itoa(item.LabelStatus),
					strconv.FormatInt(item.LabelNum, 10),
				})
			}
			if err := writeSingleRows(f, result, []string{"标签ID", "标签编码", "标签名称", "标签状态", "标签人数"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, result.PageNum, result.PageSize)
			return nil
		},
	}

	cmd.Flags().Int64Var(&labelGroupID, "label-group-id", 0, "标签组 ID")
	cmd.Flags().StringVar(&labelName, "label-name", "", "标签名称")
	cmd.Flags().IntVar(&pageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "每页条数")
	return cmd
}

func newCmdMark(f *cmdutil.Factory) *cobra.Command {
	var customerIDs []string
	var thirdMemberIDs []string
	var labelCode string
	var markDate string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "mark",
		Short: "给指定会员打标签",
		RunE: func(cmd *cobra.Command, args []string) error {
			if labelCode == "" || markDate == "" {
				return fmt.Errorf("--label-code 和 --mark-date 必填")
			}
			if len(customerIDs) == 0 && len(thirdMemberIDs) == 0 {
				return fmt.Errorf("--customer-ids 或 --third-member-ids 至少一个必填")
			}
			params := client.TagMarkParams{
				CustomerIDList: customerIDs,
				LabelCode:      labelCode,
				MarkDate:       markDate,
				ThirdMemberID:  thirdMemberIDs,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将打标签 %s 到会员 %v\n", labelCode, customerIDs)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.MarkPanoramaLabel(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交会员打标: labelCode=%s\n", labelCode)
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&customerIDs, "customer-ids", nil, "企迈会员 ID 列表")
	cmd.Flags().StringSliceVar(&thirdMemberIDs, "third-member-ids", nil, "三方会员 ID 列表")
	cmd.Flags().StringVar(&labelCode, "label-code", "", "标签编码")
	cmd.Flags().StringVar(&markDate, "mark-date", "", "打标时间，格式 yyyy-MM-dd HH:ss:mm")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdCreateTag(f *cmdutil.Factory) *cobra.Command {
	var labelCode, labelName string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建会员标签",
		RunE: func(cmd *cobra.Command, args []string) error {
			if labelCode == "" || labelName == "" {
				return fmt.Errorf("--label-code 和 --label-name 必填")
			}
			params := client.LabelCreateParams{LabelCode: labelCode, LabelName: labelName}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将创建会员标签 %s(%s)\n", labelName, labelCode)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.CreatePanoramaLabel(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已创建会员标签: %s\n", labelCode)
			return nil
		},
	}
	cmd.Flags().StringVar(&labelCode, "label-code", "", "标签编码")
	cmd.Flags().StringVar(&labelName, "label-name", "", "标签名称")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdDeleteTag(f *cmdutil.Factory) *cobra.Command {
	var labelCode string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "删除会员标签",
		RunE: func(cmd *cobra.Command, args []string) error {
			if labelCode == "" {
				return fmt.Errorf("--label-code 必填")
			}
			params := client.LabelDeleteParams{LabelCode: labelCode}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将删除会员标签 %s\n", labelCode)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.DeletePanoramaLabel(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已删除会员标签: %s\n", labelCode)
			return nil
		},
	}
	cmd.Flags().StringVar(&labelCode, "label-code", "", "标签编码")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdClearTagMembers(f *cmdutil.Factory) *cobra.Command {
	var customerIDs []string
	var thirdMemberIDs []string
	var labelCode string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "clear-members",
		Short: "清除会员标签下的会员",
		RunE: func(cmd *cobra.Command, args []string) error {
			if labelCode == "" {
				return fmt.Errorf("--label-code 必填")
			}
			if len(customerIDs) == 0 && len(thirdMemberIDs) == 0 {
				return fmt.Errorf("--customer-ids 或 --third-member-ids 至少一个必填")
			}
			params := client.LabelClearMembersParams{
				CustomerIDList: customerIDs,
				LabelCode:      labelCode,
				ThirdMemberID:  thirdMemberIDs,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将清除标签 %s 下的会员\n", labelCode)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.ClearPanoramaLabelMembers(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已清除标签会员: %s\n", labelCode)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&customerIDs, "customer-ids", nil, "企迈会员 ID 列表")
	cmd.Flags().StringSliceVar(&thirdMemberIDs, "third-member-ids", nil, "三方会员 ID 列表")
	cmd.Flags().StringVar(&labelCode, "label-code", "", "标签编码")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdDeleteCustomerTag(f *cmdutil.Factory) *cobra.Command {
	var customerID string
	var panoramaLabelID int64
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "delete-customer-label",
		Short: "清除会员的指定标签",
		RunE: func(cmd *cobra.Command, args []string) error {
			if customerID == "" || panoramaLabelID == 0 {
				return fmt.Errorf("--customer-id 和 --panorama-label-id 必填")
			}
			params := client.DeleteCustomerLabelParams{
				CustomerID:      customerID,
				PanoramaLabelID: panoramaLabelID,
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将清除会员 %s 的标签 %d\n", customerID, panoramaLabelID)
				return nil
			}
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			if err := api.DeleteCustomerLabel(cmd.Context(), params); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已清除会员标签: customerId=%s labelId=%d\n", customerID, panoramaLabelID)
			return nil
		},
	}
	cmd.Flags().StringVar(&customerID, "customer-id", "", "会员 ID")
	cmd.Flags().Int64Var(&panoramaLabelID, "panorama-label-id", 0, "标签 ID")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdTagSettings(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "查询标签设置",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newMemberAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetBrandInfo(cmd.Context())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.LabelList))
			for _, item := range result.LabelList {
				rows = append(rows, []string{
					strconv.FormatInt(item.LabelID, 10),
					item.LabelCode,
					item.LabelName,
					strconv.Itoa(item.LabelLevel),
				})
			}
			return writeSingleRows(f, result, []string{"标签ID", "标签编码", "标签名称", "标签层级"}, rows)
		},
	}
	return cmd
}

func newJSONAssetMutationCmd(f *cmdutil.Factory, use, short, action string, invoke func(api *client.MemberAPI, cmd *cobra.Command, params map[string]interface{}) error) *cobra.Command {
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
				return writeSingleRows(f, params, []string{"动作", "文件"}, [][]string{{action, fromJSON}})
			}
			api, err := newMemberAPI(f)
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

func loadJSONFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 JSON 文件失败: %w", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("解析 JSON 文件失败: %w", err)
	}
	return payload, nil
}

func writeSingleRows(f *cmdutil.Factory, data interface{}, headers []string, rows [][]string) error {
	format, err := output.ParseFormat(f.EffectiveFormat())
	if err != nil {
		return err
	}
	fmtr := output.NewFormatter(f.IOStreams.Out, format)
	return fmtr.Write(data, headers, rows)
}
