
## `github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/connectorresources` Documentation

The `connectorresources` SDK allows for interaction with Azure Resource Manager `confluent` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/connectorresources"
```


### Client Initialization

```go
client := connectorresources.NewConnectorResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectorResourcesClient.ConnectorCreateOrUpdate`

```go
ctx := context.TODO()
id := connectorresources.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId", "connectorName")

payload := connectorresources.ConnectorResource{
	// ...
}


read, err := client.ConnectorCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectorResourcesClient.ConnectorDelete`

```go
ctx := context.TODO()
id := connectorresources.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId", "connectorName")

if err := client.ConnectorDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectorResourcesClient.ConnectorGet`

```go
ctx := context.TODO()
id := connectorresources.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId", "connectorName")

read, err := client.ConnectorGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectorResourcesClient.ConnectorList`

```go
ctx := context.TODO()
id := connectorresources.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId")

// alternatively `client.ConnectorList(ctx, id, connectorresources.DefaultConnectorListOperationOptions())` can be used to do batched pagination
items, err := client.ConnectorListComplete(ctx, id, connectorresources.DefaultConnectorListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
