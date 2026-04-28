
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources` Documentation

The `metricsobjectfirewallresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
```


### Client Initialization

```go
client := metricsobjectfirewallresources.NewMetricsObjectFirewallResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MetricsObjectFirewallResourcesClient.MetricsObjectFirewallCreateOrUpdate`

```go
ctx := context.TODO()
id := metricsobjectfirewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := metricsobjectfirewallresources.MetricsObjectFirewallResource{
	// ...
}


if err := client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MetricsObjectFirewallResourcesClient.MetricsObjectFirewallDelete`

```go
ctx := context.TODO()
id := metricsobjectfirewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

if err := client.MetricsObjectFirewallDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MetricsObjectFirewallResourcesClient.MetricsObjectFirewallGet`

```go
ctx := context.TODO()
id := metricsobjectfirewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.MetricsObjectFirewallGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetricsObjectFirewallResourcesClient.MetricsObjectFirewallListByFirewalls`

```go
ctx := context.TODO()
id := metricsobjectfirewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

// alternatively `client.MetricsObjectFirewallListByFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.MetricsObjectFirewallListByFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
