package util

import (
	"time"

	"github.com/sony/sonyflake"
)

var (
	sonyFlake *sonyflake.Sonyflake
	// 定义一个全局的 machineID 模拟获取
	// 现实环境中应从 zk 或 etcd 中获取
	machineID uint16
)

// 获取 机器编码ID的 回调函数
func getMachineID() (uint16, error) {
	// machineID 返回nil, 则返回专用IP地址的低16位
	return machineID, nil
}

// 初始化 sonyFlake 配置
func Init(mID uint16) {
	machineID = mID
	st := sonyflake.Settings{StartTime: time.Time{}, MachineID: getMachineID}
	sonyFlake = sonyflake.NewSonyflake(st)
}

// 获取全局 ID 的函数
func GetUID() (id uint64, err error) {
	if sonyFlake == nil {
		Init(0)
	}
	return sonyFlake.NextID()
}
