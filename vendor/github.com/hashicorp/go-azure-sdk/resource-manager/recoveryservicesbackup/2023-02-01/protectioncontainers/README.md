
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers` Documentation

The `protectioncontainers` SDK allows for interaction with Azure Resource Manager `recoveryservicesbackup` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
```


### Client Initialization

```go
client := protectioncontainers.NewProtectionContainersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProtectionContainersClient.Get`

```go
ctx := context.TODO()
id := protectioncontainers.NewProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "backupFabricName", "protectionContainerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProtectionContainersClient.Inquire`

```go
ctx := context.TODO()
id := protectioncontainers.NewProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "backupFabricName", "protectionContainerName")

read, err := client.Inquire(ctx, id, protectioncontainers.DefaultInquireOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProtectionContainersClient.Refresh`

```go
ctx := context.TODO()
id := protectioncontainers.NewBackupFabricID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "backupFabricName")

read, err := client.Refresh(ctx, id, protectioncontainers.DefaultRefreshOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProtectionContainersClient.Register`

```go
ctx := context.TODO()
id := protectioncontainers.NewProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "backupFabricName", "protectionContainerName")

payload := protectioncontainers.ProtectionContainerResource{
	// ...
}


read, err := client.Register(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProtectionContainersClient.Unregister`

```go
ctx := context.TODO()
id := protectioncontainers.NewProtectionContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "backupFabricName", "protectionContainerName")

read, err := client.Unregister(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
