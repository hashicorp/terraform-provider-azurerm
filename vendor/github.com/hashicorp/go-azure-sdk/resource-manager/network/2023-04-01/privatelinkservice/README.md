
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privatelinkservice` Documentation

The `privatelinkservice` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privatelinkservice"
```


### Client Initialization

```go
client := privatelinkservice.NewPrivateLinkServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkServiceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatelinkservice.NewPrivateLinkServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkServiceValue")

payload := privatelinkservice.PrivateLinkService{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
