package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/miaojuncn/oss-cli/store"
	"github.com/miaojuncn/oss-cli/store/aliyun"
	"github.com/spf13/cobra"
)

var (
	fileName   string
	bucketName string
)

// uploadCmd represents the start command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件到中转站",
	Long:  `上传文件到中转站`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := getAccessKeyFromInputV2()
		if err != nil {
			return err
		}
		p, err := getProvider()
		if err != nil {
			return err
		}
		if fileName == "" {
			return fmt.Errorf("file name required")
		}

		// 为了防止文件都堆在一个文件夹里面 无法查看
		// 我们采用日期进行编码
		day := time.Now().Format("20060102")

		// 为了防止不同用户同一时间上传相同的文件
		// 我们采用用户的主机名作为前置
		hn, err := os.Hostname()
		if err != nil {
			ipAddr := getOutBindIp()
			if ipAddr == "" {
				hn = "unknown"
			} else {
				hn = ipAddr
			}
		}

		objectKey := fmt.Sprintf("%s/%s/%s", day, hn, fileName)
		downloadUrl, err := p.Upload(bucketName, objectKey, fileName)
		if err != nil {
			return err
		}

		// 正常退出
		fmt.Printf("文  件: %s 上传成功\n", fileName)

		// 打印下载链接
		fmt.Printf("下载链接: %s\n", downloadUrl)
		fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")

		return nil
	},
}

func getAccessKeyFromInputV2() error {
	prompt := &survey.Password{
		Message: "请输入 accessSecret: ",
	}
	if err := survey.AskOne(prompt, &accessSecret); err != nil {
		return err
	} else {
		// fmt.Println()
		return nil
	}
}

func getOutBindIp() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	addr := strings.Split(conn.LocalAddr().String(), ":")
	if len(addr) == 0 {
		return ""
	}

	return addr[0]
}

func getProvider() (p store.Uploader, err error) {
	switch ossProvider {
	case "aliyun":
		return aliyun.NewAliOssStore(&aliyun.Options{
			Endpoint:     ossEndpoint,
			AccessKey:    accessKey,
			AccessSecret: accessSecret,
		})
	case "qcloud":
		return nil, fmt.Errorf("not impl")
	default:
		return nil, fmt.Errorf("unknown oss privier options [aliyun/qcloud]")
	}
}

func init() {
	uploadCmd.PersistentFlags().StringVarP(&bucketName, "bucket_name", "b", "appchain-production", "上传文件的名称")
	uploadCmd.PersistentFlags().StringVarP(&fileName, "file_name", "f", "", "上传文件的名称")
	RootCmd.AddCommand(uploadCmd)
}
