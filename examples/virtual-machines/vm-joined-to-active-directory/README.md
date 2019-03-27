# Example: Joining a Windows Machine to an Active Directory Domain

This example creates an Active Directory Domain, a Windows Client; to demonstrate how to bind a Windows Client to an Active Directory Domain using a Virtual Machine Extension in Terraform (using [the `azurerm_virtual_machine_extension` resource](https://www.terraform.io/docs/providers/azurerm/r/virtual_machine_extension.html)).

This example is built around [the Virtual Machine Extension](https://www.terraform.io/docs/providers/azurerm/r/virtual_machine_extension.html) found in the `windows-client` module (documented below) - which demonstrates binding a Windows Virtual Machine to an Active Directory Domain. For the purposes of this example we create an Active Directory Domain, since it's easier to demonstrate - however you can achieve the same thing with an existing Active Directory Domain.

## Notes

- This is intended as an example of binding machines to an Active Directory Domain, and **it is not recommended for production use** as the configuration has been simplified for example purposes, e.g.:
  - The Active Directory Forest has a single node, for demonstration purposes
  - There's no security rules configured on the network, so everything's open internally etc.
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

This will take around 20m to provision - once completed you should see the Public IP Address of the Windows Client machine (which is bound to the Active Directory Domain):

```bash
windows_client_public_ip = 0.0.0.0
```

## Variables

 * `prefix` - The prefix used for all resources in this example. Needs to be a short (6 characters) alphanumeric string. Example: `addemo`.
 * `admin_username` - The username of the administrator account for both the local accounts, and Active Directory accounts. Example: `myexampleadmin`
 * `admin_password` - The password of the administrator account for both the local accounts, and Active Directory accounts. Needs to comply with the Windows Password Policy. Example: `PassW0rd1234!`

## Architecture

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                                    Internal Network                                    │
└────────────────────────────────────────────────────────────────────────────────────────┘
                                             ▲
                         ┌───────────────────┴────────────────────┐
                         │                                        │
          ┌─────────────────────────────┐          ┌─────────────────────────────┐
          │  Domain Controllers Subnet  │          │    Domain Clients Subnet    │
          └─────────────────────────────┘          └─────────────────────────────┘
                         ▲                                        ▲
                         │                                        │
         ┌───────────────────────────────┐        ┌───────────────────────────────┐
         │    Domain Controllers NIC     │        │      Domain Clients NIC       │
         │                               │        │                               │
         │     ({prefix}-dc-primary)     │        │     ({prefix}-client-nic)     │
         └───────────────────────────────┘        └───────────────────────────────┘
                         ▲                                        ▲
                         │                                        │
         ┌───────────────────────────────┐        ┌───────────────────────────────┐
         │     Domain Controller VM      │        │       Domain Client VM        │
         │                               │        │                               │
         │         ({prefix}-dc)         │        │       ({prefix}-client)       │
         └───────────────────────────────┘        └───────────────────────────────┘
                         ▲                                        ▲
                         │                             ┌──────────┴───────────┐
        ┌────────────────────────────────┐             │                      │
        │   Virtual Machine Extension    │    ┌─────────────────┐ ┌───────────────────────┐
        │                                │    │  VM Extension   │ │       Public IP       │
        │(create-active-directory-forest)│    │                 │ │                       │
        └────────────────────────────────┘    │  (join-domain)  │ │ ({prefix}-client-nic) │
                                              └─────────────────┘ └───────────────────────┘
```

## Modules

This example makes use of 3 modules:
 * [modules/active-directory](modules/active-directory)
    - This module creates an Active Directory Forest on a single Virtual Machine
 * [modules/network](modules/network)
    - This module creates the Network with 2 subnets, one for the Domain Controller and another for the Clients.
    - In a Production environment there would be [Network Security Rules](https://www.terraform.io/docs/providers/azurerm/r/network_security_rule.html) in effect which limited which ports can be used between these Subnets, however for the purposes of keeping this demonstration simple, these have been omitted.
 * [modules/windows-client](modules/windows-client)
    - This module creates a Windows Client machine that is bound to the Active Directory Domain created in the `active-directory` module above.
    - This module includes a sleep function designed to wait for 10 minutes (until the Active Directory Domain has provisioned) - however this isn't ideal for a number of reasons. In a Production Environment it's likely your Active Directory Domain already exists.
