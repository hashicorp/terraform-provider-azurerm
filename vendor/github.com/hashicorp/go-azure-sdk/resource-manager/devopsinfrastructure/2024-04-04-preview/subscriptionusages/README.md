
## `github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/subscriptionusages` Documentation

The `subscriptionusages` SDK allows for interaction with Azure Resource Manager `devopsinfrastructure` (API Version `2024-04-04-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/subscriptionusages"
```


### Client Initialization

```go
client := subscriptionusages.NewSubscriptionUsagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubscriptionUsagesClient.Usages`

```go
ctx := context.TODO()
id := subscriptionusages.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.Usages(ctx, id)` can be used to do batched pagination
items, err := client.UsagesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
