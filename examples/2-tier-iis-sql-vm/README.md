# Example: 2 Tier Architecture, IIS with MS SQL Backend joined to Active Directory

This example creates a VNet with the following subnets:

wafsubnet - To contain a WAF of your choosing (Not created as part of this example)
rpsubnet  - To contain a Reverse Proxy of your choosing (Not created as part of this example)
issubnet  - To contain IIS VM(s) (Created as part of this example)
dbsubnet  - To contain MS SQL VM(s) (Created as part of this example)

It creates the following VM's, creates an Active Directory Domain and joins the VM's to it:

- DC1 - Primary Domain Controller holding FSMO roles, Static Public & Private IP Addresses

- DC2 - Secondary Domain Controller joined to domain, Static Public & Private IP Addresses 

- IIS VM(s) - Scalable using the count function, IIS & Management Tools installed, Windows Server 2012 R2, Added to Availability Set, Static Public & Private IP Addresses Joined to Domain

- SQL VM(s) - Scalable using the count function, SQL 2014 SP2 & SSMS Installed, Windows Server 2012 R2, Windows Failover Clustering Service Installed, Added to Availability Set , Static Public & Private IP Addresses, Joined to Domain

It also creates an Azure internal load balancer (ILB) and adds the SQL VM(s) to the backend pool so can be expanded to use AlwaysOn Capability

Includes an environments folder containing a .tfvars file, to allow this Infrastructure to be deployed throughout a pipeline i.e. Dev,QA,UAT etc

## Notes

- This is intended as an example of creating a multi tier architecture joined to an Active Directory Domain, and **it is not recommended for production use** as the configuration has been simplified for example purposes, e.g.:
  - There's no security rules configured on the network, so everything's open internally etc.
  - Usernames / Passwords are in plaintext rather than using a secret store like Azure Key Vault
  - You can RDP to all servers over the public internet, rather than VPN or Bastion Host
- The numbering on the files within the modules below have no effect on which order the resources are created in - it's purely to make the examples easier to understand.


## Running this Example

Initialize the modules (and download the Azure Provider) by running `terraform init`:

```bash
$ terraform init
```

In order to run this example you'll need some kind of credentials configured - either a Service Principal or to be logged into the Azure CLI. You can find out more about this on [the Azure Provider overview page](https://www.terraform.io/docs/providers/azurerm/index.html)

Once you've initialized the Provider - you can run the sample by running:

```bash
$ terraform apply
```

This will take around 45m to provision - once completed you should see a resource group containing everything described above.

## Modules

This example makes use of 5 modules:
 * [modules/active-directory](modules/active-directory)
    - This module creates an Active Directory Forest on a single Virtual Machine
 * [modules/network](modules/network)
    - This module creates the Network with 4 subnets.
    - In a Production environment there would be [Network Security Rules](https://www.terraform.io/docs/providers/azurerm/r/network_security_rule.html) in effect which limited which ports can be used between these Subnets, however for the purposes of keeping this demonstration simple, these have been omitted.
 * [modules/dc2-vm](modules/dc2-vm)
    - This module creates a secondary domain controller machine for resiliency that is bound to the Active Directory Domain created in the `active-directory` module above.
 * [modules/iis-vm](modules/iis-vm)
    - This module creates IIS VM's - Choose how many you want using count  
 * [modules/sql-vm](modules/sql-vm)
    - This module creates SQL VM's - Choose how many you want using count
    - Also created the ILB so you could scale out to use AlwaysOn