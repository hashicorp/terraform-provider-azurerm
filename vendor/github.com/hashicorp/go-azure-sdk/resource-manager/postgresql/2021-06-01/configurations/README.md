
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/configurations` Documentation

The `configurations` SDK allows for interaction with the Azure Resource Manager Service `postgresql` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/configurations"
```


### Client Initialization

```go
client := configurations.NewConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationsClient.Get`

```go
ctx := context.TODO()
id := configurations.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerValue", "configurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationsClient.ListByServer`

```go
ctx := context.TODO()
id := configurations.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerValue")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationsClient.Put`

```go
ctx := context.TODO()
id := configurations.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerValue", "configurationValue")

payload := configurations.Configuration{
	// ...
}


if err := client.PutThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConfigurationsClient.Update`

```go
ctx := context.TODO()
id := configurations.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerValue", "configurationValue")

payload := configurations.Configuration{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
