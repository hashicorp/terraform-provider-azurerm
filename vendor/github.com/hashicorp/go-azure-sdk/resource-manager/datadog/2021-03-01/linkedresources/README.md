
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/linkedresources` Documentation

The `linkedresources` SDK allows for interaction with the Azure Resource Manager Service `datadog` (API Version `2021-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/linkedresources"
```


### Client Initialization

```go
client := linkedresources.NewLinkedResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinkedResourcesClient.MonitorsListLinkedResources`

```go
ctx := context.TODO()
id := linkedresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

// alternatively `client.MonitorsListLinkedResources(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListLinkedResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
