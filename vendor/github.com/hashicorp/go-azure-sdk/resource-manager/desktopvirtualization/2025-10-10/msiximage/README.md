
## `github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/msiximage` Documentation

The `msiximage` SDK allows for interaction with Azure Resource Manager `desktopvirtualization` (API Version `2025-10-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/msiximage"
```


### Client Initialization

```go
client := msiximage.NewMsixImageClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MsixImageClient.Expand`

```go
ctx := context.TODO()
id := msiximage.NewHostPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hostPoolName")

payload := msiximage.MSIXImageURI{
	// ...
}


// alternatively `client.Expand(ctx, id, payload)` can be used to do batched pagination
items, err := client.ExpandComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
