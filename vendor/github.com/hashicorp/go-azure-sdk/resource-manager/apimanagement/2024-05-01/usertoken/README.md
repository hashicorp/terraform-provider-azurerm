
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/usertoken` Documentation

The `usertoken` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/usertoken"
```


### Client Initialization

```go
client := usertoken.NewUserTokenClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `UserTokenClient.UserGetSharedAccessToken`

```go
ctx := context.TODO()
id := usertoken.NewUserID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "userId")

payload := usertoken.UserTokenParameters{
	// ...
}


read, err := client.UserGetSharedAccessToken(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
