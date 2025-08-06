
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations` Documentation

The `integrationaccountbatchconfigurations` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
```


### Client Initialization

```go
client := integrationaccountbatchconfigurations.NewIntegrationAccountBatchConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountBatchConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountbatchconfigurations.NewBatchConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "batchConfigurationName")

payload := integrationaccountbatchconfigurations.BatchConfiguration{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountBatchConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountbatchconfigurations.NewBatchConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "batchConfigurationName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountBatchConfigurationsClient.Get`

```go
ctx := context.TODO()
id := integrationaccountbatchconfigurations.NewBatchConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "batchConfigurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountBatchConfigurationsClient.List`

```go
ctx := context.TODO()
id := integrationaccountbatchconfigurations.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
