
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances` Documentation

The `sapapplicationserverinstances` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"
```


### Client Initialization

```go
client := sapapplicationserverinstances.NewSAPApplicationServerInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPApplicationServerInstancesClient.Create`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewApplicationInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "applicationInstanceName")

payload := sapapplicationserverinstances.SAPApplicationServerInstance{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPApplicationServerInstancesClient.Delete`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewApplicationInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "applicationInstanceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPApplicationServerInstancesClient.Get`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewApplicationInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "applicationInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPApplicationServerInstancesClient.List`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SAPApplicationServerInstancesClient.StartInstance`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewApplicationInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "applicationInstanceName")

if err := client.StartInstanceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPApplicationServerInstancesClient.StopInstance`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewApplicationInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "applicationInstanceName")

payload := sapapplicationserverinstances.StopRequest{
	// ...
}


if err := client.StopInstanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPApplicationServerInstancesClient.Update`

```go
ctx := context.TODO()
id := sapapplicationserverinstances.NewApplicationInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "applicationInstanceName")

payload := sapapplicationserverinstances.UpdateSAPApplicationInstanceRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
