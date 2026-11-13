
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/serversmigration` Documentation

The `serversmigration` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2024-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/serversmigration"
```


### Client Initialization

```go
client := serversmigration.NewServersMigrationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServersMigrationClient.CutoverMigration`

```go
ctx := context.TODO()
id := serversmigration.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

if err := client.CutoverMigrationThenPoll(ctx, id); err != nil {
	// handle the error
}
```
