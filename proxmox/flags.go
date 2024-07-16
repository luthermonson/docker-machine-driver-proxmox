package proxmox

import (
	"context"

	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/mcnflag"
)

var flags = []mcnflag.Flag{
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_METHOD",
		Name:   "proxmox-method",
		Usage:  "Method to put the ssh credentials from the box: agent (qemu-guest-agent), drive (PVE cloud-init drive), nocloud (cloud-init iso)",
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_URL",
		Name:   "proxmox-url",
		Usage:  "URL to the proxmox API Server",
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_USERNAME",
		Name:   "proxmox-username",
		Usage:  "Username for the proxmox API Server",
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_PASSWORD",
		Name:   "proxmox-password",
		Usage:  "Password for the proxmox API Server",
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_TOKENID",
		Name:   "proxmox-tokenid",
		Usage:  "Token ID for the proxmox API Server",
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_SECRET",
		Name:   "proxmox-secret",
		Usage:  "Secret for a TokenID for the proxmox API Server",
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_2FA_CODE",
		Name:   "proxmox-2fa-code",
		Usage:  "Two Factor Authentication code for logins, not required if 2fa not turned on",
	},
	mcnflag.BoolFlag{
		EnvVar: "PROXMOX_INSECURE",
		Name:   "proxmox-insecure",
		Usage:  "Skip TLS verification",
	},
	mcnflag.IntFlag{
		EnvVar: "PROXMOX_TIMEOUT",
		Name:   "proxmox-timeout",
		Usage:  "API timeout in seconds",
		Value:  30,
	},
	mcnflag.StringFlag{
		EnvVar: "PROXMOX_NODE",
		Name:   "proxmox-node",
		Usage:  "Node name the template is on",
	},
	mcnflag.IntFlag{
		EnvVar: "PROXMOX_TEMPLATE_ID",
		Name:   "proxmox-template-id",
		Usage:  "Id of the template to clone from",
	},
}

func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.Method = opts.String("proxmox-method")
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

	ctx := context.Background()
	_, err := d.client.Version(ctx) // get version info to verify credentials

	return err
}
