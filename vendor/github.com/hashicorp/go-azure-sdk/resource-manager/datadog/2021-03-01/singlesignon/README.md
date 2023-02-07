
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon` Documentation

The `singlesignon` SDK allows for interaction with the Azure Resource Manager Service `datadog` (API Version `2021-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
```


### Client Initialization

```go
client := singlesignon.NewSingleSignOnClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SingleSignOnClient.ConfigurationsCreateOrUpdate`

```go
ctx := context.TODO()
id := singlesignon.NewSingleSignOnConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "singleSignOnConfigurationValue")

payload := singlesignon.DatadogSingleSignOnResource{
	// ...
}


if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SingleSignOnClient.ConfigurationsGet`

```go
ctx := context.TODO()
id := singlesignon.NewSingleSignOnConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "singleSignOnConfigurationValue")

read, err := client.ConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SingleSignOnClient.ConfigurationsList`

```go
ctx := context.TODO()
id := singlesignon.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

// alternatively `client.ConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.ConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
