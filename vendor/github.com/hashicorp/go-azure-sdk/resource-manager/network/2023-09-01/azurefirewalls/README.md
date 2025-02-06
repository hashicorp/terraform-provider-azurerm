
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/azurefirewalls` Documentation

The `azurefirewalls` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/azurefirewalls"
```


### Client Initialization

```go
client := azurefirewalls.NewAzureFirewallsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AzureFirewallsClient.AzureFirewallsListLearnedPrefixes`

```go
ctx := context.TODO()
id := azurefirewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

if err := client.AzureFirewallsListLearnedPrefixesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AzureFirewallsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := azurefirewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

payload := azurefirewalls.AzureFirewall{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AzureFirewallsClient.Delete`

```go
ctx := context.TODO()
id := azurefirewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AzureFirewallsClient.Get`

```go
ctx := context.TODO()
id := azurefirewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AzureFirewallsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AzureFirewallsClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AzureFirewallsClient.PacketCapture`

```go
ctx := context.TODO()
id := azurefirewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

payload := azurefirewalls.FirewallPacketCaptureParameters{
	// ...
}


if err := client.PacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AzureFirewallsClient.UpdateTags`

```go
ctx := context.TODO()
id := azurefirewalls.NewAzureFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureFirewallName")

payload := azurefirewalls.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
