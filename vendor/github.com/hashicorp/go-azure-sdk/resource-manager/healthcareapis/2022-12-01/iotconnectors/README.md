
## `github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors` Documentation

The `iotconnectors` SDK allows for interaction with the Azure Resource Manager Service `healthcareapis` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
```


### Client Initialization

```go
client := iotconnectors.NewIotConnectorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IotConnectorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := iotconnectors.NewIotConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue")

payload := iotconnectors.IotConnector{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IotConnectorsClient.Delete`

```go
ctx := context.TODO()
id := iotconnectors.NewIotConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IotConnectorsClient.FhirDestinationsListByIotConnector`

```go
ctx := context.TODO()
id := iotconnectors.NewIotConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue")

// alternatively `client.FhirDestinationsListByIotConnector(ctx, id)` can be used to do batched pagination
items, err := client.FhirDestinationsListByIotConnectorComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IotConnectorsClient.Get`

```go
ctx := context.TODO()
id := iotconnectors.NewIotConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotConnectorsClient.IotConnectorFhirDestinationCreateOrUpdate`

```go
ctx := context.TODO()
id := iotconnectors.NewFhirDestinationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue", "fhirDestinationValue")

payload := iotconnectors.IotFhirDestination{
	// ...
}


if err := client.IotConnectorFhirDestinationCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IotConnectorsClient.IotConnectorFhirDestinationDelete`

```go
ctx := context.TODO()
id := iotconnectors.NewFhirDestinationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue", "fhirDestinationValue")

if err := client.IotConnectorFhirDestinationDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IotConnectorsClient.IotConnectorFhirDestinationGet`

```go
ctx := context.TODO()
id := iotconnectors.NewFhirDestinationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue", "fhirDestinationValue")

read, err := client.IotConnectorFhirDestinationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IotConnectorsClient.ListByWorkspace`

```go
ctx := context.TODO()
id := iotconnectors.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IotConnectorsClient.Update`

```go
ctx := context.TODO()
id := iotconnectors.NewIotConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "iotConnectorValue")

payload := iotconnectors.IotConnectorPatchResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
