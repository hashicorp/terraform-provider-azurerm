
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/startstopmanagedinstanceschedules` Documentation

The `startstopmanagedinstanceschedules` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/startstopmanagedinstanceschedules"
```


### Client Initialization

```go
client := startstopmanagedinstanceschedules.NewStartStopManagedInstanceSchedulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StartStopManagedInstanceSchedulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName")

payload := startstopmanagedinstanceschedules.StartStopManagedInstanceSchedule{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StartStopManagedInstanceSchedulesClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StartStopManagedInstanceSchedulesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StartStopManagedInstanceSchedulesClient.ListByInstance`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName")

// alternatively `client.ListByInstance(ctx, id)` can be used to do batched pagination
items, err := client.ListByInstanceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
