
## `github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/defaultparametergroups` Documentation

The `defaultparametergroups` SDK allows for interaction with Azure Resource Manager `horizondb` (API Version `2026-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/defaultparametergroups"
```


### Client Initialization

```go
client := defaultparametergroups.NewDefaultParameterGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DefaultParameterGroupsClient.Get`

```go
ctx := context.TODO()
id := defaultparametergroups.NewParameterGroupID("12345678-1234-9876-4563-123456789012", "locationName", "parameterGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DefaultParameterGroupsClient.ListBySubscriptionAndRegion`

```go
ctx := context.TODO()
id := defaultparametergroups.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListBySubscriptionAndRegion(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionAndRegionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
