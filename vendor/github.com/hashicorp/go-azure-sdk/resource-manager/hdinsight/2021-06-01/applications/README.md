
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications` Documentation

The `applications` SDK allows for interaction with Azure Resource Manager `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
```


### Client Initialization

```go
client := applications.NewApplicationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationsClient.Create`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "applicationName")

payload := applications.Application{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.Delete`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "applicationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.Get`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "applicationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.ListByCluster`

```go
ctx := context.TODO()
id := commonids.NewHDInsightClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
