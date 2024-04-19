package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)
var accessKeyId string
var accessKeySecret string
var portRange string
var regionId string
func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&accessKeyId, "accessKeyId", "i", "", "your accessKeyId")
	addCmd.Flags().StringVarP(&accessKeySecret, "accessKeySecret", "s", "", "your accessKeySecret")
	addCmd.Flags().StringVarP(&portRange, "portRange", "p", "1/65535", "The port range is delineated by a slash (/), separating the starting port from the ending port.")
	addCmd.Flags().StringVarP(&regionId, "regionId", "r", "cn-shanghai", "The ID of the region where the security group resides.")
	addCmd.MarkFlagRequired("accessKeyId")
	addCmd.MarkFlagRequired("accessKeySecret")
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a security group rule",
	Long: `accessKeyId, accessSecret:前往 https://ram.console.aliyun.com/manage/ak 添加 accessKey
regionId:安全组所属地域ID ,比如cn-shanghai
	 访问 [DescribeRegions:查询可以使用的阿里云地域](https://next.api.aliyun.com/api/Ecs/2014-05-26/DescribeRegions) 查阅
	 国内一般是去掉 ECS 所在可用区的后缀，比如去掉 cn-guangzhou-b 的尾号 -b`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		// fmt.Println(accessKeyId, accessKeySecret, portRange)
		add(accessKeyId,accessKeySecret,portRange,regionId)
	},
}

var ipSet map[string]struct{}

func fetch(url string, wg *sync.WaitGroup) {
	// 确保在函数退出时调用 Done 来通知 main 函数工作已经完成
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		if strings.Contains(url, "ipv6") {
			fmt.Println("It seems you don't have ipv6 enabled.")
		}
		// panic(err)
		return
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		// panic(err)
		return
	}

	ipSet[strings.TrimSpace(string(ip))] = struct{}{}
	// fmt.Printf("My IP is:%s\n", ip)

}

func addSecurityGroupRules(client *ecs.Client, protocol string, clientIP string, PortRange string) {
	fmt.Printf("protocol: %s, clientIP: %s, PortRange: %s\n", protocol, clientIP, PortRange)
	request := ecs.CreateAuthorizeSecurityGroupRequest()
	request.Scheme = "https"

	request.SecurityGroupId = "sg-uf6bdzcrtiueqk21o907" // 安全组ID
	request.IpProtocol = "tcp"                          // 协议,可选 tcp,udp, icmp, gre, all：支持所有协议
	request.PortRange = PortRange                       // 端口范围，使用斜线（/）隔开起始端口和终止端口
	request.Priority = "1"                              // 安全组规则优先级，数字越小，代表优先级越高。取值范围：1~100
	request.Policy = "accept"                           // accept:接受访问, drop: 拒绝访问
	request.NicType = "internet"                        // internet：公网网卡, intranet：内网网卡。
	request.SourceCidrIp = string(clientIP)             // 源端IPv4 CIDR地址段。支持CIDR格式和IPv4格式的IP地址范围。

	response, err := client.AuthorizeSecurityGroup(request)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Printf("Response: %#v\nClient IP: %s  was successfully added to the Security Group.\n", response, clientIP)
}

func add(accessKeyId string, accessKeySecret string,portRange string,regionId string) {
	// url := "http://ipv4.icanhazip.com/" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	ipSet = make(map[string]struct{})
	fmt.Printf("Getting IP address from  ipify ...\n")
	var wg sync.WaitGroup

	var ipArr = [...]string{"https://4.ipw.cn/", "http://ipv4.whatismyip.akamai.com/", "http://ipv6.whatismyip.akamai.com/", "http://ipv4.icanhazip.com/",
		"http://ipv6.icanhazip.com/"}

	for _, url := range ipArr {
		// 每次循环调用 Add 来增加等待的 goroutine 的数量
		wg.Add(1)

		go fetch(url, &wg)
	}

	// 等待所有的 goroutine 完成
	wg.Wait()

	// fmt.Printf("Total number of unique IP addresses: %d\n", len(ipSet))
	fmt.Print("Unique IP addresses: ")
	for key, _ := range ipSet {
		fmt.Printf("%s ", key)
	}
	fmt.Printf("\n")
	// <accessKeyId>, <accessSecret>: 前往 https://ram.console.aliyun.com/manage/ak 添加 accessKey
	// RegionId：安全组所属地域ID ，比如 `cn-guangzhou`
	// 访问 [DescribeRegions:查询可以使用的阿里云地域](https://next.api.aliyun.com/api/Ecs/2014-05-26/DescribeRegions) 查阅
	// 国内一般是去掉 ECS 所在可用区的后缀，比如去掉 cn-guangzhou-b 的尾号 -b

	client, err := ecs.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)

	if err != nil {
		fmt.Print(err.Error())
	}
	for key, _ := range ipSet {
		addSecurityGroupRules(client, "all", key, portRange)
	}
}
