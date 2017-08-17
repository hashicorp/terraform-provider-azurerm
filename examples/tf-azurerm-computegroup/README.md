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
- `vn_size` - The initial size of the virtual machine that will be used in the VM Scale Set.
- `admin_username` - The name of the administrator to access the machines part of the virtual machine scale set. 
- `admin_password` - The password of the administrator account. The password must comply with the complexity requirements for Azure virtual machines.
- `ssh_key` - The path on the local machine of the ssh public key in the case of a Linux deployment.  
- `nb_instace` - The number of instances that will be initially deployed in the virtual machine scale set.
- `protocol` - A map representing the protocols and ports to open on the load balancer in front of the virtual machine scale set.
- `os` - A map indentifying the operating system to use for that virtual machine scale set.
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
    ssh_key             = "~/.ssh/demo_key.pub"
    nb_instance         = 2
    protocol            = {
                            ssh = ["22", "Tcp", "50000", "50119"]
                            ftp = ["21", "Tcp", "51000", "51119"]
                          }
    os                  = {
                            id = ""
                            publisher = "Canonical"
                            offer = "UbuntuServer"
                            sku = "14.04.2-LTS"
                            version = "latest"
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
Originally created by Microsoft

License
=======

