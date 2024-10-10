
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountmaps` Documentation

The `integrationaccountmaps` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountmaps"
```


### Client Initialization

```go
client := integrationaccountmaps.NewIntegrationAccountMapsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountMapsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountmaps.NewMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "mapName")

payload := integrationaccountmaps.IntegrationAccountMap{
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


### Example Usage: `IntegrationAccountMapsClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountmaps.NewMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "mapName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountMapsClient.Get`

```go
ctx := context.TODO()
id := integrationaccountmaps.NewMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "mapName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountMapsClient.List`

```go
ctx := context.TODO()
id := integrationaccountmaps.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName")

// alternatively `client.List(ctx, id, integrationaccountmaps.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, integrationaccountmaps.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationAccountMapsClient.ListContentCallbackURL`

```go
ctx := context.TODO()
id := integrationaccountmaps.NewMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "mapName")

payload := integrationaccountmaps.GetCallbackURLParameters{
	// ...
}


read, err := client.ListContentCallbackURL(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
