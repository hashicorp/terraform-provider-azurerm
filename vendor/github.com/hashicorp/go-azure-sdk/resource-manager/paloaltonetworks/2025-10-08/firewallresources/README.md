
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources` Documentation

The `firewallresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
```


### Client Initialization

```go
client := firewallresources.NewFirewallResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallResourcesClient.FirewallsCreateOrUpdate`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := firewallresources.FirewallResource{
	// ...
}


if err := client.FirewallsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallResourcesClient.FirewallsDelete`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

if err := client.FirewallsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallResourcesClient.FirewallsGet`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.FirewallsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallResourcesClient.FirewallsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.FirewallsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.FirewallsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FirewallResourcesClient.FirewallsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.FirewallsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.FirewallsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FirewallResourcesClient.FirewallsUpdate`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := firewallresources.FirewallResourceUpdate{
	// ...
}


read, err := client.FirewallsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallResourcesClient.FirewallsgetGlobalRulestack`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.FirewallsgetGlobalRulestack(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallResourcesClient.FirewallsgetLogProfile`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.FirewallsgetLogProfile(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallResourcesClient.FirewallsgetSupportInfo`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.FirewallsgetSupportInfo(ctx, id, firewallresources.DefaultFirewallsgetSupportInfoOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallResourcesClient.FirewallssaveLogProfile`

```go
ctx := context.TODO()
id := firewallresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

payload := firewallresources.LogSettings{
	// ...
}


read, err := client.FirewallssaveLogProfile(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
