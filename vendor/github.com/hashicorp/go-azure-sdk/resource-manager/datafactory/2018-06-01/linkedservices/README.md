
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices` Documentation

The `linkedservices` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
```


### Client Initialization

```go
client := linkedservices.NewLinkedServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinkedServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "linkedServiceName")

payload := linkedservices.LinkedServiceResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, linkedservices.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedServicesClient.Delete`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "linkedServiceName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedServicesClient.Get`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "linkedServiceName")

read, err := client.Get(ctx, id, linkedservices.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedServicesClient.ListByFactory`

```go
ctx := context.TODO()
id := linkedservices.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
