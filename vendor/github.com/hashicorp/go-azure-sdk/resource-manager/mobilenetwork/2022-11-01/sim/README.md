
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim` Documentation

The `sim` SDK allows for interaction with Azure Resource Manager `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/sim"
```


### Client Initialization

```go
client := sim.NewSIMClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SIMClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := sim.NewSimID("12345678-1234-9876-4563-123456789012", "example-resource-group", "simGroupName", "simName")

payload := sim.Sim{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SIMClient.Delete`

```go
ctx := context.TODO()
id := sim.NewSimID("12345678-1234-9876-4563-123456789012", "example-resource-group", "simGroupName", "simName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SIMClient.Get`

```go
ctx := context.TODO()
id := sim.NewSimID("12345678-1234-9876-4563-123456789012", "example-resource-group", "simGroupName", "simName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
