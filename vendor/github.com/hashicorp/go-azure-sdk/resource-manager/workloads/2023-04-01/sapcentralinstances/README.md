
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapcentralinstances` Documentation

The `sapcentralinstances` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapcentralinstances"
```


### Client Initialization

```go
client := sapcentralinstances.NewSAPCentralInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPCentralInstancesClient.Create`

```go
ctx := context.TODO()
id := sapcentralinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralinstances.SAPCentralServerInstance{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralInstancesClient.Delete`

```go
ctx := context.TODO()
id := sapcentralinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralInstancesClient.Get`

```go
ctx := context.TODO()
id := sapcentralinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPCentralInstancesClient.List`

```go
ctx := context.TODO()
id := sapcentralinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SAPCentralInstancesClient.StartInstance`

```go
ctx := context.TODO()
id := sapcentralinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

if err := client.StartInstanceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralInstancesClient.StopInstance`

```go
ctx := context.TODO()
id := sapcentralinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralinstances.StopRequest{
	// ...
}


if err := client.StopInstanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPCentralInstancesClient.Update`

```go
ctx := context.TODO()
id := sapcentralinstances.NewCentralInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "centralInstanceName")

payload := sapcentralinstances.UpdateSAPCentralInstanceRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
