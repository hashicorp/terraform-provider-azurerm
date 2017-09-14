Deploys a group of Virtual Machines exposed to a public IP via a Load Balancer
==============================================================================

This Terraform module deploys a Virtual Machines Scale Set in Azure with the following characteristics: 

- Creates a virtual network and a subnet using the 10.0.0.0/16 and 10.0.1.0/24 address space.
- Creates a load balancer for the scale set of virtual machines
- Exposes through NAT one or several ports of the VMs on the load balancer

Module Input Variables 
----------------------

- `resource_group_name` - The name of the resource group in which the resources will be created.
- `location` - The Azure location where the resources will be created.
- `vm_size` - The initial size of the virtual machine that will be used in the VM Scale Set.
- `admin_username` - The name of the administrator to access the machines part of the virtual machine scale set. 
- `admin_password` - The password of the administrator account. The password must comply with the complexity requirements for Azure virtual machines.
- `ssh_key` - The path on the local machine of the ssh public key in the case of a Linux deployment.  
- `nb_instance` - The number of instances that will be initially deployed in the virtual machine scale set.
- `protocol` - A map representing the protocols and ports to open on the load balancer in front of the virtual machine scale set.
- `vm_os_publisher` - The name of the publisher of the image that you want to deploy, for example "Canonical".
- `vm_os_offer` - The name of the offer of the image that you want to deploy, for example "UbuntuServer"
- `vm_os_sku` - The sku of the image that you want to deploy, for example "14.04.2-LTS"
- `vm_os_id` - The ID of the image that you want to deploy if you are using a custom image.
- `lb_port` - Protocols to be used for the load balancer rules [frontend_port, protocol, backend_port]. Set to blank to disable.
- `tags` - A map of the tags to use on the resources that are deployed with this module.

Usage
-----

```hcl 
module "computegroup" { 
    source              = "./path/to/module"
    resource_group_name = "my-resource-group"
    location            = "westus"
    vm_size             = "Standard_A0"
    admin_username      = "azureuser"
    admin_password      = "ComplexPassword"
    ssh_key             = "~/.ssh/id_rsa.pub"
    nb_instance         = 2
    vm_os_publisher     = "Canonical"
    vm_os_offer         = "UbuntuServer"
    vm_os_sku           = "14.04.2-LTS"
    vm_os_id            = ""
    lb_port             = { 
                            http = ["80", "Tcp", "80"]
                            https = ["443", "Tcp", "443"]
                          }
    tags                = {
                            environment = "dev"
                            costcenter  = "it"
                          }
}

```

Outputs
=======

- `vmss_name` - Name of the virtual machine scale set
- `vnet_id` - Id of the Virtual Network deployed as part of the VM scale set
- `subnet_id` - Id of the Subnet deployed as part of the VM scale set

Authors
=======
Originally created by [Damien Caro](http://github.com/dcaro)


