
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelaglinks` Documentation

The `expressroutelaglinks` SDK allows for interaction with Azure Resource Manager `network` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutelaglinks"
```


### Client Initialization

```go
client := expressroutelaglinks.NewExpressRouteLagLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteLagLinksClient.ExpressRouteLagsLinksGet`

```go
ctx := context.TODO()
id := expressroutelaglinks.NewLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName", "linkName")

read, err := client.ExpressRouteLagsLinksGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteLagLinksClient.ExpressRouteLagsLinksList`

```go
ctx := context.TODO()
id := expressroutelaglinks.NewExpressRouteLagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteLagName")

// alternatively `client.ExpressRouteLagsLinksList(ctx, id)` can be used to do batched pagination
items, err := client.ExpressRouteLagsLinksListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
