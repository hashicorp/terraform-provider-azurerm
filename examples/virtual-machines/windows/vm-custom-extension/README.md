## Example: Virtual Machine with Custom Script Extension

This example is used to demonstrate how to create a custom script extension for a Windows Virtual Machine. In this particular example the custom extension configures the Windows Virtual Machine to be discoverable and accessible by ansible.

This example provisions:
- Windows Virtual Machine
- A Virtual Machine Extension

### Variables

- `location` - (Required) Azure Region in which all resources in this example should be provisioned
