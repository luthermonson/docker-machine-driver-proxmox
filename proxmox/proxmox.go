package proxmox

import (
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
	"github.com/rancher/machine/libmachine/state"
)

type Driver struct{}

func NewDriver() drivers.Driver {
	return &Driver{}
}

func (d *Driver) Create() error {
	return nil
}
func (d *Driver) DriverName() string {
	return ""
}
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{}
}
func (d *Driver) GetIP() (string, error) {
	return "", nil
}
func (d *Driver) GetMachineName() string {
	return ""
}
func (d *Driver) GetSSHHostname() (string, error) {
	return "", nil
}
func (d *Driver) GetSSHKeyPath() string {
	return ""
}
func (d *Driver) GetSSHPort() (int, error) {
	return 0, nil
}
func (d *Driver) GetSSHUsername() string {
	return ""
}
func (d *Driver) GetURL() (string, error) {
	return "", nil
}
func (d *Driver) GetState() (state.State, error) {
	return state.None, nil
}
func (d *Driver) Kill() error {
	return nil
}
func (d *Driver) PreCreateCheck() error {
	return nil
}
func (d *Driver) Remove() error {
	return nil
}
func (d *Driver) Restart() error {
	return nil
}
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	return nil
}
func (d *Driver) Start() error {
	return nil
}
func (d *Driver) Stop() error {
	return nil
}
