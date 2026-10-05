
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations` Documentation

The `configurations` SDK allows for interaction with Azure Resource Manager `postgresqlhsc` (API Version `2022-11-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
```


### Client Initialization

```go
client := configurations.NewConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationsClient.Get`

```go
ctx := context.TODO()
id := configurations.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name", "configurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationsClient.GetCoordinator`

```go
ctx := context.TODO()
id := configurations.NewCoordinatorConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name", "coordinatorConfigurationName")

read, err := client.GetCoordinator(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationsClient.GetNode`

```go
ctx := context.TODO()
id := configurations.NewNodeConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name", "nodeConfigurationName")

read, err := client.GetNode(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationsClient.ListByCluster`

```go
ctx := context.TODO()
id := configurations.NewServerGroupsv2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationsClient.ListByServer`

```go
ctx := context.TODO()
id := configurations.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name", "serverName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationsClient.UpdateOnCoordinator`

```go
ctx := context.TODO()
id := configurations.NewCoordinatorConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name", "coordinatorConfigurationName")

payload := configurations.ServerConfiguration{
	// ...
}


if err := client.UpdateOnCoordinatorThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConfigurationsClient.UpdateOnNode`

```go
ctx := context.TODO()
id := configurations.NewNodeConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverGroupsv2Name", "nodeConfigurationName")

payload := configurations.ServerConfiguration{
	// ...
}


if err := client.UpdateOnNodeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
