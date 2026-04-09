package product

import (
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdProduct(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "商品管理",
		Long:  "商品列表、创建、更新、上下架、批量导入导出、批量调价",
	}

	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdCreate(f))
	cmd.AddCommand(newCmdUpdate(f))
	cmd.AddCommand(newCmdDelete(f))
	cmd.AddCommand(newCmdCategory(f))
	cmd.AddCommand(newCmdImport(f))
	cmd.AddCommand(newCmdExport(f))
	cmd.AddCommand(newCmdBatchPrice(f))
	cmd.AddCommand(newCmdBatchStatus(f))
	cmd.AddCommand(newCmdEstimateClear(f))
	cmd.AddCommand(newCmdCancelEstimateClear(f))
	cmd.AddCommand(newCmdSoldOut(f))
	cmd.AddCommand(newCmdFillFull(f))
	cmd.AddCommand(newCmdPracticeStatus(f))
	cmd.AddCommand(newCmdDeleteTask(f))
	cmd.AddCommand(newCmdAttachList(f))
	cmd.AddCommand(newCmdListWithPractice(f))
	cmd.AddCommand(newCmdEnergy(f))
	cmd.AddCommand(newCmdRealtime(f))

	return cmd
}
