
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/attacheddatabaseconfigurations` Documentation

The `attacheddatabaseconfigurations` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2023-05-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/attacheddatabaseconfigurations"
```


### Client Initialization

```go
client := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AttachedDatabaseConfigurationsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := attacheddatabaseconfigurations.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

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
id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "attachedDatabaseConfigurationValue")

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
id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "attachedDatabaseConfigurationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AttachedDatabaseConfigurationsClient.Get`

```go
ctx := context.TODO()
id := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "attachedDatabaseConfigurationValue")

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
id := attacheddatabaseconfigurations.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

read, err := client.ListByCluster(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
