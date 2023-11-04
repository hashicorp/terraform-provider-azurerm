
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/longtermretentionbackups` Documentation

The `longtermretentionbackups` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-02-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/longtermretentionbackups"
```


### Client Initialization

```go
client := longtermretentionbackups.NewLongTermRetentionBackupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LongTermRetentionBackupsClient.Copy`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

payload := longtermretentionbackups.CopyLongTermRetentionBackupParameters{
	// ...
}


if err := client.CopyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LongTermRetentionBackupsClient.CopyByResourceGroup`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionDatabaseLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

payload := longtermretentionbackups.CopyLongTermRetentionBackupParameters{
	// ...
}


if err := client.CopyByResourceGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LongTermRetentionBackupsClient.Delete`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LongTermRetentionBackupsClient.DeleteByResourceGroup`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionDatabaseLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

if err := client.DeleteByResourceGroupThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LongTermRetentionBackupsClient.Get`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LongTermRetentionBackupsClient.GetByResourceGroup`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionDatabaseLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

read, err := client.GetByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LongTermRetentionBackupsClient.ListByDatabase`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionServerLongTermRetentionDatabaseID("12345678-1234-9876-4563-123456789012", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue")

// alternatively `client.ListByDatabase(ctx, id, longtermretentionbackups.DefaultListByDatabaseOperationOptions())` can be used to do batched pagination
items, err := client.ListByDatabaseComplete(ctx, id, longtermretentionbackups.DefaultListByDatabaseOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LongTermRetentionBackupsClient.ListByLocation`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.ListByLocation(ctx, id, longtermretentionbackups.DefaultListByLocationOperationOptions())` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id, longtermretentionbackups.DefaultListByLocationOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LongTermRetentionBackupsClient.ListByResourceGroupDatabase`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLocationLongTermRetentionServerLongTermRetentionDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue")

// alternatively `client.ListByResourceGroupDatabase(ctx, id, longtermretentionbackups.DefaultListByResourceGroupDatabaseOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupDatabaseComplete(ctx, id, longtermretentionbackups.DefaultListByResourceGroupDatabaseOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LongTermRetentionBackupsClient.ListByResourceGroupLocation`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue")

// alternatively `client.ListByResourceGroupLocation(ctx, id, longtermretentionbackups.DefaultListByResourceGroupLocationOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupLocationComplete(ctx, id, longtermretentionbackups.DefaultListByResourceGroupLocationOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LongTermRetentionBackupsClient.ListByResourceGroupServer`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLocationLongTermRetentionServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "longTermRetentionServerValue")

// alternatively `client.ListByResourceGroupServer(ctx, id, longtermretentionbackups.DefaultListByResourceGroupServerOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupServerComplete(ctx, id, longtermretentionbackups.DefaultListByResourceGroupServerOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LongTermRetentionBackupsClient.ListByServer`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionServerID("12345678-1234-9876-4563-123456789012", "locationValue", "longTermRetentionServerValue")

// alternatively `client.ListByServer(ctx, id, longtermretentionbackups.DefaultListByServerOperationOptions())` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id, longtermretentionbackups.DefaultListByServerOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LongTermRetentionBackupsClient.Update`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

payload := longtermretentionbackups.UpdateLongTermRetentionBackupParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LongTermRetentionBackupsClient.UpdateByResourceGroup`

```go
ctx := context.TODO()
id := longtermretentionbackups.NewLongTermRetentionDatabaseLongTermRetentionBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationValue", "longTermRetentionServerValue", "longTermRetentionDatabaseValue", "longTermRetentionBackupValue")

payload := longtermretentionbackups.UpdateLongTermRetentionBackupParameters{
	// ...
}


if err := client.UpdateByResourceGroupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
