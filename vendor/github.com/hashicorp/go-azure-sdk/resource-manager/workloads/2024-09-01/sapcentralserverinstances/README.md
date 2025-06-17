
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapcentralserverinstances` Documentation

The `sapcentralserverinstances` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2024-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapcentralserverinstances"
```


### Client Initialization

```go
client := sapcentralserverinstances.NewSAPCentralServerInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPCentralServerInstancesClient.Create`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralserverinstances.SAPCentralServerInstance{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralServerInstancesClient.Delete`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralServerInstancesClient.Get`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPCentralServerInstancesClient.List`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SAPCentralServerInstancesClient.Start`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralserverinstances.StartRequest{
	// ...
}


if err := client.StartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralServerInstancesClient.Stop`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralserverinstances.StopRequest{
	// ...
}


if err := client.StopThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralServerInstancesClient.Update`

```go
ctx := context.TODO()
id := sapcentralserverinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralserverinstances.UpdateSAPCentralInstanceRequest{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
