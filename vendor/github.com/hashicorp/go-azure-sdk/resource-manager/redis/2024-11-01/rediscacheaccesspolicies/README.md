
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicies` Documentation

The `rediscacheaccesspolicies` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicies"
```


### Client Initialization

```go
client := rediscacheaccesspolicies.NewRedisCacheAccessPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisCacheAccessPoliciesClient.AccessPolicyCreateUpdate`

```go
ctx := context.TODO()
id := rediscacheaccesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "accessPolicyName")

payload := rediscacheaccesspolicies.RedisCacheAccessPolicy{
	// ...
}


if err := client.AccessPolicyCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisCacheAccessPoliciesClient.AccessPolicyDelete`

```go
ctx := context.TODO()
id := rediscacheaccesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "accessPolicyName")

if err := client.AccessPolicyDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisCacheAccessPoliciesClient.AccessPolicyGet`

```go
ctx := context.TODO()
id := rediscacheaccesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "accessPolicyName")

read, err := client.AccessPolicyGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisCacheAccessPoliciesClient.AccessPolicyList`

```go
ctx := context.TODO()
id := rediscacheaccesspolicies.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

// alternatively `client.AccessPolicyList(ctx, id)` can be used to do batched pagination
items, err := client.AccessPolicyListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
