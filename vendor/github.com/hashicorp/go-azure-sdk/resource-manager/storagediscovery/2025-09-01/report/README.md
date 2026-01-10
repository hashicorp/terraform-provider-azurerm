
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/report` Documentation

The `report` SDK allows for interaction with Azure Resource Manager `storagediscovery` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/report"
```


### Client Initialization

```go
client := report.NewReportClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReportClient.GenerateReport`

```go
ctx := context.TODO()
id := report.NewReportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageDiscoveryWorkspaceName", "reportName")

payload := report.GetReportContent{
	// ...
}


if err := client.GenerateReportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReportClient.Get`

```go
ctx := context.TODO()
id := report.NewReportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageDiscoveryWorkspaceName", "reportName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReportClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := report.NewProviderStorageDiscoveryWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageDiscoveryWorkspaceName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReportClient.ListBySubscription`

```go
ctx := context.TODO()
id := report.NewStorageDiscoveryWorkspaceID("12345678-1234-9876-4563-123456789012", "storageDiscoveryWorkspaceName")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
