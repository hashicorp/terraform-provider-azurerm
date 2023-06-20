
## `github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/creators` Documentation

The `creators` SDK allows for interaction with the Azure Resource Manager Service `maps` (API Version `2021-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/creators"
```


### Client Initialization

```go
client := creators.NewCreatorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CreatorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := creators.NewCreatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "creatorValue")

payload := creators.Creator{
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


### Example Usage: `CreatorsClient.Delete`

```go
ctx := context.TODO()
id := creators.NewCreatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "creatorValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CreatorsClient.Get`

```go
ctx := context.TODO()
id := creators.NewCreatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "creatorValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CreatorsClient.ListByAccount`

```go
ctx := context.TODO()
id := creators.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

// alternatively `client.ListByAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CreatorsClient.Update`

```go
ctx := context.TODO()
id := creators.NewCreatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "creatorValue")

payload := creators.CreatorUpdateParameters{
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
