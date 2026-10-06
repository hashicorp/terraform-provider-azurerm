
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimeenableinteractivequery` Documentation

The `integrationruntimeenableinteractivequery` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimeenableinteractivequery"
```


### Client Initialization

```go
client := integrationruntimeenableinteractivequery.NewIntegrationRuntimeEnableInteractiveQueryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationRuntimeEnableInteractiveQueryClient.IntegrationRuntimeEnableInteractiveQuery`

```go
ctx := context.TODO()
id := integrationruntimeenableinteractivequery.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

payload := integrationruntimeenableinteractivequery.EnableInteractiveQueryRequest{
	// ...
}


if err := client.IntegrationRuntimeEnableInteractiveQueryThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
