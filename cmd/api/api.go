package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/madaima/qmai-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAPI creates the api command for raw API passthrough
func NewCmdAPI(f *cmdutil.Factory) *cobra.Command {
	var body string

	cmd := &cobra.Command{
		Use:   "api <path>",
		Short: "Raw API 请求透传",
		Long: `直接发送开放平台 API 请求（自动签名）。

示例:
  qmai api v3/goods/item/getItemList --body '{"shopCode":"S001"}'
  qmai api v3/goods/item/getItemList --body @params.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			apiClient, err := f.ApiClient()
			if err != nil {
				return err
			}

			// Parse params from body
			var params interface{}
			if body != "" {
				var bodyStr string
				if strings.HasPrefix(body, "@") {
					data, err := os.ReadFile(body[1:])
					if err != nil {
						return fmt.Errorf("读取文件失败: %w", err)
					}
					bodyStr = string(data)
				} else {
					bodyStr = body
				}
				if err := json.Unmarshal([]byte(bodyStr), &params); err != nil {
					return fmt.Errorf("解析 JSON 失败: %w", err)
				}
			}

			resp, err := apiClient.Call(cmd.Context(), path, params)
			if err != nil {
				return fmt.Errorf("请求失败: %w", err)
			}

			// Pretty-print the response data
			enc := json.NewEncoder(f.IOStreams.Out)
			enc.SetIndent("", "  ")
			if resp.Data != nil {
				var data interface{}
				if err := json.Unmarshal(resp.Data, &data); err == nil {
					enc.Encode(data)
				} else {
					fmt.Fprintln(f.IOStreams.Out, string(resp.Data))
				}
			} else {
				fmt.Fprintln(f.IOStreams.Out, resp.Message)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "业务参数 JSON 或 @filename")

	return cmd
}
