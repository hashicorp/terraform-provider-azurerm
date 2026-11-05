
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/commits` Documentation

The `commits` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/commits"
```


### Client Initialization

```go
client := commits.NewCommitsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CommitsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commits.NewCommitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "commitName")

payload := commits.Commit{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CommitsClient.Delete`

```go
ctx := context.TODO()
id := commits.NewCommitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "commitName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CommitsClient.Get`

```go
ctx := context.TODO()
id := commits.NewCommitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "commitName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommitsClient.List`

```go
ctx := context.TODO()
id := commits.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName")

// alternatively `client.List(ctx, id, commits.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, commits.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
