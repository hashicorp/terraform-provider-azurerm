
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverresetgtid` Documentation

The `serverresetgtid` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverresetgtid"
```


### Client Initialization

```go
client := serverresetgtid.NewServerResetGtidClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerResetGtidClient.ServersResetGtid`

```go
ctx := context.TODO()
id := serverresetgtid.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := serverresetgtid.ServerGtidSetParameter{
	// ...
}


if err := client.ServersResetGtidThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
