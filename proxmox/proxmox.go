package proxmox

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/luthermonson/go-proxmox"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
	"github.com/rancher/machine/libmachine/ssh"
	"github.com/rancher/machine/libmachine/state"
)

const (
	DriverName = "proxmox"
	B2DUser    = "docker"
)

type Driver struct {
	*drivers.BaseDriver
	client   *proxmox.Client
	ID       int
	node     *proxmox.Node
	template *proxmox.VirtualMachine
	vm       *proxmox.VirtualMachine

	Method            string
	ApiUrl            string
	Username          string
	Password          string
	TwoFactorAuthCode string
	Insecure          bool
	Timeout           int
	TemplateId        int
	Node              string
	TokenID           string
	Secret            string
}

func NewDriver() drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			SSHUser: B2DUser,
		},
	}
}

func (d *Driver) Create() error {
	if d.node == nil || d.template == nil {
		return errors.New("node and template required")
	}

	newid, task, err := d.template.Clone(&proxmox.VirtualMachineCloneOptions{
		Name: d.MachineName,
		Full: 1,
	})

	if err != nil {
		return err
	}

	if err := task.WaitFor(30); err != nil {
		return err
	}

	d.ID = newid
	d.vm, err = d.node.VirtualMachine(d.ID)
	if err != nil {
		return err
	}

	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	var starter func() error

	switch d.Method {
	case "agent":
		starter = d.startAgent
	case "drive":
		starter = d.startDrive
	case "nocloud":
		starter = d.startNoCloud
	default:
		return fmt.Errorf("method %s is not supported", d.Method)
	}

	if err := starter(); err != nil {
		return err
	}

	return d.waitForIP()
}

func (d *Driver) waitForIP() error {
	// todo only supports agent, add more methods to find ip
	// todo only supports Net0
	// todo only supports ipv4

	if err := d.vm.WaitForAgent(10); err != nil {
		return err
	}

	net := d.vm.VirtualMachineConfig.Net0

RETRY:
	ifaces, err := d.vm.AgentGetNetworkIFaces()
	if err != nil {
		return err
	}

	for _, iface := range ifaces {
		if strings.Contains(strings.ToLower(net), strings.ToLower(iface.HardwareAddress)) {
			for _, ip := range iface.IPAddresses {
				if ip.IPAddressType == "ipv4" {
					d.IPAddress = ip.IPAddress
				}
			}
		}
	}

	if d.IPAddress == "" {
		time.Sleep(2 * time.Second)
		goto RETRY
	}

	return nil
}

func (d *Driver) DriverName() string {
	return DriverName
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return flags
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.IPAddress, nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetState() (state.State, error) {
	if err := d.setup(); err != nil {
		return state.None, err
	}

	if err := d.vm.Ping(); err != nil {
		return state.None, err
	}

	if d.vm.IsStopped() {
		return state.Stopped, nil
	}

	if d.vm.IsRunning() {
		return state.Running, nil
	}

	return state.None, nil
}

func (d *Driver) Kill() error {
	t, err := d.vm.Stop()
	if err != nil {
		return err
	}

	if err := t.WaitFor(15); err != nil {
		return err
	}

	t, err = d.vm.Delete()
	if err != nil {
		return err
	}

	return t.WaitFor(15)
}

func (d *Driver) PreCreateCheck() error {
	if d.client == nil {
		return fmt.Errorf("no api client was created")
	}

	if d.Node == "" {
		return fmt.Errorf("template node has to be set")
	}

	return d.setup()
}

func (d *Driver) Remove() error {
	if err := d.setup(); err != nil {
		return err
	}

	if d.vm == nil {
		return nil
	}

	return d.Kill()
}

func (d *Driver) Restart() error {
	t, err := d.vm.Reboot()
	if err != nil {
		return err
	}

	return t.WaitFor(15)
}

func (d *Driver) Start() error {
	t, err := d.vm.Start()
	if err != nil {
		return err
	}

	return t.WaitFor(15)
}

func (d *Driver) Stop() error {
	t, err := d.vm.Stop()
	if err != nil {
		return err
	}

	return t.WaitFor(15)
}

func (d *Driver) proxmoxClient() *proxmox.Client {
	var options []proxmox.Option
	if d.Insecure {
		options = append(options, proxmox.WithClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}))
	}

	if d.TokenID != "" && d.Secret != "" {
		options = append(options, proxmox.WithAPIToken(d.TokenID, d.Secret))
	}

	if d.Username != "" && d.Password != "" {
		options = append(options, proxmox.WithLogins(d.Username, d.Password))
	}

	options = append(options, proxmox.WithLogger(logger))
	return proxmox.NewClient(d.ApiUrl, options...)
}

func (d *Driver) setup() (err error) {
	if d.TemplateId == 0 {
		return fmt.Errorf("template id has to be set")
	}

	if d.client == nil {
		d.client = d.proxmoxClient()
	}

	logger.Debugf("finding node: %s", d.Node)
	d.node, err = d.client.Node(d.Node)
	if err != nil {
		return err
	}

	logger.Debugf("finding template: %d", d.TemplateId)
	d.template, err = d.node.VirtualMachine(d.TemplateId)
	if err != nil {
		return err
	}

	if d.ID == 0 {
		return nil
	}

	d.vm, err = d.node.VirtualMachine(d.ID)
	if err != nil {
		return err
	}

	return err
}
