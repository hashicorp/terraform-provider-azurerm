
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backendreconnect` Documentation

The `backendreconnect` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backendreconnect"
```


### Client Initialization

```go
client := backendreconnect.NewBackendReconnectClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackendReconnectClient.BackendReconnect`

```go
ctx := context.TODO()
id := backendreconnect.NewBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "backendId")

payload := backendreconnect.BackendReconnectContract{
	// ...
}


read, err := client.BackendReconnect(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
