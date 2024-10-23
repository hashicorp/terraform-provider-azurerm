
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/availabledelegations` Documentation

The `availabledelegations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/availabledelegations"
```


### Client Initialization

```go
client := availabledelegations.NewAvailableDelegationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AvailableDelegationsClient.AvailableDelegationsList`

```go
ctx := context.TODO()
id := availabledelegations.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.AvailableDelegationsList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableDelegationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AvailableDelegationsClient.AvailableResourceGroupDelegationsList`

```go
ctx := context.TODO()
id := availabledelegations.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.AvailableResourceGroupDelegationsList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableResourceGroupDelegationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
