
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databases` Documentation

The `databases` SDK allows for interaction with Azure Resource Manager `kusto` (API Version `2023-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databases"
```


### Client Initialization

```go
client := databases.NewDatabasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasesClient.AddPrincipals`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

payload := databases.DatabasePrincipalListRequest{
	// ...
}


read, err := client.AddPrincipals(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := databases.CheckNameRequest{
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


### Example Usage: `DatabasesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

payload := databases.Database{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, databases.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.DatabaseInviteFollower`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

payload := databases.DatabaseInviteFollowerRequest{
	// ...
}


read, err := client.DatabaseInviteFollower(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.ListByCluster`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.ListByCluster(ctx, id, databases.DefaultListByClusterOperationOptions())` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id, databases.DefaultListByClusterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DatabasesClient.ListPrincipals`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

read, err := client.ListPrincipals(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.RemovePrincipals`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

payload := databases.DatabasePrincipalListRequest{
	// ...
}


read, err := client.RemovePrincipals(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.Update`

```go
ctx := context.TODO()
id := commonids.NewKustoDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "databaseName")

payload := databases.Database{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, databases.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
