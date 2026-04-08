
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/verifierworkspaces` Documentation

The `verifierworkspaces` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/verifierworkspaces"
```


### Client Initialization

```go
client := verifierworkspaces.NewVerifierWorkspacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VerifierWorkspacesClient.Create`

```go
ctx := context.TODO()
id := verifierworkspaces.NewVerifierWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName")

payload := verifierworkspaces.VerifierWorkspace{
	// ...
}


read, err := client.Create(ctx, id, payload, verifierworkspaces.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VerifierWorkspacesClient.Delete`

```go
ctx := context.TODO()
id := verifierworkspaces.NewVerifierWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName")

if err := client.DeleteThenPoll(ctx, id, verifierworkspaces.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VerifierWorkspacesClient.Get`

```go
ctx := context.TODO()
id := verifierworkspaces.NewVerifierWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VerifierWorkspacesClient.List`

```go
ctx := context.TODO()
id := verifierworkspaces.NewNetworkManagerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName")

// alternatively `client.List(ctx, id, verifierworkspaces.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, verifierworkspaces.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VerifierWorkspacesClient.Update`

```go
ctx := context.TODO()
id := verifierworkspaces.NewVerifierWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "verifierWorkspaceName")

payload := verifierworkspaces.VerifierWorkspaceUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload, verifierworkspaces.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
