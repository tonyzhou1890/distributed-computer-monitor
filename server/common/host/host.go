package host

// 主机数据类型
type HostDataType struct {
	ID       int    `field:"id" json:"id"`
	Hostname string `field:"hostname" json:"hostname"`
	Created  string `field:"created" json:"created"`
	Updated  string `field:"updated" json:"updated"`
}
