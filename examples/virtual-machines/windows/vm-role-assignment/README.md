## Example: Virtual Machine with Role Assignment

This example is used to demonstrate how to perform a role assigned to a Virtual Machine (with System Assigned identity) to grant it permissions to access another Azure resource.

This example provisions:
- Windows Virtual Machine with System Assigned identity
- An Azure Storage Account
- Performs a Role Assignment granting Reader role to the Windows Virtual Machine to access the storage account.

Similar pattern can be used:
- to grant any other Azure built-in role (e.g. Contributor, Owner), or a custom user defined role.
- to grant access to Azure resources other than a storage account (e.g. Resource Group, Service Bus, Event Hubs)
- with Virtual Machines with User Assigned identity
- with Linux Virtual Machines

### Variables

- `admin` - (Required) Virtual Machine Admin username
- `adminPassword` - (Required) Virtual Machine Admin Password
- `location` - (Required) Azure Region in which all resources in this example should be provisioned
- `storageaccount` - (Required) Name of the storage account to be provisioned
