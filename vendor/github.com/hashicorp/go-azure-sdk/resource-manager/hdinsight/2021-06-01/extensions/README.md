
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions` Documentation

The `extensions` SDK allows for interaction with the Azure Resource Manager Service `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
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
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "extensionValue")

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
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "extensionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.DisableAzureMonitor`

```go
ctx := context.TODO()
id := extensions.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.DisableAzureMonitorThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.DisableMonitoring`

```go
ctx := context.TODO()
id := extensions.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.DisableMonitoringThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.EnableAzureMonitor`

```go
ctx := context.TODO()
id := extensions.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

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
id := extensions.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

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
id := extensions.NewExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "extensionValue")

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
id := extensions.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

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
id := extensions.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.GetMonitoringStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
