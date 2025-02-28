
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationloginlinks` Documentation

The `authorizationloginlinks` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/authorizationloginlinks"
```


### Client Initialization

```go
client := authorizationloginlinks.NewAuthorizationLoginLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AuthorizationLoginLinksClient.Post`

```go
ctx := context.TODO()
id := authorizationloginlinks.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "authorizationProviderId", "authorizationId")

payload := authorizationloginlinks.AuthorizationLoginRequestContract{
	// ...
}


read, err := client.Post(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
