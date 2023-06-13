# Example: Assert a VM's power state using Terraform Continuous Validation Checks

This example provisions a Linux Virtual Machine with a Terraform Cloud (TFC) Check on the VM's power state.

You can force the check to fail in this example by provisioning the VM and manually stopping it, and then triggering a health check in TFC. The check will fail and report that the VM is not running.

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.
