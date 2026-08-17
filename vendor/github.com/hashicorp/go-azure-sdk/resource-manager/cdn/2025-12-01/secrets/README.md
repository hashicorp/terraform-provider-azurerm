
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/secrets` Documentation

The `secrets` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/secrets"
```


### Client Initialization

```go
client := secrets.NewSecretsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecretsClient.Create`

```go
ctx := context.TODO()
id := secrets.NewSecretID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "secretName")

payload := secrets.Secret{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecretsClient.Delete`

```go
ctx := context.TODO()
id := secrets.NewSecretID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "secretName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SecretsClient.Get`

```go
ctx := context.TODO()
id := secrets.NewSecretID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "secretName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecretsClient.ListByProfile`

```go
ctx := context.TODO()
id := secrets.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
