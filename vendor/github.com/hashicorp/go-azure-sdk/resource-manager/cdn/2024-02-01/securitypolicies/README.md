
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies` Documentation

The `securitypolicies` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2024-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
```


### Client Initialization

```go
client := securitypolicies.NewSecurityPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityPoliciesClient.Create`

```go
ctx := context.TODO()
id := securitypolicies.NewSecurityPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "securityPolicyName")

payload := securitypolicies.SecurityPolicy{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityPoliciesClient.Delete`

```go
ctx := context.TODO()
id := securitypolicies.NewSecurityPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "securityPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityPoliciesClient.Get`

```go
ctx := context.TODO()
id := securitypolicies.NewSecurityPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "securityPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityPoliciesClient.ListByProfile`

```go
ctx := context.TODO()
id := securitypolicies.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecurityPoliciesClient.Patch`

```go
ctx := context.TODO()
id := securitypolicies.NewSecurityPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "securityPolicyName")

payload := securitypolicies.SecurityPolicyUpdateParameters{
	// ...
}


if err := client.PatchThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
