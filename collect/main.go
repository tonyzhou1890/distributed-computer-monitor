package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Config struct {
	Hostname string
	Server   Server
}

type Server struct {
	Address string
}

// 记录数据类型
type RecordDataType struct {
	Hostname string  `json:"hostname"`
	CpuType  string  `json:"cpuType"`
	CpuCore  int     `json:"cpuCore"`
	CpuLoad  float64 `json:"cpuLoad"`
	CpuTemp  float64 `json:"cpuTemp"`
	MemCap   int32   `json:"memCap"`
	MemLoad  float64 `json:"memLoad"`
	MemTemp  float64 `json:"memTemp"`
	SwapCap  int32   `json:"swapCap"`
	SwapLoad float64 `json:"swapLoad"`
	GpuType  string  `json:"gpuType"`
	GpuLoad  float64 `json:"gpuLoad"`
	GpuTemp  float64 `json:"gpuTemp"`
	NetUp    int     `json:"netUp"`
	NetDown  int     `json:"netDown"`
}

var config = Config{}

func main() {
	// 初始化配置
	var file = "./config.json"

	jsonConfig, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Print(jsonConfig)

	err = json.Unmarshal(jsonConfig, &config)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	if config.Hostname == "" {
		// log.Fatal("未配置主机名")
		config.Hostname = getHostname()
	}
	jsonstr, err := json.MarshalIndent(config, "", "  ")
	if err == nil {
		fmt.Print(string(jsonstr) + "\n")
	}
	if config.Server.Address == "" {
		log.Fatal("未配置服务器地址")
	}

	// 定时器
	cron2 := cron.New() //创建一个cron实例

	//执行定时任务（每分钟执行一次）
	sendInfo()
	_, err = cron2.AddFunc("@every 1m", sendInfo)
	if err != nil {
		log.Fatal(err)
	}

	//启动/关闭
	cron2.Start()
	defer cron2.Stop()
	fmt.Print("采集器运行中\r\n")
	select {
	//查询语句，保持程序运行，在这里等同于for{}
	}
}

// 发送信息
func sendInfo() {
	var info RecordDataType
	info.Hostname = config.Hostname
	// 获取信息
	info.CpuType = getCpuType()
	info.CpuCore = getCpuCore()
	info.CpuLoad = getCpuLoad()
	info.CpuTemp = getCpuTemp()
	info.MemCap = getMemCap()
	info.MemLoad = getMemLoad()
	info.SwapCap = getSwapCap()
	info.SwapLoad = getSwapLoad()
	// 打印 json
	// jsonstr, err := json.MarshalIndent(info, "", "  ")
	// if err == nil {
	// 	fmt.Print(string(jsonstr))
	// }
	// 发送请求
	jsonstr, err := json.Marshal(info)
	if err != nil {
		return
	}
	res, err := http.Post(config.Server.Address+"/api/collect/record/create", "application/json", bytes.NewBuffer(jsonstr))
	if err != nil {
		fmt.Print(err.Error() + "\n")
		return
	}
	if res.StatusCode != 200 {
		fmt.Print("接口请求失败" + res.Status + "\n")
	}
	res.Body.Close()
	return
}

// 获取主机名称
func getHostname() string {
	hostInfo, err := host.Info()
	if err != nil {
		return ""
	}
	return hostInfo.Hostname
}

// 获取 cpu 类型
func getCpuType() string {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return ""
	}
	if len(cpuInfo) == 0 {
		return ""
	}
	return cpuInfo[0].ModelName
}

// 获取 cpu 核心数
func getCpuCore() int {
	cores, err := cpu.Counts(true)
	if err != nil {
		return 0
	}
	return cores
}

// 获取 cpu 使用率
func getCpuLoad() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0
	}
	return percent[0]
}

// 获取 cpu 温度
func getCpuTemp() float64 {
	// windows 暂无法取得温度
	sys := runtime.GOOS
	if sys == "windows" {
		return 0
	}

	cmd := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0
	}
	tempStr := strings.Replace(out.String(), "\n", "", -1)
	temp, err := strconv.Atoi(tempStr)
	if err != nil {
		return 0
	}
	return float64(temp) / 1000
}

// 获取内存容量
func getMemCap() int32 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	// 内存容量需要除以 1024 得到 kb。
	return int32(memInfo.Total / 1024)
}

// 获取内存负载
func getMemLoad() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return memInfo.UsedPercent
}

// 获取 swap 容量
func getSwapCap() int32 {
	memInfo, err := mem.SwapMemory()
	if err != nil {
		return 0
	}
	return int32(memInfo.Total / 1024)
}

// 获取 swap 负载
func getSwapLoad() float64 {
	memInfo, err := mem.SwapMemory()
	if err != nil {
		return 0
	}
	return memInfo.UsedPercent
}
