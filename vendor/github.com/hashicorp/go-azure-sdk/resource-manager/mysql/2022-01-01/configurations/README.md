
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/configurations` Documentation

The `configurations` SDK allows for interaction with the Azure Resource Manager Service `mysql` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/configurations"
```


### Client Initialization

```go
client := configurations.NewConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationsClient.BatchUpdate`

```go
ctx := context.TODO()
id := configurations.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerValue")

payload := configurations.ConfigurationListForBatchUpdate{
	// ...
}


// alternatively `client.BatchUpdate(ctx, id, payload)` can be used to do batched pagination
items, err := client.BatchUpdateComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurations.NewConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerValue", "configurationValue")

payload := configurations.Configuration{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
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

// alternatively `client.ListByServer(ctx, id, configurations.DefaultListByServerOperationOptions())` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id, configurations.DefaultListByServerOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
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
