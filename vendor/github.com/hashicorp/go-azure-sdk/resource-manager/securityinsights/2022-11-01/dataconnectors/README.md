
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/dataconnectors` Documentation

The `dataconnectors` SDK allows for interaction with the Azure Resource Manager Service `securityinsights` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/dataconnectors"
```


### Client Initialization

```go
client := dataconnectors.NewDataConnectorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataConnectorsClient.DataConnectorsCreateOrUpdate`

```go
ctx := context.TODO()
id := dataconnectors.NewDataConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataConnectorIdValue")

payload := dataconnectors.DataConnector{
	// ...
}


read, err := client.DataConnectorsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataConnectorsClient.DataConnectorsDelete`

```go
ctx := context.TODO()
id := dataconnectors.NewDataConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataConnectorIdValue")

read, err := client.DataConnectorsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataConnectorsClient.DataConnectorsGet`

```go
ctx := context.TODO()
id := dataconnectors.NewDataConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataConnectorIdValue")

read, err := client.DataConnectorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataConnectorsClient.DataConnectorsList`

```go
ctx := context.TODO()
id := dataconnectors.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.DataConnectorsList(ctx, id)` can be used to do batched pagination
items, err := client.DataConnectorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
