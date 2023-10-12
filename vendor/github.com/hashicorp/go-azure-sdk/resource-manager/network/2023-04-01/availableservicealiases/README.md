
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/availableservicealiases` Documentation

The `availableservicealiases` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/availableservicealiases"
```


### Client Initialization

```go
client := availableservicealiases.NewAvailableServiceAliasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AvailableServiceAliasesClient.List`

```go
ctx := context.TODO()
id := availableservicealiases.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AvailableServiceAliasesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := availableservicealiases.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
