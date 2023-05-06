
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/tokens` Documentation

The `tokens` SDK allows for interaction with the Azure Resource Manager Service `containerregistry` (API Version `2021-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/tokens"
```


### Client Initialization

```go
client := tokens.NewTokensClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TokensClient.Create`

```go
ctx := context.TODO()
id := tokens.NewTokenID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "tokenValue")

payload := tokens.Token{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TokensClient.Delete`

```go
ctx := context.TODO()
id := tokens.NewTokenID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "tokenValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TokensClient.Get`

```go
ctx := context.TODO()
id := tokens.NewTokenID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "tokenValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TokensClient.List`

```go
ctx := context.TODO()
id := tokens.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TokensClient.Update`

```go
ctx := context.TODO()
id := tokens.NewTokenID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "tokenValue")

payload := tokens.TokenUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
