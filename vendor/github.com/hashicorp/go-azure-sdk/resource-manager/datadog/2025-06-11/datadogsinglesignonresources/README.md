
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogsinglesignonresources` Documentation

The `datadogsinglesignonresources` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2025-06-11`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogsinglesignonresources"
```


### Client Initialization

```go
client := datadogsinglesignonresources.NewDatadogSingleSignOnResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatadogSingleSignOnResourcesClient.SingleSignOnConfigurationsCreateOrUpdate`

```go
ctx := context.TODO()
id := datadogsinglesignonresources.NewSingleSignOnConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "singleSignOnConfigurationName")

payload := datadogsinglesignonresources.DatadogSingleSignOnResource{
	// ...
}


if err := client.SingleSignOnConfigurationsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatadogSingleSignOnResourcesClient.SingleSignOnConfigurationsGet`

```go
ctx := context.TODO()
id := datadogsinglesignonresources.NewSingleSignOnConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "singleSignOnConfigurationName")

read, err := client.SingleSignOnConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatadogSingleSignOnResourcesClient.SingleSignOnConfigurationsList`

```go
ctx := context.TODO()
id := datadogsinglesignonresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.SingleSignOnConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.SingleSignOnConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
