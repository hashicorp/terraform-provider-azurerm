
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/backuppolicies` Documentation

The `backuppolicies` SDK allows for interaction with Azure Resource Manager `recoveryservicesbackup` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/backuppolicies"
```


### Client Initialization

```go
client := backuppolicies.NewBackupPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupPoliciesClient.List`

```go
ctx := context.TODO()
id := backuppolicies.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName")

// alternatively `client.List(ctx, id, backuppolicies.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, backuppolicies.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
