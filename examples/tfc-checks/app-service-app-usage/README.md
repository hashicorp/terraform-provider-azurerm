`# Example: Assert if an App Service Function or Web App has exceeded its usage limit using Terraform Continuous Validation Checks

This example provisions aa App Service Function with a Terraform Cloud (TFC) Check on the usage state to ensure it has not been exceeded

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.
