
## `github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers` Documentation

The `servers` SDK allows for interaction with the Azure Resource Manager Service `analysisservices` (API Version `2017-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
```


### Client Initialization

```go
client := servers.NewServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
if err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := servers.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := servers.CheckServerNameAvailabilityParameters{
	// ...
}

read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.Create`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.AnalysisServicesServer{
	// ...
}

future, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Delete`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
future, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.DissociateGateway`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
read, err := client.DissociateGateway(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.GetDetails`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
read, err := client.GetDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.List`

```go
ctx := context.TODO()
id := servers.NewSubscriptionID()
read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := servers.NewResourceGroupID()
read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ListGatewayStatus`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
read, err := client.ListGatewayStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.ListSkusForExisting`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
read, err := client.ListSkusForExisting(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServersClient.Resume`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
future, err := client.Resume(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Suspend`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")
future, err := client.Suspend(ctx, id)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `ServersClient.Update`

```go
ctx := context.TODO()
id := servers.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := servers.AnalysisServicesServerUpdateParameters{
	// ...
}

future, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```
