# Example: a Function App with Azure RBAC Role Assignment

This example provisions a Function App with Azure RBAC Role Assignment.

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.
- `role_definition_name` - (Optional) Desired role to assign your function (Reader, Contributor, Owner, etc.) Defaults to `Reader`.

## Outputs

- `account_id` - The Principal ID of the RBAC identity.
