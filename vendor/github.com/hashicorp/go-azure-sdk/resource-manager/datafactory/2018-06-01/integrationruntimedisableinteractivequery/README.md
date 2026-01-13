
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimedisableinteractivequery` Documentation

The `integrationruntimedisableinteractivequery` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimedisableinteractivequery"
```


### Client Initialization

```go
client := integrationruntimedisableinteractivequery.NewIntegrationRuntimeDisableInteractiveQueryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationRuntimeDisableInteractiveQueryClient.IntegrationRuntimeDisableInteractiveQuery`

```go
ctx := context.TODO()
id := integrationruntimedisableinteractivequery.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

if err := client.IntegrationRuntimeDisableInteractiveQueryThenPoll(ctx, id); err != nil {
	// handle the error
}
```
