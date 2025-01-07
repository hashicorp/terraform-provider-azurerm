
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverrestart` Documentation

The `serverrestart` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverrestart"
```


### Client Initialization

```go
client := serverrestart.NewServerRestartClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerRestartClient.ServersRestart`

```go
ctx := context.TODO()
id := serverrestart.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := serverrestart.ServerRestartParameter{
	// ...
}


if err := client.ServersRestartThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
