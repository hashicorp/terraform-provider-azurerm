
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/firewallpolicies` Documentation

The `firewallpolicies` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/firewallpolicies"
```


### Client Initialization

```go
client := firewallpolicies.NewFirewallPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

payload := firewallpolicies.FirewallPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallPoliciesClient.Delete`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallPoliciesClient.FirewallPolicyIdpsSignaturesFilterValuesList`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

payload := firewallpolicies.SignatureOverridesFilterValuesQuery{
	// ...
}


read, err := client.FirewallPolicyIdpsSignaturesFilterValuesList(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.FirewallPolicyIdpsSignaturesList`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

payload := firewallpolicies.IDPSQueryObject{
	// ...
}


read, err := client.FirewallPolicyIdpsSignaturesList(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.FirewallPolicyIdpsSignaturesOverridesGet`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

read, err := client.FirewallPolicyIdpsSignaturesOverridesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.FirewallPolicyIdpsSignaturesOverridesList`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

read, err := client.FirewallPolicyIdpsSignaturesOverridesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.FirewallPolicyIdpsSignaturesOverridesPatch`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

payload := firewallpolicies.SignaturesOverrides{
	// ...
}


read, err := client.FirewallPolicyIdpsSignaturesOverridesPatch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.FirewallPolicyIdpsSignaturesOverridesPut`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

payload := firewallpolicies.SignaturesOverrides{
	// ...
}


read, err := client.FirewallPolicyIdpsSignaturesOverridesPut(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.Get`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

read, err := client.Get(ctx, id, firewallpolicies.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallPoliciesClient.List`

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


### Example Usage: `FirewallPoliciesClient.ListAll`

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


### Example Usage: `FirewallPoliciesClient.UpdateTags`

```go
ctx := context.TODO()
id := firewallpolicies.NewFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallPolicyName")

payload := firewallpolicies.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
