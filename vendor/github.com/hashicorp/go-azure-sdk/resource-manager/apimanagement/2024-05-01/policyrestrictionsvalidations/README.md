
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestrictionsvalidations` Documentation

The `policyrestrictionsvalidations` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestrictionsvalidations"
```


### Client Initialization

```go
client := policyrestrictionsvalidations.NewPolicyRestrictionsValidationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyRestrictionsValidationsClient.PolicyRestrictionValidationsByService`

```go
ctx := context.TODO()
id := policyrestrictionsvalidations.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

if err := client.PolicyRestrictionValidationsByServiceThenPoll(ctx, id); err != nil {
	// handle the error
}
```
