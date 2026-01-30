
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratorsmicrosoftentra` Documentation

The `administratorsmicrosoftentra` SDK allows for interaction with Azure Resource Manager `postgresql` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratorsmicrosoftentra"
```


### Client Initialization

```go
client := administratorsmicrosoftentra.NewAdministratorsMicrosoftEntraClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AdministratorsMicrosoftEntraClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := administratorsmicrosoftentra.NewAdministratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "objectId")

payload := administratorsmicrosoftentra.AdministratorMicrosoftEntraAdd{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AdministratorsMicrosoftEntraClient.Delete`

```go
ctx := context.TODO()
id := administratorsmicrosoftentra.NewAdministratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "objectId")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AdministratorsMicrosoftEntraClient.Get`

```go
ctx := context.TODO()
id := administratorsmicrosoftentra.NewAdministratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "objectId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdministratorsMicrosoftEntraClient.ListByServer`

```go
ctx := context.TODO()
id := administratorsmicrosoftentra.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
