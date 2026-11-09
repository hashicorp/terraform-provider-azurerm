
## `github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes` Documentation

The `integrationruntimes` SDK allows for interaction with Azure Resource Manager `synapse` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/integrationruntimes"
```


### Client Initialization

```go
client := integrationruntimes.NewIntegrationRuntimesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "integrationRuntimeName")

read, err := client.ListOutboundNetworkDependenciesEndpoints(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
