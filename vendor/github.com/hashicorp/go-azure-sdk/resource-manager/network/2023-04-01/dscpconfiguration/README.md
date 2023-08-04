
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/dscpconfiguration` Documentation

The `dscpconfiguration` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/dscpconfiguration"
```


### Client Initialization

```go
client := dscpconfiguration.NewDscpConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DscpConfigurationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dscpconfiguration.NewDscpConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dscpConfigurationValue")

payload := dscpconfiguration.DscpConfiguration{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DscpConfigurationClient.Delete`

```go
ctx := context.TODO()
id := dscpconfiguration.NewDscpConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dscpConfigurationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DscpConfigurationClient.Get`

```go
ctx := context.TODO()
id := dscpconfiguration.NewDscpConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dscpConfigurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
