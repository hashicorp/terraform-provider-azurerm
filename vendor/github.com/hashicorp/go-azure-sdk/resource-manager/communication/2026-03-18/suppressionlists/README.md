
## `github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/suppressionlists` Documentation

The `suppressionlists` SDK allows for interaction with Azure Resource Manager `communication` (API Version `2026-03-18`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/suppressionlists"
```


### Client Initialization

```go
client := suppressionlists.NewSuppressionListsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SuppressionListsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName")

payload := suppressionlists.SuppressionListResource{
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


### Example Usage: `SuppressionListsClient.Delete`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionListsClient.Get`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionListsClient.ListByDomain`

```go
ctx := context.TODO()
id := suppressionlists.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

// alternatively `client.ListByDomain(ctx, id)` can be used to do batched pagination
items, err := client.ListByDomainComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SuppressionListsClient.SuppressionListAddressesCreateOrUpdate`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName", "addressId")

payload := suppressionlists.SuppressionListAddressResource{
	// ...
}


read, err := client.SuppressionListAddressesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionListsClient.SuppressionListAddressesDelete`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName", "addressId")

read, err := client.SuppressionListAddressesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionListsClient.SuppressionListAddressesGet`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName", "addressId")

read, err := client.SuppressionListAddressesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SuppressionListsClient.SuppressionListAddressesList`

```go
ctx := context.TODO()
id := suppressionlists.NewSuppressionListID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "suppressionListName")

// alternatively `client.SuppressionListAddressesList(ctx, id)` can be used to do batched pagination
items, err := client.SuppressionListAddressesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
