
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/replicationlinks` Documentation

The `replicationlinks` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/replicationlinks"
```


### Client Initialization

```go
client := replicationlinks.NewReplicationLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicationLinksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := replicationlinks.NewReplicationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName", "linkId")

payload := replicationlinks.ReplicationLink{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationLinksClient.Delete`

```go
ctx := context.TODO()
id := replicationlinks.NewReplicationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName", "linkId")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationLinksClient.Failover`

```go
ctx := context.TODO()
id := replicationlinks.NewReplicationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName", "linkId")

if err := client.FailoverThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationLinksClient.FailoverAllowDataLoss`

```go
ctx := context.TODO()
id := replicationlinks.NewReplicationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName", "linkId")

if err := client.FailoverAllowDataLossThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ReplicationLinksClient.Get`

```go
ctx := context.TODO()
id := replicationlinks.NewReplicationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName", "linkId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReplicationLinksClient.ListByDatabase`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName")

// alternatively `client.ListByDatabase(ctx, id)` can be used to do batched pagination
items, err := client.ListByDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationLinksClient.ListByServer`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ReplicationLinksClient.Update`

```go
ctx := context.TODO()
id := replicationlinks.NewReplicationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName", "linkId")

payload := replicationlinks.ReplicationLinkUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
