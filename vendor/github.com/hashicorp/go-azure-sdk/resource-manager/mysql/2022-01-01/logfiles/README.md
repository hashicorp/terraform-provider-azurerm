
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/logfiles` Documentation

The `logfiles` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/logfiles"
```


### Client Initialization

```go
client := logfiles.NewLogFilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LogFilesClient.ListByServer`

```go
ctx := context.TODO()
id := logfiles.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
