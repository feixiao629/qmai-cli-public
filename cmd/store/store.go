package store

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/madaima/qmai-cli/internal/client"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/madaima/qmai-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdStore creates the store command group for org service APIs.
func NewCmdStore(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "门店与组织管理",
		Long:  "门店详情、门店列表、组织树、门店组、配置查询、门店状态与门店同步",
	}

	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdGetByID(f))
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdTakeoutMapList(f))
	cmd.AddCommand(newCmdID(f))
	cmd.AddCommand(newCmdSetStatus(f))
	cmd.AddCommand(newCmdSync(f))
	cmd.AddCommand(newCmdMoveTeam(f))
	cmd.AddCommand(newCmdConfig(f))
	cmd.AddCommand(newCmdBrandConfig(f))
	cmd.AddCommand(newCmdExtData(f))
	cmd.AddCommand(newCmdLabels(f))
	cmd.AddCommand(newCmdOrgTree(f))
	cmd.AddCommand(newCmdTeamList(f))
	cmd.AddCommand(newCmdDeptTree(f))

	return cmd
}

func newStoreAPI(f *cmdutil.Factory) (*client.StoreAPI, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return nil, err
	}
	return client.NewStoreAPI(apiClient), nil
}

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <shop-code>",
		Short: "根据门店编码查询门店详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			detail, err := api.GetShopDetailByCode(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return writeShopDetail(f, detail)
		},
	}
}

func newCmdGetByID(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get-by-id <shop-id>",
		Short: "根据门店 ID 查询门店详情",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shopID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("门店 ID 格式错误: %s", args[0])
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			detail, err := api.GetShopDetailByID(cmd.Context(), shopID)
			if err != nil {
				return err
			}
			return writeShopDetail(f, detail)
		},
	}
}

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	var params client.ShopListParams

	cmd := &cobra.Command{
		Use:   "list",
		Short: "批量查询门店",
		RunE: func(cmd *cobra.Command, args []string) error {
			if params.PageNum == 0 {
				params.PageNum = 1
			}
			if params.PageSize == 0 {
				params.PageSize = 10
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListShops(cmd.Context(), params)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					fmt.Sprintf("%d", item.ID),
					item.Code,
					item.Name,
					item.OperateStatus,
					item.ContactPhone,
					item.FullAddress,
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			if err := fmtr.Write(result, []string{"ID", "编码", "名称", "营业状态", "联系电话", "地址"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, params.PageNum, params.PageSize)
			return nil
		},
	}

	cmd.Flags().StringVar(&params.Keyfield, "keyfield", "", "地址/联系电话")
	cmd.Flags().StringVar(&params.Keyword, "keyword", "", "门店名称/编码/ID")
	cmd.Flags().Int64Var(&params.LabelID, "label-id", 0, "门店标签 ID")
	cmd.Flags().IntVar(&params.PageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&params.PageSize, "page-size", 10, "每页条数")
	cmd.Flags().StringVar(&params.Search, "search", "", "搜索关键字，匹配门店名称或地址")
	cmd.Flags().IntVar(&params.Type, "type", 0, "机构类型，1=品牌 2=机构 3=门店")
	cmd.Flags().Int64Var(&params.TypeID, "type-id", 0, "机构类型 ID")
	cmd.Flags().IntVar(&params.ContainCloseFlag, "contain-close", 0, "是否包含停业门店，1=包含")

	return cmd
}

func newCmdTakeoutMapList(f *cmdutil.Factory) *cobra.Command {
	var params client.TakeoutShopMappingParams

	cmd := &cobra.Command{
		Use:   "takeout-map-list",
		Short: "查询平台外卖门店映射明细",
		RunE: func(cmd *cobra.Command, args []string) error {
			if params.PlatformType == 0 {
				return fmt.Errorf("--platform-type 必填，4=美团 5=饿了么")
			}
			if params.PageNum == 0 {
				params.PageNum = 1
			}
			if params.PageSize == 0 {
				params.PageSize = 10
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListTakeoutMappings(cmd.Context(), params)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.List))
			for _, item := range result.List {
				rows = append(rows, []string{
					item.ShopCode,
					fmt.Sprintf("%d", item.ShopID),
					fmt.Sprintf("%d", item.PlatformType),
					item.ExternalShopID,
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			if err := fmtr.Write(result, []string{"门店编码", "门店ID", "平台类型", "三方门店ID"}, rows); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "\n共 %d 条 (第 %d 页，每页 %d 条)\n", result.Total, params.PageNum, params.PageSize)
			return nil
		},
	}

	cmd.Flags().IntVar(&params.PlatformType, "platform-type", 0, "平台类型，4=美团 5=饿了么")
	cmd.Flags().IntVar(&params.PageNum, "page", 1, "页码")
	cmd.Flags().IntVar(&params.PageSize, "page-size", 10, "每页条数")
	return cmd
}

func newCmdID(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "id <shop-code>",
		Short: "根据门店编码查询门店 ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ShopCodeToID(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"门店编码", "门店ID"}, [][]string{{args[0], fmt.Sprintf("%d", result.ID)}})
		},
	}
}

