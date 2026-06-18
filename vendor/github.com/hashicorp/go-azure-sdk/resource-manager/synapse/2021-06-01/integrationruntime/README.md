
## `github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntime` Documentation

The `integrationruntime` SDK allows for interaction with Azure Resource Manager `synapse` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntime"
```


### Client Initialization

```go
client := integrationruntime.NewIntegrationRuntimeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationRuntimeClient.AuthKeysList`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.AuthKeysList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.AuthKeysRegenerate`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

payload := integrationruntime.IntegrationRuntimeRegenerateKeyParameters{
	// ...
}


read, err := client.AuthKeysRegenerate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.ConnectionInfosGet`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.ConnectionInfosGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.Create`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

payload := integrationruntime.IntegrationRuntimeResource{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload, integrationruntime.DefaultCreateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.CredentialsSync`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.CredentialsSync(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.Delete`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.DisableInteractiveQuery`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

if err := client.DisableInteractiveQueryThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.EnableInteractiveQuery`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

if err := client.EnableInteractiveQueryThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.Get`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.Get(ctx, id, integrationruntime.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.ListByWorkspace`

```go
ctx := context.TODO()
id := integrationruntime.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationRuntimeClient.MonitoringDataList`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.MonitoringDataList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.NodeIPAddressGet`

```go
ctx := context.TODO()
id := integrationruntime.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName", "nodeName")

read, err := client.NodeIPAddressGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.NodesDelete`

```go
ctx := context.TODO()
id := integrationruntime.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName", "nodeName")

read, err := client.NodesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.NodesGet`

```go
ctx := context.TODO()
id := integrationruntime.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName", "nodeName")

read, err := client.NodesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.NodesUpdate`

```go
ctx := context.TODO()
id := integrationruntime.NewNodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName", "nodeName")

payload := integrationruntime.UpdateIntegrationRuntimeNodeRequest{
	// ...
}


read, err := client.NodesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.ObjectMetadataList`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

payload := integrationruntime.GetSsisObjectMetadataRequest{
	// ...
}


// alternatively `client.ObjectMetadataList(ctx, id, payload)` can be used to do batched pagination
items, err := client.ObjectMetadataListComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationRuntimeClient.ObjectMetadataRefresh`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

if err := client.ObjectMetadataRefreshThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.Start`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.StatusGet`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.StatusGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimeClient.Stop`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimeClient.Update`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

payload := integrationruntime.UpdateIntegrationRuntimeRequest{
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


### Example Usage: `IntegrationRuntimeClient.Upgrade`

```go
ctx := context.TODO()
id := integrationruntime.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.Upgrade(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
