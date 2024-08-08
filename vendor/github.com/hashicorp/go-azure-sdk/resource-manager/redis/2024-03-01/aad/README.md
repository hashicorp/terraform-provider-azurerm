
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/aad` Documentation

The `aad` SDK allows for interaction with the Azure Resource Manager Service `redis` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/aad"
```


### Client Initialization

```go
client := aad.NewAADClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AADClient.AccessPolicyAssignmentCreateUpdate`

```go
ctx := context.TODO()
id := aad.NewAccessPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "accessPolicyAssignmentValue")

payload := aad.RedisCacheAccessPolicyAssignment{
	// ...
}


if err := client.AccessPolicyAssignmentCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AADClient.AccessPolicyAssignmentDelete`

```go
ctx := context.TODO()
id := aad.NewAccessPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "accessPolicyAssignmentValue")

if err := client.AccessPolicyAssignmentDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AADClient.AccessPolicyAssignmentGet`

```go
ctx := context.TODO()
id := aad.NewAccessPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "accessPolicyAssignmentValue")

read, err := client.AccessPolicyAssignmentGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AADClient.AccessPolicyAssignmentList`

```go
ctx := context.TODO()
id := aad.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

// alternatively `client.AccessPolicyAssignmentList(ctx, id)` can be used to do batched pagination
items, err := client.AccessPolicyAssignmentListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AADClient.AccessPolicyCreateUpdate`

```go
ctx := context.TODO()
id := aad.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "accessPolicyValue")

payload := aad.RedisCacheAccessPolicy{
	// ...
}


if err := client.AccessPolicyCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AADClient.AccessPolicyDelete`

```go
ctx := context.TODO()
id := aad.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "accessPolicyValue")

if err := client.AccessPolicyDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AADClient.AccessPolicyGet`

```go
ctx := context.TODO()
id := aad.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "accessPolicyValue")

read, err := client.AccessPolicyGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AADClient.AccessPolicyList`

```go
ctx := context.TODO()
id := aad.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

// alternatively `client.AccessPolicyList(ctx, id)` can be used to do batched pagination
items, err := client.AccessPolicyListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
