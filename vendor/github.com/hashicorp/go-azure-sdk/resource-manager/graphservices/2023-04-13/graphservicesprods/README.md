
## `github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods` Documentation

The `graphservicesprods` SDK allows for interaction with Azure Resource Manager `graphservices` (API Version `2023-04-13`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
```


### Client Initialization

```go
client := graphservicesprods.NewGraphservicesprodsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GraphservicesprodsClient.AccountsCreateAndUpdate`

```go
ctx := context.TODO()
id := graphservicesprods.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

payload := graphservicesprods.AccountResource{
	// ...
}


if err := client.AccountsCreateAndUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GraphservicesprodsClient.AccountsDelete`

```go
ctx := context.TODO()
id := graphservicesprods.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

read, err := client.AccountsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphservicesprodsClient.AccountsGet`

```go
ctx := context.TODO()
id := graphservicesprods.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

read, err := client.AccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphservicesprodsClient.AccountsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.AccountsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GraphservicesprodsClient.AccountsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AccountsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GraphservicesprodsClient.AccountsUpdate`

```go
ctx := context.TODO()
id := graphservicesprods.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

payload := graphservicesprods.TagUpdate{
	// ...
}


read, err := client.AccountsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
