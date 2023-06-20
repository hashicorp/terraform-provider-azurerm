
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/tableserviceproperties` Documentation

The `tableserviceproperties` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/tableserviceproperties"
```


### Client Initialization

```go
client := tableserviceproperties.NewTableServicePropertiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TableServicePropertiesClient.TableServicesGetServiceProperties`

```go
ctx := context.TODO()
id := tableserviceproperties.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.TableServicesGetServiceProperties(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TableServicePropertiesClient.TableServicesList`

```go
ctx := context.TODO()
id := tableserviceproperties.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.TableServicesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TableServicePropertiesClient.TableServicesSetServiceProperties`

```go
ctx := context.TODO()
id := tableserviceproperties.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

payload := tableserviceproperties.TableServiceProperties{
	// ...
}


read, err := client.TableServicesSetServiceProperties(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
