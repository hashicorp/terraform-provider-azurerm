
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverupgrade` Documentation

The `serverupgrade` SDK allows for interaction with the Azure Resource Manager Service `mysql` (API Version `2020-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverupgrade"
```


### Client Initialization

```go
client := serverupgrade.NewServerUpgradeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerUpgradeClient.ServersUpgrade`

```go
ctx := context.TODO()
id := serverupgrade.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := serverupgrade.ServerUpgradeParameters{
	// ...
}


if err := client.ServersUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
