
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raipolicies` Documentation

The `raipolicies` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/raipolicies"
```


### Client Initialization

```go
client := raipolicies.NewRaiPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RaiPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := raipolicies.NewRaiPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiPolicyName")

payload := raipolicies.RaiPolicy{
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


### Example Usage: `RaiPoliciesClient.Delete`

```go
ctx := context.TODO()
id := raipolicies.NewRaiPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RaiPoliciesClient.Get`

```go
ctx := context.TODO()
id := raipolicies.NewRaiPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiPoliciesClient.List`

```go
ctx := context.TODO()
id := raipolicies.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
