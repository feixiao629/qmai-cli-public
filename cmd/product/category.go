package product

import (
	"fmt"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdCategory(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category",
		Short: "分类管理",
		Long:  "商品分类管理（开放平台暂未提供独立分类 API，分类通过商品同步时 className 自动创建）",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(f.IOStreams.Out, "分类管理功能尚未适配开放平台。")
			fmt.Fprintln(f.IOStreams.Out, "提示: 通过 'qmai product create --class 分类名' 创建商品时会自动关联分类。")
			return nil
		},
	}

	return cmd
}
