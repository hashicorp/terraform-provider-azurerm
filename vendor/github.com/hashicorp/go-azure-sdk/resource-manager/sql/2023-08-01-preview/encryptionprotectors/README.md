
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/encryptionprotectors` Documentation

The `encryptionprotectors` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/encryptionprotectors"
```


### Client Initialization

```go
client := encryptionprotectors.NewEncryptionProtectorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EncryptionProtectorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

payload := encryptionprotectors.EncryptionProtector{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EncryptionProtectorsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncryptionProtectorsClient.ListByServer`

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


### Example Usage: `EncryptionProtectorsClient.Revalidate`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

if err := client.RevalidateThenPoll(ctx, id); err != nil {
	// handle the error
}
```
