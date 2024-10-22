
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/servicetags` Documentation

The `servicetags` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/servicetags"
```


### Client Initialization

```go
client := servicetags.NewServiceTagsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServiceTagsClient.ServiceTagInformationList`

```go
ctx := context.TODO()
id := servicetags.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ServiceTagInformationList(ctx, id, servicetags.DefaultServiceTagInformationListOperationOptions())` can be used to do batched pagination
items, err := client.ServiceTagInformationListComplete(ctx, id, servicetags.DefaultServiceTagInformationListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServiceTagsClient.ServiceTagsList`

```go
ctx := context.TODO()
id := servicetags.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.ServiceTagsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
