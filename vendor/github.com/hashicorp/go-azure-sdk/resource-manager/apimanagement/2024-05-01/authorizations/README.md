
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizations` Documentation

The `authorizations` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizations"
```


### Client Initialization

```go
client := authorizations.NewAuthorizationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationsClient.AuthorizationListByAuthorizationProvider`

```go
ctx := context.TODO()
id := authorizations.NewAuthorizationProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId")

// alternatively `client.AuthorizationListByAuthorizationProvider(ctx, id, authorizations.DefaultAuthorizationListByAuthorizationProviderOperationOptions())` can be used to do batched pagination
items, err := client.AuthorizationListByAuthorizationProviderComplete(ctx, id, authorizations.DefaultAuthorizationListByAuthorizationProviderOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
