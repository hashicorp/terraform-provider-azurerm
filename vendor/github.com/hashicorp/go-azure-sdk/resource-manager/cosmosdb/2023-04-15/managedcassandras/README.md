
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras` Documentation

The `managedcassandras` SDK allows for interaction with Azure Resource Manager `cosmosdb` (API Version `2023-04-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
```


### Client Initialization

```go
client := managedcassandras.NewManagedCassandrasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersCreateUpdate`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

payload := managedcassandras.ClusterResource{
	// ...
}


if err := client.CassandraClustersCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersDeallocate`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

if err := client.CassandraClustersDeallocateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersDelete`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

if err := client.CassandraClustersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersGet`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

read, err := client.CassandraClustersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersInvokeCommand`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

payload := managedcassandras.CommandPostBody{
	// ...
}


if err := client.CassandraClustersInvokeCommandThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.CassandraClustersListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.CassandraClustersListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersStart`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

if err := client.CassandraClustersStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersStatus`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

read, err := client.CassandraClustersStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedCassandrasClient.CassandraClustersUpdate`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

payload := managedcassandras.ClusterResource{
	// ...
}


if err := client.CassandraClustersUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraDataCentersCreateUpdate`

```go
ctx := context.TODO()
id := managedcassandras.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

payload := managedcassandras.DataCenterResource{
	// ...
}


if err := client.CassandraDataCentersCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraDataCentersDelete`

```go
ctx := context.TODO()
id := managedcassandras.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

if err := client.CassandraDataCentersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedCassandrasClient.CassandraDataCentersGet`

```go
ctx := context.TODO()
id := managedcassandras.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

read, err := client.CassandraDataCentersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedCassandrasClient.CassandraDataCentersList`

```go
ctx := context.TODO()
id := managedcassandras.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

read, err := client.CassandraDataCentersList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedCassandrasClient.CassandraDataCentersUpdate`

```go
ctx := context.TODO()
id := managedcassandras.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

payload := managedcassandras.DataCenterResource{
	// ...
}


if err := client.CassandraDataCentersUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
