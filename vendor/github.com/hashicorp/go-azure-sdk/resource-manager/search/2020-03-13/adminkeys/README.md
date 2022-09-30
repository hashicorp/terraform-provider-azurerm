
## `github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/adminkeys` Documentation

The `adminkeys` SDK allows for interaction with the Azure Resource Manager Service `search` (API Version `2020-03-13`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/adminkeys"
```


### Client Initialization

```go
client := adminkeys.NewAdminKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AdminKeysClient.Get`

```go
ctx := context.TODO()
id := adminkeys.NewSearchServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue")

read, err := client.Get(ctx, id, adminkeys.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdminKeysClient.Regenerate`

```go
ctx := context.TODO()
id := adminkeys.NewKeyKindID("12345678-1234-9876-4563-123456789012", "example-resource-group", "searchServiceValue", "primary")

read, err := client.Regenerate(ctx, id, adminkeys.DefaultRegenerateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
