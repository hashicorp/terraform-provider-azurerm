
## `github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies` Documentation

The `collectorpolicies` SDK allows for interaction with Azure Resource Manager `networkfunction` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
```


### Client Initialization

```go
client := collectorpolicies.NewCollectorPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CollectorPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := collectorpolicies.NewCollectorPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorName", "collectorPolicyName")

payload := collectorpolicies.CollectorPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CollectorPoliciesClient.Delete`

```go
ctx := context.TODO()
id := collectorpolicies.NewCollectorPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorName", "collectorPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CollectorPoliciesClient.Get`

```go
ctx := context.TODO()
id := collectorpolicies.NewCollectorPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorName", "collectorPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CollectorPoliciesClient.List`

```go
ctx := context.TODO()
id := collectorpolicies.NewAzureTrafficCollectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CollectorPoliciesClient.UpdateTags`

```go
ctx := context.TODO()
id := collectorpolicies.NewCollectorPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorName", "collectorPolicyName")

payload := collectorpolicies.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
