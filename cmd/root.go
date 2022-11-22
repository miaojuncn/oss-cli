package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ossProvider  string
	version      bool
	ossEndpoint  string
	accessKey    string
	accessSecret string
)

var RootCmd = &cobra.Command{
	Use:   "oss-cli",
	Short: "oss-cli short",
	Long:  "oss-cli long",
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println("cloud-station-cli 0.0.1")
		}
		return nil
	},
}

func init() {
	f := RootCmd.PersistentFlags()
	f.BoolVarP(&version, "version", "v", false, "oss-cli version")
	f.StringVarP(&ossProvider, "provider", "p", "aliyun", "oss provider")
	f.StringVarP(&ossEndpoint, "endpoint", "e", "", "oss endpoint")
	f.StringVarP(&accessKey, "key", "k", "", "oss key")
}
