
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgatewaywafdynamicmanifests` Documentation

The `applicationgatewaywafdynamicmanifests` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgatewaywafdynamicmanifests"
```


### Client Initialization

```go
client := applicationgatewaywafdynamicmanifests.NewApplicationGatewayWafDynamicManifestsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationGatewayWafDynamicManifestsClient.ApplicationGatewayWafDynamicManifestsDefaultGet`

```go
ctx := context.TODO()
id := applicationgatewaywafdynamicmanifests.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.ApplicationGatewayWafDynamicManifestsDefaultGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewayWafDynamicManifestsClient.ApplicationGatewayWafDynamicManifestsGet`

```go
ctx := context.TODO()
id := applicationgatewaywafdynamicmanifests.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.ApplicationGatewayWafDynamicManifestsGet(ctx, id)` can be used to do batched pagination
items, err := client.ApplicationGatewayWafDynamicManifestsGetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
