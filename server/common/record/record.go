package record

import "time"

// 记录数据类型
type RecordDataType struct {
	ID       int     `field:"id" json:"id"`
	HostId   int     `field:"host_id" json:"hostId"`
	Hostname string  `json:"hostname" binding:"required"`
	CpuType  string  `field:"cpu_type" json:"cpuType"`
	CpuCore  int     `field:"cpu_core" json:"cpuCore"`
	CpuLoad  float32 `field:"cpu_load" json:"cpuLoad"`
	CpuTemp  float32 `field:"cpu_temp" json:"cpuTemp"`
	MemCap   int     `field:"mem_cap" json:"memCap"`
	MemLoad  float32 `field:"mem_load" json:"memLoad"`
	MemTemp  float32 `field:"mem_temp" json:"memTemp"`
	SwapCap  int     `field:"swap_cap" json:"swapCap"`
	SwapLoad float32 `field:"swap_load" json:"swapLoad"`
	GpuType  string  `field:"gpu_type" json:"gpuType"`
	GpuLoad  float32 `field:"gpu_load" json:"gpuLoad"`
	GpuTemp  float32 `field:"gpu_temp" json:"gpuTemp"`
	NetUp    int     `field:"net_up" json:"netUp"`
	NetDown  int     `field:"net_down" json:"netDown"`
	Created  string  `field:"created" json:"created"`
}

// 获取记录列表
type RecordListReqType struct {
	HostId    int       `form:"hostId" binding:"required"`
	StartTime time.Time `form:"startTime" time_format:"2006-01-02 15:04:05"`
	EndTime   time.Time `form:"endTime" time_format:"2006-01-02 15:04:05"`
}
