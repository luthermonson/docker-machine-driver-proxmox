package proxmox

import (
	"crypto/tls"
	"fmt"

	api "github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
	"github.com/rancher/machine/libmachine/state"
)

const DriverName = "proxmox"

type Driver struct {
	ApiUrl            string
	Username          string
	Password          string
	TwoFactorAuthCode string
	Insecure          bool
	Timeout           int
	TemplateId        int

	client *api.Client
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
	if d.client == nil {
		return fmt.Errorf("no api client was created and we can not communicate with proxmox")
	}

	if d.TemplateId == 0 {
		return fmt.Errorf("template id has to be set")
	}

	vmref := api.NewVmRef(d.TemplateId)
	vminfo, err := d.client.GetVmInfo(vmref)
	if err != nil {
		return err
	}

	fmt.Println(vminfo)

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

	return d.login()
}
func (d *Driver) Start() error {
	return nil
}
func (d *Driver) Stop() error {
	return nil
}

func (d *Driver) login() (err error) {
	d.client, err = api.NewClient(d.ApiUrl, nil, &tls.Config{
		InsecureSkipVerify: d.Insecure,
	}, d.Timeout)

	if err != nil {
		return err
	}

	return d.client.Login(d.Username, d.Password, d.TwoFactorAuthCode)
}
