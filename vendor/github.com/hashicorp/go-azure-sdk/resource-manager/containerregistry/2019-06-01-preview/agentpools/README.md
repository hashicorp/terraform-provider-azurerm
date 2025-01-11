
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/agentpools` Documentation

The `agentpools` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2019-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/agentpools"
```


### Client Initialization

```go
client := agentpools.NewAgentPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AgentPoolsClient.Create`

```go
ctx := context.TODO()
id := agentpools.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "agentPoolName")

payload := agentpools.AgentPool{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AgentPoolsClient.Delete`

```go
ctx := context.TODO()
id := agentpools.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "agentPoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AgentPoolsClient.Get`

```go
ctx := context.TODO()
id := agentpools.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "agentPoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgentPoolsClient.GetQueueStatus`

```go
ctx := context.TODO()
id := agentpools.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "agentPoolName")

read, err := client.GetQueueStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AgentPoolsClient.List`

```go
ctx := context.TODO()
id := agentpools.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AgentPoolsClient.Update`

```go
ctx := context.TODO()
id := agentpools.NewAgentPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "agentPoolName")

payload := agentpools.AgentPoolUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
