
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets` Documentation

The `credentialsets` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
```


### Client Initialization

```go
client := credentialsets.NewCredentialSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CredentialSetsClient.Create`

```go
ctx := context.TODO()
id := credentialsets.NewCredentialSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "credentialSetName")

payload := credentialsets.CredentialSet{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CredentialSetsClient.Delete`

```go
ctx := context.TODO()
id := credentialsets.NewCredentialSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "credentialSetName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CredentialSetsClient.Get`

```go
ctx := context.TODO()
id := credentialsets.NewCredentialSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "credentialSetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CredentialSetsClient.List`

```go
ctx := context.TODO()
id := credentialsets.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CredentialSetsClient.Update`

```go
ctx := context.TODO()
id := credentialsets.NewCredentialSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "credentialSetName")

payload := credentialsets.CredentialSetUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
