
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/authenticationpolicies` Documentation

The `authenticationpolicies` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/authenticationpolicies"
```


### Client Initialization

```go
client := authenticationpolicies.NewAuthenticationPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthenticationPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := authenticationpolicies.NewAuthenticationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "authenticationPolicyName")

payload := authenticationpolicies.AuthenticationPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AuthenticationPoliciesClient.Delete`

```go
ctx := context.TODO()
id := authenticationpolicies.NewAuthenticationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "authenticationPolicyName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthenticationPoliciesClient.Get`

```go
ctx := context.TODO()
id := authenticationpolicies.NewAuthenticationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "authenticationPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AuthenticationPoliciesClient.List`

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


### Example Usage: `AuthenticationPoliciesClient.ListAll`

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


### Example Usage: `AuthenticationPoliciesClient.Update`

```go
ctx := context.TODO()
id := authenticationpolicies.NewAuthenticationPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "authenticationPolicyName")

payload := authenticationpolicies.AuthenticationPolicyUpdateParameters{
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
