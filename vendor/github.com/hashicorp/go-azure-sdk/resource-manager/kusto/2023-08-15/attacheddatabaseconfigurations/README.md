
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/attacheddatabaseconfigurations` Documentation

The `attacheddatabaseconfigurations` SDK allows for interaction with Azure Resource Manager `kusto` (API Version `2023-08-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/attacheddatabaseconfigurations"
```


### Client Initialization

```go
client := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AttachedDatabaseConfigurationsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

payload := attacheddatabaseconfigurations.AttachedDatabaseConfigurationsCheckNameRequest{
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


### Example Usage: `AttachedDatabaseConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "attachedDatabaseConfigurationName")

payload := attacheddatabaseconfigurations.AttachedDatabaseConfiguration{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedDatabaseConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "attachedDatabaseConfigurationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedDatabaseConfigurationsClient.Get`

```go
ctx := context.TODO()
id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "attachedDatabaseConfigurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttachedDatabaseConfigurationsClient.ListByCluster`

```go
ctx := context.TODO()
id := commonids.NewKustoClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

read, err := client.ListByCluster(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
