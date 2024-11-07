
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountsessions` Documentation

The `integrationaccountsessions` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountsessions"
```


### Client Initialization

```go
client := integrationaccountsessions.NewIntegrationAccountSessionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountSessionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountsessions.NewSessionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "sessionName")

payload := integrationaccountsessions.IntegrationAccountSession{
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


### Example Usage: `IntegrationAccountSessionsClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountsessions.NewSessionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "sessionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountSessionsClient.Get`

```go
ctx := context.TODO()
id := integrationaccountsessions.NewSessionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "sessionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountSessionsClient.List`

```go
ctx := context.TODO()
id := integrationaccountsessions.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName")

// alternatively `client.List(ctx, id, integrationaccountsessions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, integrationaccountsessions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
