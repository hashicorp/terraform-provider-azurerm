
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestriction` Documentation

The `policyrestriction` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyrestriction"
```


### Client Initialization

```go
client := policyrestriction.NewPolicyRestrictionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyRestrictionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := policyrestriction.NewPolicyRestrictionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "policyRestrictionId")

payload := policyrestriction.PolicyRestrictionContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, policyrestriction.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyRestrictionClient.Delete`

```go
ctx := context.TODO()
id := policyrestriction.NewPolicyRestrictionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "policyRestrictionId")

read, err := client.Delete(ctx, id, policyrestriction.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyRestrictionClient.Get`

```go
ctx := context.TODO()
id := policyrestriction.NewPolicyRestrictionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "policyRestrictionId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyRestrictionClient.GetEntityTag`

```go
ctx := context.TODO()
id := policyrestriction.NewPolicyRestrictionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "policyRestrictionId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyRestrictionClient.Update`

```go
ctx := context.TODO()
id := policyrestriction.NewPolicyRestrictionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "policyRestrictionId")

payload := policyrestriction.PolicyRestrictionUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, policyrestriction.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
