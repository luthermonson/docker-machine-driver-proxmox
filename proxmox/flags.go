package proxmox

import "github.com/rancher/machine/libmachine/mcnflag"

var flags = []mcnflag.Flag{
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
