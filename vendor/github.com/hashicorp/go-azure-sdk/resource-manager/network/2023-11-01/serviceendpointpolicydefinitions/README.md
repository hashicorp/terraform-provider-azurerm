
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicydefinitions` Documentation

The `serviceendpointpolicydefinitions` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicydefinitions"
```


### Client Initialization

```go
client := serviceendpointpolicydefinitions.NewServiceEndpointPolicyDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServiceEndpointPolicyDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := serviceendpointpolicydefinitions.NewServiceEndpointPolicyDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceEndpointPolicyName", "serviceEndpointPolicyDefinitionName")

payload := serviceendpointpolicydefinitions.ServiceEndpointPolicyDefinition{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceEndpointPolicyDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := serviceendpointpolicydefinitions.NewServiceEndpointPolicyDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceEndpointPolicyName", "serviceEndpointPolicyDefinitionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceEndpointPolicyDefinitionsClient.Get`

```go
ctx := context.TODO()
id := serviceendpointpolicydefinitions.NewServiceEndpointPolicyDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceEndpointPolicyName", "serviceEndpointPolicyDefinitionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServiceEndpointPolicyDefinitionsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := serviceendpointpolicydefinitions.NewServiceEndpointPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceEndpointPolicyName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
