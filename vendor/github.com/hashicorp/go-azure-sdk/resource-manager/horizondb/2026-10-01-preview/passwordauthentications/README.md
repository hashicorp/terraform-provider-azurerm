
## `github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/passwordauthentications` Documentation

The `passwordauthentications` SDK allows for interaction with Azure Resource Manager `horizondb` (API Version `2026-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/passwordauthentications"
```


### Client Initialization

```go
client := passwordauthentications.NewPasswordAuthenticationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PasswordAuthenticationsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := passwordauthentications.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := passwordauthentications.PasswordAuthentication{
	// ...
}


if err := client.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PasswordAuthenticationsClient.Get`

```go
ctx := context.TODO()
id := passwordauthentications.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PasswordAuthenticationsClient.Update`

```go
ctx := context.TODO()
id := passwordauthentications.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := passwordauthentications.PasswordAuthenticationForPatchUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
