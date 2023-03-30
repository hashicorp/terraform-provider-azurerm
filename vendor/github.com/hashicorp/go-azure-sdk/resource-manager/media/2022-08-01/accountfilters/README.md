
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/accountfilters` Documentation

The `accountfilters` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/accountfilters"
```


### Client Initialization

```go
client := accountfilters.NewAccountFiltersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountFiltersClient.AccountFiltersCreateOrUpdate`

```go
ctx := context.TODO()
id := accountfilters.NewAccountFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "accountFilterValue")

payload := accountfilters.AccountFilter{
	// ...
}


read, err := client.AccountFiltersCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountFiltersClient.AccountFiltersDelete`

```go
ctx := context.TODO()
id := accountfilters.NewAccountFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "accountFilterValue")

read, err := client.AccountFiltersDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountFiltersClient.AccountFiltersGet`

```go
ctx := context.TODO()
id := accountfilters.NewAccountFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "accountFilterValue")

read, err := client.AccountFiltersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountFiltersClient.AccountFiltersList`

```go
ctx := context.TODO()
id := accountfilters.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.AccountFiltersList(ctx, id)` can be used to do batched pagination
items, err := client.AccountFiltersListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountFiltersClient.AccountFiltersUpdate`

```go
ctx := context.TODO()
id := accountfilters.NewAccountFilterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "accountFilterValue")

payload := accountfilters.AccountFilter{
	// ...
}


read, err := client.AccountFiltersUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
