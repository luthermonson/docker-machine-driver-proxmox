package proxmox

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/luthermonson/go-proxmox"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/machine/libmachine/mcnflag"
	"github.com/rancher/machine/libmachine/state"
)

const DriverName = "proxmox"

type Driver struct {
	*drivers.BaseDriver
	client   *proxmox.Client
	ID       int
	node     *proxmox.Node
	template *proxmox.VirtualMachine
	vm       *proxmox.VirtualMachine

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
	return &Driver{}
}

func (d *Driver) Create() error {
	if d.node == nil || d.template == nil {
		return errors.New("node and template required")
	}

	newid, task, err := d.template.Clone(d.MachineName, d.Node)
	if err != nil {
		return err
	}

	if err := task.WaitFor(10); err != nil {
		return err
	}

	d.ID = newid
	d.vm, err = d.node.VirtualMachine(d.ID)
	if err != nil {
		return err
	}

	return d.Start()
}

func (d *Driver) DriverName() string {
	return DriverName
}

func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return flags
}

func (d *Driver) GetSSHHostname() (string, error) {
	return "", nil
}

func (d *Driver) GetURL() (string, error) {
	return "", nil
}

func (d *Driver) GetState() (state.State, error) {
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

	task, err := d.vm.Delete()
	if err != nil {
		return err
	}

	return task.Wait(1*time.Second, 30*time.Second)
}

func (d *Driver) Restart() error {
	t, err := d.vm.Reboot()
	if err != nil {
		return err
	}

	return t.WaitFor(15)
}

func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.ApiUrl = opts.String("proxmox-url")
	d.Username = opts.String("proxmox-username")
	d.Password = opts.String("proxmox-password")
	d.TwoFactorAuthCode = opts.String("proxmox-2fa-code")
	d.Insecure = opts.Bool("proxmox-insecure")
	d.TemplateId = opts.Int("proxmox-template-id")
	d.Node = opts.String("proxmox-node")
	d.TokenID = opts.String("proxmox-tokenid")
	d.Secret = opts.String("proxmox-secret")
	d.client = d.proxmoxClient()
	_, err := d.client.Version() // get version info to verify credentials

	return err
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

	return proxmox.NewClient(d.ApiUrl, options...)
}

func (d *Driver) setup() (err error) {
	if d.client == nil {
		d.client = d.proxmoxClient()
	}

	log.Debugf("finding node: %s", d.Node)
	d.node, err = d.client.Node(d.Node)
	if err != nil {
		return err
	}

	log.Debugf("finding template: %d", d.TemplateId)
	if d.TemplateId == 0 {
		return fmt.Errorf("template id has to be set")
	}

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
