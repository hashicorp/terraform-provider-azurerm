
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/configurationsupdate` Documentation

The `configurationsupdate` SDK allows for interaction with the Azure Resource Manager Service `mysql` (API Version `2017-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/configurationsupdate"
```


### Client Initialization

```go
client := configurationsupdate.NewConfigurationsUpdateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationsUpdateClient.ServerParametersListUpdateConfigurations`

```go
ctx := context.TODO()
id := configurationsupdate.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := configurationsupdate.ConfigurationListResult{
	// ...
}


if err := client.ServerParametersListUpdateConfigurationsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