func newCmdSetStatus(f *cmdutil.Factory) *cobra.Command {
	var status int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "set-status <shop-code>",
		Short: "设置门店营业状态",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if status != 1 && status != 2 {
				return fmt.Errorf("--status 必须为 1(开启) 或 2(关闭)")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将设置门店 %s 营业状态为 %d\n", args[0], status)
				return nil
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			if err := api.ChangeShopStatus(cmd.Context(), args[0], status); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已更新门店 %s 营业状态为 %d\n", args[0], status)
			return nil
		},
	}

	cmd.Flags().IntVar(&status, "status", 0, "营业状态，1=开启 2=关闭")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdSync(f *cmdutil.Factory) *cobra.Command {
	var fromJSON string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "同步门店信息",
		Long:  "通过 JSON 文件提交门店信息同步，适合复杂字段结构和批量准备。",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromJSON == "" {
				return fmt.Errorf("--from-json 必填")
			}
			data, err := os.ReadFile(fromJSON)
			if err != nil {
				return fmt.Errorf("读取 JSON 文件失败: %w", err)
			}
			var payload client.ShopSyncRequest
			if err := json.Unmarshal(data, &payload); err != nil {
				return fmt.Errorf("解析 JSON 失败: %w", err)
			}
			if payload.Name == "" || payload.ProvinceID == 0 || payload.CityID == 0 || payload.DistrictID == 0 || payload.Lat == "" || payload.Lng == "" || payload.ManagerStatus == 0 {
				return fmt.Errorf("同步门店信息缺少核心必填字段，至少需包含 name/provinceId/cityId/districtId/lat/lng/managerStatus")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将同步门店信息: name=%s code=%s\n", payload.Name, payload.Code)
				return nil
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			if err := api.SyncShopInfo(cmd.Context(), payload); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 已提交门店同步: %s\n", payload.Name)
			return nil
		},
	}

	cmd.Flags().StringVar(&fromJSON, "from-json", "", "门店同步 JSON 文件")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdMoveTeam(f *cmdutil.Factory) *cobra.Command {
	var shopID int64
	var teamID int64
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "move-team",
		Short: "修改门店组别",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopID == 0 || teamID == 0 {
				return fmt.Errorf("--shop-id 和 --team-id 必填")
			}
			if dryRun {
				fmt.Fprintf(f.IOStreams.Out, "[dry-run] 将把门店 %d 调整到门店组 %d\n", shopID, teamID)
				return nil
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ModifyShopTeam(cmd.Context(), shopID, teamID)
			if err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "✓ 门店组别变更结果: %v\n", result.Data)
			return nil
		},
	}

	cmd.Flags().Int64Var(&shopID, "shop-id", 0, "门店 ID")
	cmd.Flags().Int64Var(&teamID, "team-id", 0, "门店组 ID")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览，不实际执行")
	return cmd
}

func newCmdConfig(f *cmdutil.Factory) *cobra.Command {
	var storeID int64
	var fieldCodes []string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "查询门店配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			if storeID == 0 {
				return fmt.Errorf("--store-id 必填")
			}
			if len(fieldCodes) == 0 {
				return fmt.Errorf("--field-codes 必填")
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QueryStoreConfigBatch(cmd.Context(), storeID, fieldCodes)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.ConfigList))
			for _, item := range result.ConfigList {
				rows = append(rows, []string{fmt.Sprintf("%d", result.StoreID), item.FieldCode, item.FieldValue})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"门店ID", "配置编码", "配置值"}, rows)
		},
	}

	cmd.Flags().Int64Var(&storeID, "store-id", 0, "门店 ID")
	cmd.Flags().StringSliceVar(&fieldCodes, "field-codes", nil, "配置编码列表")
	return cmd
}

func newCmdBrandConfig(f *cmdutil.Factory) *cobra.Command {
	var fieldCodes []string

	cmd := &cobra.Command{
		Use:   "brand-config",
		Short: "查询品牌配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(fieldCodes) == 0 {
				return fmt.Errorf("--field-codes 必填")
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.QuerySellerConfigBatch(cmd.Context(), fieldCodes)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result.ConfigList))
			for _, item := range result.ConfigList {
				rows = append(rows, []string{item.FieldCode, item.FieldValue})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"配置编码", "配置值"}, rows)
		},
	}

	cmd.Flags().StringSliceVar(&fieldCodes, "field-codes", nil, "配置编码列表")
	return cmd
}

