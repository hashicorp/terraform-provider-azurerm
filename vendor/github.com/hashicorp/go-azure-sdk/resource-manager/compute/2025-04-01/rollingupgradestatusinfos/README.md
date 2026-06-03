
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/rollingupgradestatusinfos` Documentation

The `rollingupgradestatusinfos` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2025-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/rollingupgradestatusinfos"
```


### Client Initialization

```go
client := rollingupgradestatusinfos.NewRollingUpgradeStatusInfosClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RollingUpgradeStatusInfosClient.VirtualMachineScaleSetRollingUpgradesGetLatest`

```go
ctx := context.TODO()
id := rollingupgradestatusinfos.NewVirtualMachineScaleSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualMachineScaleSetName")

read, err := client.VirtualMachineScaleSetRollingUpgradesGetLatest(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
