
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratormicrosoftentras` Documentation

The `administratormicrosoftentras` SDK allows for interaction with Azure Resource Manager `postgresql` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratormicrosoftentras"
```


### Client Initialization

```go
client := administratormicrosoftentras.NewAdministratorMicrosoftEntrasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AdministratorMicrosoftEntrasClient.AdministratorsMicrosoftEntraCreateOrUpdate`

```go
ctx := context.TODO()
id := administratormicrosoftentras.NewAdministratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "objectId")

payload := administratormicrosoftentras.AdministratorMicrosoftEntraAdd{
	// ...
}


if err := client.AdministratorsMicrosoftEntraCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AdministratorMicrosoftEntrasClient.AdministratorsMicrosoftEntraDelete`

```go
ctx := context.TODO()
id := administratormicrosoftentras.NewAdministratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "objectId")

if err := client.AdministratorsMicrosoftEntraDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AdministratorMicrosoftEntrasClient.AdministratorsMicrosoftEntraGet`

```go
ctx := context.TODO()
id := administratormicrosoftentras.NewAdministratorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "objectId")

read, err := client.AdministratorsMicrosoftEntraGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdministratorMicrosoftEntrasClient.AdministratorsMicrosoftEntraListByServer`

```go
ctx := context.TODO()
id := administratormicrosoftentras.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.AdministratorsMicrosoftEntraListByServer(ctx, id)` can be used to do batched pagination
items, err := client.AdministratorsMicrosoftEntraListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
