
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules` Documentation

The `volumequotarules` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
```


### Client Initialization

```go
client := volumequotarules.NewVolumeQuotaRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VolumeQuotaRulesClient.Create`

```go
ctx := context.TODO()
id := volumequotarules.NewVolumeQuotaRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "volumeQuotaRuleName")

payload := volumequotarules.VolumeQuotaRule{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumeQuotaRulesClient.Delete`

```go
ctx := context.TODO()
id := volumequotarules.NewVolumeQuotaRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "volumeQuotaRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumeQuotaRulesClient.Get`

```go
ctx := context.TODO()
id := volumequotarules.NewVolumeQuotaRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "volumeQuotaRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumeQuotaRulesClient.ListByVolume`

```go
ctx := context.TODO()
id := volumequotarules.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

read, err := client.ListByVolume(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumeQuotaRulesClient.Update`

```go
ctx := context.TODO()
id := volumequotarules.NewVolumeQuotaRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "volumeQuotaRuleName")

payload := volumequotarules.VolumeQuotaRulePatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
