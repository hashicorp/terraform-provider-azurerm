
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelagmembers` Documentation

The `expressroutelagmembers` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelagmembers"
```


### Client Initialization

```go
client := expressroutelagmembers.NewExpressRouteLagMembersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteLagMembersClient.ExpressRouteLagsMembersGet`

```go
ctx := context.TODO()
id := expressroutelagmembers.NewMemberID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "linkName", "memberName")

read, err := client.ExpressRouteLagsMembersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteLagMembersClient.ExpressRouteLagsMembersList`

```go
ctx := context.TODO()
id := expressroutelagmembers.NewLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "linkName")

// alternatively `client.ExpressRouteLagsMembersList(ctx, id)` can be used to do batched pagination
items, err := client.ExpressRouteLagsMembersListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
