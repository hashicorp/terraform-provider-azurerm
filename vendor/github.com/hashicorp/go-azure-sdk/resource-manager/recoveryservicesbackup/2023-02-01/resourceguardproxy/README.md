
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy` Documentation

The `resourceguardproxy` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicesbackup` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
```


### Client Initialization

```go
client := resourceguardproxy.NewResourceGuardProxyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceGuardProxyClient.Delete`

```go
ctx := context.TODO()
id := resourceguardproxy.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "backupResourceGuardProxyValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardProxyClient.Get`

```go
ctx := context.TODO()
id := resourceguardproxy.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "backupResourceGuardProxyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardProxyClient.Put`

```go
ctx := context.TODO()
id := resourceguardproxy.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "backupResourceGuardProxyValue")

payload := resourceguardproxy.ResourceGuardProxyBaseResource{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceGuardProxyClient.UnlockDelete`

```go
ctx := context.TODO()
id := resourceguardproxy.NewBackupResourceGuardProxyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue", "backupResourceGuardProxyValue")

payload := resourceguardproxy.UnlockDeleteRequest{
	// ...
}


read, err := client.UnlockDelete(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
