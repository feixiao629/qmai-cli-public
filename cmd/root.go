package cmd

import (
	"github.com/madaima/qmai-cli/cmd/api"
	"github.com/madaima/qmai-cli/cmd/auth"
	"github.com/madaima/qmai-cli/cmd/completion"
	configcmd "github.com/madaima/qmai-cli/cmd/config"
	"github.com/madaima/qmai-cli/cmd/delivery"
	"github.com/madaima/qmai-cli/cmd/doctor"
	"github.com/madaima/qmai-cli/cmd/finance"
	"github.com/madaima/qmai-cli/cmd/inventory"
	"github.com/madaima/qmai-cli/cmd/marketing"
	"github.com/madaima/qmai-cli/cmd/member"
	"github.com/madaima/qmai-cli/cmd/order"
	"github.com/madaima/qmai-cli/cmd/product"
	"github.com/madaima/qmai-cli/cmd/queue"
	"github.com/madaima/qmai-cli/cmd/shortcuts"
	"github.com/madaima/qmai-cli/cmd/store"
	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd(version string) *cobra.Command {
	f := cmdutil.NewFactory()

	cmd := &cobra.Command{
		Use:           "qmai",
		Short:         "门店经营 CLI 工具",
		Long:          "qmai — 门店经营命令行工具，覆盖商品管理、门店与组织管理、会员服务、营销服务、订单服务、财务服务、聚合配送、进销存、排队服务",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	cmd.PersistentFlags().StringVar(&f.Format, "format", "", "Output format: json, table, csv")
	cmd.PersistentFlags().StringVar(&f.Profile, "profile", "", "Use a specific store profile")
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Enable debug output")

	// Core commands
	productCmd := product.NewCmdProduct(f)
	storeCmd := store.NewCmdStore(f)
	memberCmd := member.NewCmdMember(f)
	marketingCmd := marketing.NewCmdMarketing(f)
	orderCmd := order.NewCmdOrder(f)
	financeCmd := finance.NewCmdFinance(f)
	deliveryCmd := delivery.NewCmdDelivery(f)
	inventoryCmd := inventory.NewCmdInventory(f)
	queueCmd := queue.NewCmdQueue(f)

	cmd.AddCommand(newVersionCmd(version))
	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(configcmd.NewCmdConfig(f))
	cmd.AddCommand(productCmd)
	cmd.AddCommand(storeCmd)
	cmd.AddCommand(memberCmd)
	cmd.AddCommand(marketingCmd)
	cmd.AddCommand(orderCmd)
	cmd.AddCommand(financeCmd)
	cmd.AddCommand(deliveryCmd)
	cmd.AddCommand(inventoryCmd)
	cmd.AddCommand(queueCmd)
	cmd.AddCommand(api.NewCmdAPI(f))
	cmd.AddCommand(doctor.NewCmdDoctor(f))
	cmd.AddCommand(completion.NewCmdCompletion())

	// Register shortcuts under domain commands
	shortcuts.RegisterShortcuts(productCmd, f)

	return cmd
}

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("qmai version %s\n", version)
		},
	}
}
