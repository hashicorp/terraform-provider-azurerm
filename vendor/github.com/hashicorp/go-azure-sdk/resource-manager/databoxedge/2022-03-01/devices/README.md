
## `github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices` Documentation

The `devices` SDK allows for interaction with the Azure Resource Manager Service `databoxedge` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
```


### Client Initialization

```go
client := devices.NewDevicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DevicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

payload := devices.DataBoxEdgeDevice{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.CreateOrUpdateSecuritySettings`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

payload := devices.SecuritySettings{
	// ...
}


if err := client.CreateOrUpdateSecuritySettingsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DevicesClient.Delete`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DevicesClient.DownloadUpdates`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

if err := client.DownloadUpdatesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DevicesClient.GenerateCertificate`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

read, err := client.GenerateCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.Get`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.GetExtendedInformation`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

read, err := client.GetExtendedInformation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.GetNetworkSettings`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

read, err := client.GetNetworkSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.GetUpdateSummary`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

read, err := client.GetUpdateSummary(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.InstallUpdates`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

if err := client.InstallUpdatesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DevicesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := devices.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, devices.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, devices.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DevicesClient.ListBySubscription`

```go
ctx := context.TODO()
id := devices.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, devices.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, devices.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DevicesClient.ScanForUpdates`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

if err := client.ScanForUpdatesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DevicesClient.Update`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

payload := devices.DataBoxEdgeDevicePatch{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.UpdateExtendedInformation`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

payload := devices.DataBoxEdgeDeviceExtendedInfoPatch{
	// ...
}


read, err := client.UpdateExtendedInformation(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DevicesClient.UploadCertificate`

```go
ctx := context.TODO()
id := devices.NewDataBoxEdgeDeviceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataBoxEdgeDeviceValue")

payload := devices.UploadCertificateRequest{
	// ...
}


read, err := client.UploadCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
