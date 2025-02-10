
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/clusterprincipalassignments` Documentation

The `clusterprincipalassignments` SDK allows for interaction with Azure Resource Manager `kusto` (API Version `2023-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/clusterprincipalassignments"
```


### Client Initialization

```go
client := clusterprincipalassignments.NewClusterPrincipalAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ClusterPrincipalAssignmentsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := clusterprincipalassignments.ClusterPrincipalAssignmentCheckNameRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClusterPrincipalAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := clusterprincipalassignments.NewPrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "principalAssignmentName")

payload := clusterprincipalassignments.ClusterPrincipalAssignment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClusterPrincipalAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := clusterprincipalassignments.NewPrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "principalAssignmentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClusterPrincipalAssignmentsClient.Get`

```go
ctx := context.TODO()
id := clusterprincipalassignments.NewPrincipalAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "principalAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClusterPrincipalAssignmentsClient.List`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
