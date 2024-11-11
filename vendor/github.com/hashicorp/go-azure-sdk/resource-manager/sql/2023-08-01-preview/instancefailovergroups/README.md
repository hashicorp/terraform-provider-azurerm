
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/instancefailovergroups` Documentation

The `instancefailovergroups` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/instancefailovergroups"
```


### Client Initialization

```go
client := instancefailovergroups.NewInstanceFailoverGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `InstanceFailoverGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := instancefailovergroups.NewInstanceFailoverGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "instanceFailoverGroupName")

payload := instancefailovergroups.InstanceFailoverGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `InstanceFailoverGroupsClient.Delete`

```go
ctx := context.TODO()
id := instancefailovergroups.NewInstanceFailoverGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "instanceFailoverGroupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `InstanceFailoverGroupsClient.Failover`

```go
ctx := context.TODO()
id := instancefailovergroups.NewInstanceFailoverGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "instanceFailoverGroupName")

if err := client.FailoverThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `InstanceFailoverGroupsClient.ForceFailoverAllowDataLoss`

```go
ctx := context.TODO()
id := instancefailovergroups.NewInstanceFailoverGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "instanceFailoverGroupName")

if err := client.ForceFailoverAllowDataLossThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `InstanceFailoverGroupsClient.Get`

```go
ctx := context.TODO()
id := instancefailovergroups.NewInstanceFailoverGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName", "instanceFailoverGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InstanceFailoverGroupsClient.ListByLocation`

```go
ctx := context.TODO()
id := instancefailovergroups.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
