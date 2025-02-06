
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions` Documentation

The `extensions` SDK allows for interaction with Azure Resource Manager `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
```


### Client Initialization

```go
client := extensions.NewExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExtensionsClient.Create`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "extensionName")

payload := extensions.Extension{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Delete`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "extensionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.DisableAzureMonitor`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.DisableAzureMonitorThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.DisableMonitoring`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

if err := client.DisableMonitoringThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.EnableAzureMonitor`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := extensions.AzureMonitorRequest{
	// ...
}


if err := client.EnableAzureMonitorThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.EnableMonitoring`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := extensions.ClusterMonitoringRequest{
	// ...
}


if err := client.EnableMonitoringThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Get`

```go
ctx := context.TODO()
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "extensionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExtensionsClient.GetAzureMonitorStatus`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.GetAzureMonitorStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExtensionsClient.GetMonitoringStatus`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.GetMonitoringStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
