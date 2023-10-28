
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadonlyauthentications` Documentation

The `serverazureadonlyauthentications` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-02-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadonlyauthentications"
```


### Client Initialization

```go
client := serverazureadonlyauthentications.NewServerAzureADOnlyAuthenticationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerAzureADOnlyAuthenticationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := serverazureadonlyauthentications.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := serverazureadonlyauthentications.ServerAzureADOnlyAuthentication{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServerAzureADOnlyAuthenticationsClient.Delete`

```go
ctx := context.TODO()
id := serverazureadonlyauthentications.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServerAzureADOnlyAuthenticationsClient.Get`

```go
ctx := context.TODO()
id := serverazureadonlyauthentications.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServerAzureADOnlyAuthenticationsClient.ListByServer`

```go
ctx := context.TODO()
id := serverazureadonlyauthentications.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
