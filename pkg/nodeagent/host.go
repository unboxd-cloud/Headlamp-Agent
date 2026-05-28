package nodeagent

import (
	"os"
	"runtime"
	"time"
)

type HostInventory struct {
	Hostname string    `json:"hostname"`
	OS       string    `json:"os"`
	Arch     string    `json:"arch"`
	Go       string    `json:"goVersion"`
	Time     time.Time `json:"time"`
}

func CollectHostInventory() HostInventory {
	hostname, _ := os.Hostname()
	return HostInventory{
		Hostname: hostname,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Go:       runtime.Version(),
		Time:     time.Now().UTC(),
	}
}
