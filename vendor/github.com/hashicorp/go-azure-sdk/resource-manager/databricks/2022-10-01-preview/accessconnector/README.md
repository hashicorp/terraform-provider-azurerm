
## `github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector` Documentation

The `accessconnector` SDK allows for interaction with Azure Resource Manager `databricks` (API Version `2022-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector"
```


### Client Initialization

```go
client := accessconnector.NewAccessConnectorClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccessConnectorClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := accessconnector.NewAccessConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accessConnectorName")

payload := accessconnector.AccessConnector{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AccessConnectorClient.Delete`

```go
ctx := context.TODO()
id := accessconnector.NewAccessConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accessConnectorName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AccessConnectorClient.Get`

```go
ctx := context.TODO()
id := accessconnector.NewAccessConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accessConnectorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccessConnectorClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccessConnectorClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccessConnectorClient.Update`

```go
ctx := context.TODO()
id := accessconnector.NewAccessConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accessConnectorName")

payload := accessconnector.AccessConnectorUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
