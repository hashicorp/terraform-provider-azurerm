
## `github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets` Documentation

The `dnsforwardingrulesets` SDK allows for interaction with Azure Resource Manager `dnsresolver` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
```


### Client Initialization

```go
client := dnsforwardingrulesets.NewDnsForwardingRulesetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DnsForwardingRulesetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dnsforwardingrulesets.NewDnsForwardingRulesetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName")

payload := dnsforwardingrulesets.DnsForwardingRuleset{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, dnsforwardingrulesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DnsForwardingRulesetsClient.Delete`

```go
ctx := context.TODO()
id := dnsforwardingrulesets.NewDnsForwardingRulesetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName")

if err := client.DeleteThenPoll(ctx, id, dnsforwardingrulesets.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DnsForwardingRulesetsClient.Get`

```go
ctx := context.TODO()
id := dnsforwardingrulesets.NewDnsForwardingRulesetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DnsForwardingRulesetsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, dnsforwardingrulesets.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, dnsforwardingrulesets.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DnsForwardingRulesetsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, dnsforwardingrulesets.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, dnsforwardingrulesets.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DnsForwardingRulesetsClient.ListByVirtualNetwork`

```go
ctx := context.TODO()
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

// alternatively `client.ListByVirtualNetwork(ctx, id, dnsforwardingrulesets.DefaultListByVirtualNetworkOperationOptions())` can be used to do batched pagination
items, err := client.ListByVirtualNetworkComplete(ctx, id, dnsforwardingrulesets.DefaultListByVirtualNetworkOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DnsForwardingRulesetsClient.Update`

```go
ctx := context.TODO()
id := dnsforwardingrulesets.NewDnsForwardingRulesetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName")

payload := dnsforwardingrulesets.DnsForwardingRulesetPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, dnsforwardingrulesets.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
