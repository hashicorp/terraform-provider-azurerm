
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/webapplicationfirewallpolicies` Documentation

The `webapplicationfirewallpolicies` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/webapplicationfirewallpolicies"
```


### Client Initialization

```go
client := webapplicationfirewallpolicies.NewWebApplicationFirewallPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebApplicationFirewallPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayWebApplicationFirewallPolicyValue")

payload := webapplicationfirewallpolicies.WebApplicationFirewallPolicy{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebApplicationFirewallPoliciesClient.Delete`

```go
ctx := context.TODO()
id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayWebApplicationFirewallPolicyValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebApplicationFirewallPoliciesClient.Get`

```go
ctx := context.TODO()
id := webapplicationfirewallpolicies.NewApplicationGatewayWebApplicationFirewallPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayWebApplicationFirewallPolicyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebApplicationFirewallPoliciesClient.List`

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


### Example Usage: `WebApplicationFirewallPoliciesClient.ListAll`

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
