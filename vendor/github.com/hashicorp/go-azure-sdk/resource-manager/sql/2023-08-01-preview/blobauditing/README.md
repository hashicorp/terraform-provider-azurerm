
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/blobauditing` Documentation

The `blobauditing` SDK allows for interaction with the Azure Resource Manager Service `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/blobauditing"
```


### Client Initialization

```go
client := blobauditing.NewBlobAuditingClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BlobAuditingClient.DatabaseBlobAuditingPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

payload := blobauditing.DatabaseBlobAuditingPolicy{
	// ...
}


read, err := client.DatabaseBlobAuditingPoliciesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobAuditingClient.DatabaseBlobAuditingPoliciesGet`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

read, err := client.DatabaseBlobAuditingPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobAuditingClient.DatabaseBlobAuditingPoliciesListByDatabase`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

// alternatively `client.DatabaseBlobAuditingPoliciesListByDatabase(ctx, id)` can be used to do batched pagination
items, err := client.DatabaseBlobAuditingPoliciesListByDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BlobAuditingClient.ExtendedDatabaseBlobAuditingPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

payload := blobauditing.ExtendedDatabaseBlobAuditingPolicy{
	// ...
}


read, err := client.ExtendedDatabaseBlobAuditingPoliciesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobAuditingClient.ExtendedDatabaseBlobAuditingPoliciesGet`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

read, err := client.ExtendedDatabaseBlobAuditingPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobAuditingClient.ExtendedDatabaseBlobAuditingPoliciesListByDatabase`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "databaseValue")

// alternatively `client.ExtendedDatabaseBlobAuditingPoliciesListByDatabase(ctx, id)` can be used to do batched pagination
items, err := client.ExtendedDatabaseBlobAuditingPoliciesListByDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BlobAuditingClient.ExtendedServerBlobAuditingPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := blobauditing.ExtendedServerBlobAuditingPolicy{
	// ...
}


if err := client.ExtendedServerBlobAuditingPoliciesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BlobAuditingClient.ExtendedServerBlobAuditingPoliciesGet`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.ExtendedServerBlobAuditingPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobAuditingClient.ExtendedServerBlobAuditingPoliciesListByServer`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

// alternatively `client.ExtendedServerBlobAuditingPoliciesListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ExtendedServerBlobAuditingPoliciesListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BlobAuditingClient.ServerBlobAuditingPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := blobauditing.ServerBlobAuditingPolicy{
	// ...
}


if err := client.ServerBlobAuditingPoliciesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BlobAuditingClient.ServerBlobAuditingPoliciesGet`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.ServerBlobAuditingPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobAuditingClient.ServerBlobAuditingPoliciesListByServer`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

// alternatively `client.ServerBlobAuditingPoliciesListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ServerBlobAuditingPoliciesListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
