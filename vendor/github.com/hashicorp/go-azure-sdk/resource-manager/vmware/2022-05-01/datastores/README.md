
## `github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores` Documentation

The `datastores` SDK allows for interaction with Azure Resource Manager `vmware` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/datastores"
```


### Client Initialization

```go
client := datastores.NewDataStoresClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataStoresClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datastores.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudName", "clusterName", "dataStoreName")

payload := datastores.Datastore{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataStoresClient.Delete`

```go
ctx := context.TODO()
id := datastores.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudName", "clusterName", "dataStoreName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataStoresClient.Get`

```go
ctx := context.TODO()
id := datastores.NewDataStoreID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudName", "clusterName", "dataStoreName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataStoresClient.List`

```go
ctx := context.TODO()
id := datastores.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudName", "clusterName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
