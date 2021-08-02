package proxmox

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/luthermonson/go-proxmox"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
	"github.com/rancher/machine/libmachine/state"
)

const DriverName = "proxmox"

type Driver struct {
	*drivers.BaseDriver
	client            *proxmox.Client
	ApiUrl            string
	Username          string
	Password          string
	TwoFactorAuthCode string
	Insecure          bool
	Timeout           int
	TemplateId        int
	TemplateNode      string
	TokenID           string
	Secret            string
}

func NewDriver() drivers.Driver {
	return &Driver{}
}

func (d *Driver) Create() error {

	return nil
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
	return state.None, nil
}

func (d *Driver) Kill() error {
	return nil
}

func (d *Driver) PreCreateCheck() error {
	if d.client == nil {
		return fmt.Errorf("no api client was created")
	}

	if d.TemplateNode == "" {
		return fmt.Errorf("template node has to be set")
	}
	node, err := d.client.Node(d.TemplateNode)
	if err != nil {
		return err
	}

	if d.TemplateId == 0 {
		return fmt.Errorf("template id has to be set")
	}

	vm, err := node.VirtualMachine(d.TemplateId)
	if err != nil {
		return err
	}

	if !vm.Template {
		return fmt.Errorf("virtual machine id %d was not a template", d.TemplateId)
	}

	return nil
}

func (d *Driver) Remove() error {
	return nil
}

func (d *Driver) Restart() error {
	return nil
}

func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.ApiUrl = opts.String("proxmox-api-url")
	d.Username = opts.String("proxmox-username")
	d.Password = opts.String("proxmox-password")
	d.TwoFactorAuthCode = opts.String("proxmox-2fa-code")
	d.Insecure = opts.Bool("proxmox-insecure")
	d.TemplateId = opts.Int("proxmox-template-id")
	d.TemplateNode = opts.String("proxmox-template-node")
	d.TokenID = opts.String("proxmox-tokenid")
	d.Secret = opts.String("proxmox-secret")
	d.client = d.proxmoxClient()
	_, err := d.client.Version() // get version info to verify credentials

	return err
}

func (d *Driver) Start() error {
	return nil
}

func (d *Driver) Stop() error {
	return nil
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
