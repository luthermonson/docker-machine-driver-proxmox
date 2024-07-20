package proxmox

import (
	"context"
	"errors"
	"os"
)

func (d *Driver) startAgent() error {
	//if err := d.vm.WaitForAgent(10); err != nil {
	//	return err
	//}
	//
	//if err := d.Start(); err != nil {
	//	return err
	//}

	return nil
}

func (d *Driver) startDrive() error {
	// create a snippet and do something like this
	//task, err := d.vm.Config(proxmox.VirtualMachineOption{
	//	Name:  "cicustom",
	//	Value:  "user=local:snippets/userconfig.yaml",
	//})

	return errors.New("method not supported: drive")
}

func (d *Driver) startNoCloud() error {
	userData, metaData, err := d.buildCloudInit()
	ctx := context.Background()
	if err != nil {
		return err
	}
	if err := d.vm.CloudInit(ctx, "scsi1", userData, metaData, "", ""); err != nil {
		return err
	}

	return d.Start()
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) buildCloudInit() (string, string, error) {
	sshkey, err := os.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", "", err
	}

	return `#cloud-config
users:
  - name: ` + d.GetSSHUsername() + `
    sudo: "ALL=(ALL) NOPASSWD:ALL"
    lock_passwd: true
    groups: staff
    create_groups: false
    no_user_group: true
    ssh_authorized_keys:
    -  ` + string(sshkey) + `
groups:
  - staff
`, `instance-id: iid-` + d.MachineName + `
hostname: ` + d.MachineName + `
`, nil
}
