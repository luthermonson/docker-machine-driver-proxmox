# Docker Machine Driver Proxmox

## Methods
There are three methods for this machine driver to give credentials to the virtual machine it creates.
* cloud-init NoCloud
* cloud-init Drive
* qemu-guest-agent

### cloud-init NoCloud
Dynamically create the ssh user login credentials on the virtual machine using an ISO mounted as a CDROM. This machine driver will create the ISO and upload it as `user-data-<vmid>.iso` to the ISO storage for the node the VM is going to be created on. 

#### Template Requirements
* cloud-init
* qemu-guest-agent (trying to remove, need today for node ip)
* virtio-scsi-pci

### CloudInit Drive (pending snippets api support)
Use the built in CloudInit drive functionality in ProxmoxVE. This can only work when API access to snippet drives is available.

### qemu-guest-agent (pending implementation)
Use the qemu guest agane to send commands to the virtual machine to add the user and SSH key.
