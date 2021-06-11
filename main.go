package main

import (
	"github.com/luthermonson/docker-machine-driver-proxmox/proxmox"
	"github.com/rancher/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(proxmox.NewDriver())
}