func newCmdExtData(f *cmdutil.Factory) *cobra.Command {
	var shopID int64

	cmd := &cobra.Command{
		Use:   "ext-data",
		Short: "查询门店扩展字段信息",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopID == 0 {
				return fmt.Errorf("--shop-id 必填")
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetShopExtData(cmd.Context(), shopID)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result))
			for _, item := range result {
				rows = append(rows, []string{fmt.Sprintf("%d", item.ID), item.Name, item.Value})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"扩展ID", "字段名", "字段值"}, rows)
		},
	}

	cmd.Flags().Int64Var(&shopID, "shop-id", 0, "门店 ID")
	return cmd
}

func newCmdLabels(f *cmdutil.Factory) *cobra.Command {
	var shopID int64

	cmd := &cobra.Command{
		Use:   "labels",
		Short: "查询门店标签列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shopID == 0 {
				return fmt.Errorf("--shop-id 必填")
			}
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListShopLabels(cmd.Context(), shopID)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(result))
			for _, item := range result {
				rows = append(rows, []string{fmt.Sprintf("%d", item.ID), item.Name})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"标签ID", "标签名"}, rows)
		},
	}

	cmd.Flags().Int64Var(&shopID, "shop-id", 0, "门店 ID")
	return cmd
}

func newCmdOrgTree(f *cmdutil.Factory) *cobra.Command {
	var containClose int

	cmd := &cobra.Command{
		Use:   "org-tree",
		Short: "查询组织机构列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetOrgTree(cmd.Context(), containClose)
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
					fmt.Sprintf("%d", item.BrandID),
					item.BrandName,
					fmt.Sprintf("%d", item.UnBoundShopCount),
				})
			}
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"品牌ID", "品牌名", "未绑定门店数"}, rows)
		},
	}

	cmd.Flags().IntVar(&containClose, "contain-close", 0, "是否包含已停业门店，1=包含")
	return cmd
}

func newCmdTeamList(f *cmdutil.Factory) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "team-list",
		Short: "查询门店组列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.ListShopTeams(cmd.Context(), name)
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0)
			appendShopTeamRows(&rows, result, 0)
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"门店组ID", "名称", "层级", "门店数量", "管理员", "路径"}, rows)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "门店组名称")
	return cmd
}

func newCmdDeptTree(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "dept-tree",
		Short: "查询门店树结构数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := newStoreAPI(f)
			if err != nil {
				return err
			}
			result, err := api.GetShopDeptTree(cmd.Context())
			if err != nil {
				return err
			}
			format, err := output.ParseFormat(f.EffectiveFormat())
			if err != nil {
				return err
			}
			rows := make([][]string, 0)
			appendDeptRows(&rows, *result, 0)
			fmtr := output.NewFormatter(f.IOStreams.Out, format)
			return fmtr.Write(result, []string{"部门ID", "名称", "层级", "负责人", "类型", "路径"}, rows)
		},
	}
}

func writeShopDetail(f *cmdutil.Factory, detail *client.ShopDetail) error {
	format, err := output.ParseFormat(f.EffectiveFormat())
	if err != nil {
		return err
	}
	fmtr := output.NewFormatter(f.IOStreams.Out, format)
	return fmtr.Write(detail, []string{"ID", "编码", "名称", "营业状态", "联系电话", "地址"}, [][]string{{
		fmt.Sprintf("%d", detail.ID),
		detail.Code,
		detail.Name,
		detail.OperateStatus,
		detail.ContactPhone,
		firstNonEmpty(detail.FullAddress, detail.Address),
	}})
}

func appendShopTeamRows(rows *[][]string, teams []client.ShopTeam, level int) {
	for _, item := range teams {
		*rows = append(*rows, []string{
			fmt.Sprintf("%d", item.ID),
			item.Name,
			strconv.Itoa(level),
			fmt.Sprintf("%d", item.Num),
			item.ManagerAccountName,
			item.Path,
		})
		if len(item.Children) > 0 {
			appendShopTeamRows(rows, item.Children, level+1)
		}
	}
}

func appendDeptRows(rows *[][]string, node client.ShopDeptNode, level int) {
	*rows = append(*rows, []string{
		fmt.Sprintf("%d", node.ID),
		node.Name,
		strconv.Itoa(level),
		node.LeaderName,
		fmt.Sprintf("%d", node.Type),
		node.Path,
	})
	for _, child := range node.SubDeptTree {
		appendDeptRows(rows, child, level+1)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
