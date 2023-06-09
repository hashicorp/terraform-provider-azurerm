# Example: Assert an App Service Certificate's expiry date using Terraform Continuous Validation Checks

This example provisions an App Service certificate with a user supplied certificate and a Terraform Cloud (TFC) Check on the certificates expiry date to ensure it is valid for at least 30 days.

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.
