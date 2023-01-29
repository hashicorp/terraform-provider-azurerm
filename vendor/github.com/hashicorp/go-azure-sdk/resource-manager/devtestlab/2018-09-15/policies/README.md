
## `github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies` Documentation

The `policies` SDK allows for interaction with the Azure Resource Manager Service `devtestlab` (API Version `2018-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
```


### Client Initialization

```go
client := policies.NewPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := policies.NewPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "policySetValue", "policyValue")

payload := policies.Policy{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoliciesClient.Delete`

```go
ctx := context.TODO()
id := policies.NewPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "policySetValue", "policyValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoliciesClient.Get`

```go
ctx := context.TODO()
id := policies.NewPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "policySetValue", "policyValue")

read, err := client.Get(ctx, id, policies.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoliciesClient.List`

```go
ctx := context.TODO()
id := policies.NewPolicySetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "policySetValue")

// alternatively `client.List(ctx, id, policies.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, policies.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoliciesClient.Update`

```go
ctx := context.TODO()
id := policies.NewPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "policySetValue", "policyValue")

payload := policies.UpdateResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
