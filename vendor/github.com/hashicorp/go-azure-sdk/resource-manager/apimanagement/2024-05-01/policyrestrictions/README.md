
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestrictions` Documentation

The `policyrestrictions` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestrictions"
```


### Client Initialization

```go
client := policyrestrictions.NewPolicyRestrictionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyRestrictionsClient.PolicyRestrictionListByService`

```go
ctx := context.TODO()
id := policyrestrictions.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.PolicyRestrictionListByService(ctx, id)` can be used to do batched pagination
items, err := client.PolicyRestrictionListByServiceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
