
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/deletedservice` Documentation

The `deletedservice` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/deletedservice"
```


### Client Initialization

```go
client := deletedservice.NewDeletedServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedServiceClient.GetByName`

```go
ctx := context.TODO()
id := deletedservice.NewDeletedServiceID("12345678-1234-9876-4563-123456789012", "locationName", "deletedServiceName")

read, err := client.GetByName(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedServiceClient.ListBySubscription`

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


### Example Usage: `DeletedServiceClient.Purge`

```go
ctx := context.TODO()
id := deletedservice.NewDeletedServiceID("12345678-1234-9876-4563-123456789012", "locationName", "deletedServiceName")

if err := client.PurgeThenPoll(ctx, id); err != nil {
	// handle the error
}
```
