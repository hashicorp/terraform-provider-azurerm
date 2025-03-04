
## `github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource` Documentation

The `serviceresource` SDK allows for interaction with Azure Resource Manager `datamigration` (API Version `2021-06-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
```


### Client Initialization

```go
client := serviceresource.NewServiceResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServiceResourceClient.ServiceTasksList`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

// alternatively `client.ServiceTasksList(ctx, id, serviceresource.DefaultServiceTasksListOperationOptions())` can be used to do batched pagination
items, err := client.ServiceTasksListComplete(ctx, id, serviceresource.DefaultServiceTasksListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServiceResourceClient.ServicesCheckStatus`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

read, err := client.ServicesCheckStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServiceResourceClient.ServicesCreateOrUpdate`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

payload := serviceresource.DataMigrationService{
	// ...
}


if err := client.ServicesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceResourceClient.ServicesDelete`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

if err := client.ServicesDeleteThenPoll(ctx, id, serviceresource.DefaultServicesDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceResourceClient.ServicesGet`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

read, err := client.ServicesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServiceResourceClient.ServicesList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ServicesList(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServiceResourceClient.ServicesListByResourceGroup`

```go
ctx := context.TODO()
id := serviceresource.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "resourceGroupName")

// alternatively `client.ServicesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServiceResourceClient.ServicesListSkus`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

// alternatively `client.ServicesListSkus(ctx, id)` can be used to do batched pagination
items, err := client.ServicesListSkusComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ServiceResourceClient.ServicesStart`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

if err := client.ServicesStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceResourceClient.ServicesStop`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

if err := client.ServicesStopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceResourceClient.ServicesUpdate`

```go
ctx := context.TODO()
id := serviceresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName")

payload := serviceresource.DataMigrationService{
	// ...
}


if err := client.ServicesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceResourceClient.TasksList`

```go
ctx := context.TODO()
id := serviceresource.NewProjectID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "serviceName", "projectName")

// alternatively `client.TasksList(ctx, id, serviceresource.DefaultTasksListOperationOptions())` can be used to do batched pagination
items, err := client.TasksListComplete(ctx, id, serviceresource.DefaultTasksListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
