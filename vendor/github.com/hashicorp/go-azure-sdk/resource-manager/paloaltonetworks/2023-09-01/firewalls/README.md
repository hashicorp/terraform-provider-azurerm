
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls` Documentation

The `firewalls` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
```


### Client Initialization

```go
client := firewalls.NewFirewallsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallsClient.AveLogProfile`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := firewalls.LogSettings{
	// ...
}


read, err := client.AveLogProfile(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := firewalls.FirewallResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallsClient.Delete`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallsClient.Get`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallsClient.GetGlobalRulestack`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.GetGlobalRulestack(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallsClient.GetLogProfile`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.GetLogProfile(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallsClient.GetSupportInfo`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.GetSupportInfo(ctx, id, firewalls.DefaultGetSupportInfoOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FirewallsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FirewallsClient.Update`

```go
ctx := context.TODO()
id := firewalls.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := firewalls.FirewallResourceUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
