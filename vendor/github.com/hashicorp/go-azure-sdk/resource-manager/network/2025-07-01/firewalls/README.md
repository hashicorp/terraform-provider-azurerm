
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewalls` Documentation

The `firewalls` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewalls"
```


### Client Initialization

```go
client := firewalls.NewFirewallsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallsClient.AzureFirewallsListLearnedPrefixes`

```go
ctx := context.TODO()
id := firewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

if err := client.AzureFirewallsListLearnedPrefixesThenPoll(ctx, id); err != nil {
	// handle the error
}
```
