
## `github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools` Documentation

The `bigdatapools` SDK allows for interaction with Azure Resource Manager `synapse` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
```


### Client Initialization

```go
client := bigdatapools.NewBigDataPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BigDataPoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := bigdatapools.NewBigDataPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "bigDataPoolName")

payload := bigdatapools.BigDataPoolResourceInfo{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, bigdatapools.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BigDataPoolsClient.Delete`

```go
ctx := context.TODO()
id := bigdatapools.NewBigDataPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "bigDataPoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BigDataPoolsClient.Get`

```go
ctx := context.TODO()
id := bigdatapools.NewBigDataPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "bigDataPoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BigDataPoolsClient.ListByWorkspace`

```go
ctx := context.TODO()
id := bigdatapools.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BigDataPoolsClient.Update`

```go
ctx := context.TODO()
id := bigdatapools.NewBigDataPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "bigDataPoolName")

payload := bigdatapools.BigDataPoolPatchInfo{
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
