
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofileoperations` Documentation

The `autoupgradeprofileoperations` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofileoperations"
```


### Client Initialization

```go
client := autoupgradeprofileoperations.NewAutoUpgradeProfileOperationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutoUpgradeProfileOperationsClient.GenerateUpdateRun`

```go
ctx := context.TODO()
id := autoupgradeprofileoperations.NewAutoUpgradeProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "autoUpgradeProfileName")

if err := client.GenerateUpdateRunThenPoll(ctx, id); err != nil {
	// handle the error
}
```
