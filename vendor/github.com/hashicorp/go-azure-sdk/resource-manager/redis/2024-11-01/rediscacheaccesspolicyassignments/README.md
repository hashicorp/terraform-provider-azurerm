
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicyassignments` Documentation

The `rediscacheaccesspolicyassignments` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicyassignments"
```


### Client Initialization

```go
client := rediscacheaccesspolicyassignments.NewRedisCacheAccessPolicyAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisCacheAccessPolicyAssignmentsClient.AccessPolicyAssignmentCreateUpdate`

```go
ctx := context.TODO()
id := rediscacheaccesspolicyassignments.NewAccessPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "accessPolicyAssignmentName")

payload := rediscacheaccesspolicyassignments.RedisCacheAccessPolicyAssignment{
	// ...
}


if err := client.AccessPolicyAssignmentCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisCacheAccessPolicyAssignmentsClient.AccessPolicyAssignmentDelete`

```go
ctx := context.TODO()
id := rediscacheaccesspolicyassignments.NewAccessPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "accessPolicyAssignmentName")

if err := client.AccessPolicyAssignmentDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisCacheAccessPolicyAssignmentsClient.AccessPolicyAssignmentGet`

```go
ctx := context.TODO()
id := rediscacheaccesspolicyassignments.NewAccessPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "accessPolicyAssignmentName")

read, err := client.AccessPolicyAssignmentGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisCacheAccessPolicyAssignmentsClient.AccessPolicyAssignmentList`

```go
ctx := context.TODO()
id := rediscacheaccesspolicyassignments.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

// alternatively `client.AccessPolicyAssignmentList(ctx, id)` can be used to do batched pagination
items, err := client.AccessPolicyAssignmentListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
