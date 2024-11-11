
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps` Documentation

The `webapps` SDK allows for interaction with Azure Resource Manager `web` (API Version `2023-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
```


### Client Initialization

```go
client := webapps.NewWebAppsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WebAppsClient.AddPremierAddOn`

```go
ctx := context.TODO()
id := webapps.NewPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "premierAddonName")

payload := webapps.PremierAddOn{
	// ...
}


read, err := client.AddPremierAddOn(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.AddPremierAddOnSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "premierAddonName")

payload := webapps.PremierAddOn{
	// ...
}


read, err := client.AddPremierAddOnSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.AnalyzeCustomHostname`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.AnalyzeCustomHostname(ctx, id, webapps.DefaultAnalyzeCustomHostnameOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.AnalyzeCustomHostnameSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.AnalyzeCustomHostnameSlot(ctx, id, webapps.DefaultAnalyzeCustomHostnameSlotOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ApplySlotConfigToProduction`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.CsmSlotEntity{
	// ...
}


read, err := client.ApplySlotConfigToProduction(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ApplySlotConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.CsmSlotEntity{
	// ...
}


read, err := client.ApplySlotConfigurationSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ApproveOrRejectPrivateEndpointConnection`

```go
ctx := context.TODO()
id := webapps.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "privateEndpointConnectionName")

payload := webapps.RemotePrivateEndpointConnectionARMResource{
	// ...
}


if err := client.ApproveOrRejectPrivateEndpointConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.ApproveOrRejectPrivateEndpointConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "privateEndpointConnectionName")

payload := webapps.RemotePrivateEndpointConnectionARMResource{
	// ...
}


if err := client.ApproveOrRejectPrivateEndpointConnectionSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.Backup`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.BackupRequest{
	// ...
}


read, err := client.Backup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.BackupSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.BackupRequest{
	// ...
}


read, err := client.BackupSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateDeployment`

```go
ctx := context.TODO()
id := webapps.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "deploymentName")

payload := webapps.Deployment{
	// ...
}


read, err := client.CreateDeployment(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateDeploymentSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "deploymentName")

payload := webapps.Deployment{
	// ...
}


read, err := client.CreateDeploymentSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateFunction`

```go
ctx := context.TODO()
id := webapps.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName")

payload := webapps.FunctionEnvelope{
	// ...
}


if err := client.CreateFunctionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateInstanceFunctionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName")

payload := webapps.FunctionEnvelope{
	// ...
}


if err := client.CreateInstanceFunctionSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateInstanceMSDeployOperation`

```go
ctx := context.TODO()
id := webapps.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId")

payload := webapps.MSDeploy{
	// ...
}


if err := client.CreateInstanceMSDeployOperationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateInstanceMSDeployOperationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId")

payload := webapps.MSDeploy{
	// ...
}


if err := client.CreateInstanceMSDeployOperationSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateMSDeployOperation`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.MSDeploy{
	// ...
}


if err := client.CreateMSDeployOperationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateMSDeployOperationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.MSDeploy{
	// ...
}


if err := client.CreateMSDeployOperationSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateOneDeployOperation`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.CreateOneDeployOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.Site{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteConfigResource{
	// ...
}


read, err := client.CreateOrUpdateConfiguration(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteConfigResource{
	// ...
}


read, err := client.CreateOrUpdateConfigurationSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateDomainOwnershipIdentifier`

```go
ctx := context.TODO()
id := webapps.NewSiteDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "domainOwnershipIdentifierName")

payload := webapps.Identifier{
	// ...
}


read, err := client.CreateOrUpdateDomainOwnershipIdentifier(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateDomainOwnershipIdentifierSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "domainOwnershipIdentifierName")

payload := webapps.Identifier{
	// ...
}


read, err := client.CreateOrUpdateDomainOwnershipIdentifierSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateFunctionSecret`

```go
ctx := context.TODO()
id := webapps.NewKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName", "keyName")

payload := webapps.KeyInfo{
	// ...
}


read, err := client.CreateOrUpdateFunctionSecret(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateFunctionSecretSlot`

```go
ctx := context.TODO()
id := webapps.NewFunctionKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName", "keyName")

payload := webapps.KeyInfo{
	// ...
}


read, err := client.CreateOrUpdateFunctionSecretSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateHostNameBinding`

```go
ctx := context.TODO()
id := webapps.NewHostNameBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hostNameBindingName")

payload := webapps.HostNameBinding{
	// ...
}


read, err := client.CreateOrUpdateHostNameBinding(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateHostNameBindingSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHostNameBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hostNameBindingName")

payload := webapps.HostNameBinding{
	// ...
}


read, err := client.CreateOrUpdateHostNameBindingSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateHostSecret`

```go
ctx := context.TODO()
id := webapps.NewDefaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "defaultName", "keyName")

payload := webapps.KeyInfo{
	// ...
}


read, err := client.CreateOrUpdateHostSecret(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateHostSecretSlot`

```go
ctx := context.TODO()
id := webapps.NewHostDefaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "defaultName", "keyName")

payload := webapps.KeyInfo{
	// ...
}


read, err := client.CreateOrUpdateHostSecretSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateHybridConnection`

```go
ctx := context.TODO()
id := webapps.NewRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionNamespaceName", "relayName")

payload := webapps.HybridConnection{
	// ...
}


read, err := client.CreateOrUpdateHybridConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateHybridConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionNamespaceRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionNamespaceName", "relayName")

payload := webapps.HybridConnection{
	// ...
}


read, err := client.CreateOrUpdateHybridConnectionSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdatePublicCertificate`

```go
ctx := context.TODO()
id := webapps.NewPublicCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "publicCertificateName")

payload := webapps.PublicCertificate{
	// ...
}


read, err := client.CreateOrUpdatePublicCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdatePublicCertificateSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPublicCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "publicCertificateName")

payload := webapps.PublicCertificate{
	// ...
}


read, err := client.CreateOrUpdatePublicCertificateSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateRelayServiceConnection`

```go
ctx := context.TODO()
id := webapps.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionName")

payload := webapps.RelayServiceConnectionEntity{
	// ...
}


read, err := client.CreateOrUpdateRelayServiceConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateRelayServiceConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionName")

payload := webapps.RelayServiceConnectionEntity{
	// ...
}


read, err := client.CreateOrUpdateRelayServiceConnectionSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSiteContainer`

```go
ctx := context.TODO()
id := webapps.NewSitecontainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "sitecontainerName")

payload := webapps.SiteContainer{
	// ...
}


read, err := client.CreateOrUpdateSiteContainer(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSiteContainerSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotSitecontainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "sitecontainerName")

payload := webapps.SiteContainer{
	// ...
}


read, err := client.CreateOrUpdateSiteContainerSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.Site{
	// ...
}


if err := client.CreateOrUpdateSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSourceControl`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteSourceControl{
	// ...
}


if err := client.CreateOrUpdateSourceControlThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSourceControlSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteSourceControl{
	// ...
}


if err := client.CreateOrUpdateSourceControlSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SwiftVirtualNetwork{
	// ...
}


read, err := client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheck(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheckSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SwiftVirtualNetwork{
	// ...
}


read, err := client.CreateOrUpdateSwiftVirtualNetworkConnectionWithCheckSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateVnetConnection`

```go
ctx := context.TODO()
id := webapps.NewVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName")

payload := webapps.VnetInfoResource{
	// ...
}


read, err := client.CreateOrUpdateVnetConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateVnetConnectionGateway`

```go
ctx := context.TODO()
id := webapps.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName", "gatewayName")

payload := webapps.VnetGateway{
	// ...
}


read, err := client.CreateOrUpdateVnetConnectionGateway(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateVnetConnectionGatewaySlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName", "gatewayName")

payload := webapps.VnetGateway{
	// ...
}


read, err := client.CreateOrUpdateVnetConnectionGatewaySlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.CreateOrUpdateVnetConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName")

payload := webapps.VnetInfoResource{
	// ...
}


read, err := client.CreateOrUpdateVnetConnectionSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.Delete(ctx, id, webapps.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteBackup`

```go
ctx := context.TODO()
id := webapps.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "backupId")

read, err := client.DeleteBackup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteBackupConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.DeleteBackupConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteBackupConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.DeleteBackupConfigurationSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteBackupSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "backupId")

read, err := client.DeleteBackupSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteContinuousWebJob`

```go
ctx := context.TODO()
id := webapps.NewContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "continuousWebJobName")

read, err := client.DeleteContinuousWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteContinuousWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "continuousWebJobName")

read, err := client.DeleteContinuousWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteDeployment`

```go
ctx := context.TODO()
id := webapps.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "deploymentName")

read, err := client.DeleteDeployment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteDeploymentSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "deploymentName")

read, err := client.DeleteDeploymentSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteDomainOwnershipIdentifier`

```go
ctx := context.TODO()
id := webapps.NewSiteDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "domainOwnershipIdentifierName")

read, err := client.DeleteDomainOwnershipIdentifier(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteDomainOwnershipIdentifierSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "domainOwnershipIdentifierName")

read, err := client.DeleteDomainOwnershipIdentifierSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteFunction`

```go
ctx := context.TODO()
id := webapps.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName")

read, err := client.DeleteFunction(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteFunctionSecret`

```go
ctx := context.TODO()
id := webapps.NewKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName", "keyName")

read, err := client.DeleteFunctionSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteFunctionSecretSlot`

```go
ctx := context.TODO()
id := webapps.NewFunctionKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName", "keyName")

read, err := client.DeleteFunctionSecretSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteHostNameBinding`

```go
ctx := context.TODO()
id := webapps.NewHostNameBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hostNameBindingName")

read, err := client.DeleteHostNameBinding(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteHostNameBindingSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHostNameBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hostNameBindingName")

read, err := client.DeleteHostNameBindingSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteHostSecret`

```go
ctx := context.TODO()
id := webapps.NewDefaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "defaultName", "keyName")

read, err := client.DeleteHostSecret(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteHostSecretSlot`

```go
ctx := context.TODO()
id := webapps.NewHostDefaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "defaultName", "keyName")

read, err := client.DeleteHostSecretSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteHybridConnection`

```go
ctx := context.TODO()
id := webapps.NewRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionNamespaceName", "relayName")

read, err := client.DeleteHybridConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteHybridConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionNamespaceRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionNamespaceName", "relayName")

read, err := client.DeleteHybridConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteInstanceFunctionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName")

read, err := client.DeleteInstanceFunctionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteInstanceProcess`

```go
ctx := context.TODO()
id := webapps.NewInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId", "processId")

read, err := client.DeleteInstanceProcess(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteInstanceProcessSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId", "processId")

read, err := client.DeleteInstanceProcessSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeletePremierAddOn`

```go
ctx := context.TODO()
id := webapps.NewPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "premierAddonName")

read, err := client.DeletePremierAddOn(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeletePremierAddOnSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "premierAddonName")

read, err := client.DeletePremierAddOnSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeletePrivateEndpointConnection`

```go
ctx := context.TODO()
id := webapps.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "privateEndpointConnectionName")

if err := client.DeletePrivateEndpointConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.DeletePrivateEndpointConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "privateEndpointConnectionName")

if err := client.DeletePrivateEndpointConnectionSlotThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.DeleteProcess`

```go
ctx := context.TODO()
id := webapps.NewProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "processId")

read, err := client.DeleteProcess(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteProcessSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "processId")

read, err := client.DeleteProcessSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeletePublicCertificate`

```go
ctx := context.TODO()
id := webapps.NewPublicCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "publicCertificateName")

read, err := client.DeletePublicCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeletePublicCertificateSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPublicCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "publicCertificateName")

read, err := client.DeletePublicCertificateSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteRelayServiceConnection`

```go
ctx := context.TODO()
id := webapps.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionName")

read, err := client.DeleteRelayServiceConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteRelayServiceConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionName")

read, err := client.DeleteRelayServiceConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSiteContainer`

```go
ctx := context.TODO()
id := webapps.NewSitecontainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "sitecontainerName")

read, err := client.DeleteSiteContainer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSiteContainerSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotSitecontainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "sitecontainerName")

read, err := client.DeleteSiteContainerSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSiteExtension`

```go
ctx := context.TODO()
id := webapps.NewSiteExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "siteExtensionId")

read, err := client.DeleteSiteExtension(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSiteExtensionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotSiteExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "siteExtensionId")

read, err := client.DeleteSiteExtensionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.DeleteSlot(ctx, id, webapps.DefaultDeleteSlotOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSourceControl`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.DeleteSourceControl(ctx, id, webapps.DefaultDeleteSourceControlOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSourceControlSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.DeleteSourceControlSlot(ctx, id, webapps.DefaultDeleteSourceControlSlotOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSwiftVirtualNetwork`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.DeleteSwiftVirtualNetwork(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteSwiftVirtualNetworkSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.DeleteSwiftVirtualNetworkSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteTriggeredWebJob`

```go
ctx := context.TODO()
id := webapps.NewTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "triggeredWebJobName")

read, err := client.DeleteTriggeredWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteTriggeredWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "triggeredWebJobName")

read, err := client.DeleteTriggeredWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteVnetConnection`

```go
ctx := context.TODO()
id := webapps.NewVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName")

read, err := client.DeleteVnetConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeleteVnetConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName")

read, err := client.DeleteVnetConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeployWorkflowArtifacts`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.WorkflowArtifacts{
	// ...
}


read, err := client.DeployWorkflowArtifacts(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DeployWorkflowArtifactsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.WorkflowArtifacts{
	// ...
}


read, err := client.DeployWorkflowArtifactsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DiscoverBackup`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.RestoreRequest{
	// ...
}


read, err := client.DiscoverBackup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.DiscoverBackupSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.RestoreRequest{
	// ...
}


read, err := client.DiscoverBackupSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GenerateNewSitePublishingPassword`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GenerateNewSitePublishingPassword(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GenerateNewSitePublishingPasswordSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GenerateNewSitePublishingPasswordSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAppSettingKeyVaultReference`

```go
ctx := context.TODO()
id := webapps.NewAppSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "appSettingKey")

read, err := client.GetAppSettingKeyVaultReference(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAppSettingKeyVaultReferenceSlot`

```go
ctx := context.TODO()
id := webapps.NewConfigReferenceAppSettingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "appSettingKey")

read, err := client.GetAppSettingKeyVaultReferenceSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAppSettingsKeyVaultReferences`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.GetAppSettingsKeyVaultReferences(ctx, id)` can be used to do batched pagination
items, err := client.GetAppSettingsKeyVaultReferencesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.GetAppSettingsKeyVaultReferencesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.GetAppSettingsKeyVaultReferencesSlot(ctx, id)` can be used to do batched pagination
items, err := client.GetAppSettingsKeyVaultReferencesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.GetAuthSettings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetAuthSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAuthSettingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetAuthSettingsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAuthSettingsV2`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetAuthSettingsV2(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAuthSettingsV2Slot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetAuthSettingsV2Slot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAuthSettingsV2WithoutSecrets`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetAuthSettingsV2WithoutSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetAuthSettingsV2WithoutSecretsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetAuthSettingsV2WithoutSecretsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetBackupConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetBackupConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetBackupConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetBackupConfigurationSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetBackupStatus`

```go
ctx := context.TODO()
id := webapps.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "backupId")

read, err := client.GetBackupStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetBackupStatusSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "backupId")

read, err := client.GetBackupStatusSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetConfigurationSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetConfigurationSnapshot`

```go
ctx := context.TODO()
id := webapps.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "snapshotId")

read, err := client.GetConfigurationSnapshot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetConfigurationSnapshotSlot`

```go
ctx := context.TODO()
id := webapps.NewWebSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "snapshotId")

read, err := client.GetConfigurationSnapshotSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetContainerLogsZip`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetContainerLogsZip(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetContainerLogsZipSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetContainerLogsZipSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetContinuousWebJob`

```go
ctx := context.TODO()
id := webapps.NewContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "continuousWebJobName")

read, err := client.GetContinuousWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetContinuousWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "continuousWebJobName")

read, err := client.GetContinuousWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetDeployment`

```go
ctx := context.TODO()
id := webapps.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "deploymentName")

read, err := client.GetDeployment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetDeploymentSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "deploymentName")

read, err := client.GetDeploymentSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetDiagnosticLogsConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetDiagnosticLogsConfiguration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetDiagnosticLogsConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetDiagnosticLogsConfigurationSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetDomainOwnershipIdentifier`

```go
ctx := context.TODO()
id := webapps.NewSiteDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "domainOwnershipIdentifierName")

read, err := client.GetDomainOwnershipIdentifier(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetDomainOwnershipIdentifierSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "domainOwnershipIdentifierName")

read, err := client.GetDomainOwnershipIdentifierSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetFtpAllowed`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetFtpAllowed(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetFtpAllowedSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetFtpAllowedSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetFunction`

```go
ctx := context.TODO()
id := webapps.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName")

read, err := client.GetFunction(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetFunctionsAdminToken`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetFunctionsAdminToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetFunctionsAdminTokenSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetFunctionsAdminTokenSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetHostNameBinding`

```go
ctx := context.TODO()
id := webapps.NewHostNameBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hostNameBindingName")

read, err := client.GetHostNameBinding(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetHostNameBindingSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHostNameBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hostNameBindingName")

read, err := client.GetHostNameBindingSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetHybridConnection`

```go
ctx := context.TODO()
id := webapps.NewRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionNamespaceName", "relayName")

read, err := client.GetHybridConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetHybridConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionNamespaceRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionNamespaceName", "relayName")

read, err := client.GetHybridConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceFunctionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName")

read, err := client.GetInstanceFunctionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceInfo`

```go
ctx := context.TODO()
id := webapps.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId")

read, err := client.GetInstanceInfo(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceInfoSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId")

read, err := client.GetInstanceInfoSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceMSDeployLog`

```go
ctx := context.TODO()
id := webapps.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId")

read, err := client.GetInstanceMSDeployLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceMSDeployLogSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId")

read, err := client.GetInstanceMSDeployLogSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceMsDeployStatus`

```go
ctx := context.TODO()
id := webapps.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId")

read, err := client.GetInstanceMsDeployStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceMsDeployStatusSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId")

read, err := client.GetInstanceMsDeployStatusSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceProcess`

```go
ctx := context.TODO()
id := webapps.NewInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId", "processId")

read, err := client.GetInstanceProcess(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceProcessDump`

```go
ctx := context.TODO()
id := webapps.NewInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId", "processId")

read, err := client.GetInstanceProcessDump(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceProcessDumpSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId", "processId")

read, err := client.GetInstanceProcessDumpSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceProcessModule`

```go
ctx := context.TODO()
id := webapps.NewInstanceProcessModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId", "processId", "moduleName")

read, err := client.GetInstanceProcessModule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceProcessModuleSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceProcessModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId", "processId", "moduleName")

read, err := client.GetInstanceProcessModuleSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceProcessSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId", "processId")

read, err := client.GetInstanceProcessSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetInstanceWorkflowSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "workflowName")

read, err := client.GetInstanceWorkflowSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetMSDeployLog`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetMSDeployLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetMSDeployLogSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetMSDeployLogSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetMSDeployStatus`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetMSDeployStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetMSDeployStatusSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetMSDeployStatusSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetMigrateMySqlStatus`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetMigrateMySqlStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetMigrateMySqlStatusSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetMigrateMySqlStatusSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetNetworkTraces`

```go
ctx := context.TODO()
id := webapps.NewNetworkTraceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "operationId")

read, err := client.GetNetworkTraces(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetNetworkTracesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotNetworkTraceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "operationId")

read, err := client.GetNetworkTracesSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetNetworkTracesSlotV2`

```go
ctx := context.TODO()
id := webapps.NewSiteSlotNetworkTraceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "operationId")

read, err := client.GetNetworkTracesSlotV2(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetNetworkTracesV2`

```go
ctx := context.TODO()
id := webapps.NewSiteNetworkTraceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "operationId")

read, err := client.GetNetworkTracesV2(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetOneDeployStatus`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetOneDeployStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPremierAddOn`

```go
ctx := context.TODO()
id := webapps.NewPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "premierAddonName")

read, err := client.GetPremierAddOn(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPremierAddOnSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "premierAddonName")

read, err := client.GetPremierAddOnSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPrivateAccess`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetPrivateAccess(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPrivateAccessSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetPrivateAccessSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPrivateEndpointConnection`

```go
ctx := context.TODO()
id := webapps.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "privateEndpointConnectionName")

read, err := client.GetPrivateEndpointConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPrivateEndpointConnectionList`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.GetPrivateEndpointConnectionList(ctx, id)` can be used to do batched pagination
items, err := client.GetPrivateEndpointConnectionListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.GetPrivateEndpointConnectionListSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.GetPrivateEndpointConnectionListSlot(ctx, id)` can be used to do batched pagination
items, err := client.GetPrivateEndpointConnectionListSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.GetPrivateEndpointConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "privateEndpointConnectionName")

read, err := client.GetPrivateEndpointConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPrivateLinkResources`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetPrivateLinkResources(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPrivateLinkResourcesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetPrivateLinkResourcesSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetProcess`

```go
ctx := context.TODO()
id := webapps.NewProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "processId")

read, err := client.GetProcess(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetProcessDump`

```go
ctx := context.TODO()
id := webapps.NewProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "processId")

read, err := client.GetProcessDump(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetProcessDumpSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "processId")

read, err := client.GetProcessDumpSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetProcessModule`

```go
ctx := context.TODO()
id := webapps.NewModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "processId", "moduleName")

read, err := client.GetProcessModule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetProcessModuleSlot`

```go
ctx := context.TODO()
id := webapps.NewProcessModuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "processId", "moduleName")

read, err := client.GetProcessModuleSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetProcessSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "processId")

read, err := client.GetProcessSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPublicCertificate`

```go
ctx := context.TODO()
id := webapps.NewPublicCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "publicCertificateName")

read, err := client.GetPublicCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetPublicCertificateSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPublicCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "publicCertificateName")

read, err := client.GetPublicCertificateSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetRelayServiceConnection`

```go
ctx := context.TODO()
id := webapps.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionName")

read, err := client.GetRelayServiceConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetRelayServiceConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionName")

read, err := client.GetRelayServiceConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetScmAllowed`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetScmAllowed(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetScmAllowedSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetScmAllowedSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSiteConnectionStringKeyVaultReference`

```go
ctx := context.TODO()
id := webapps.NewConnectionStringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "connectionStringKey")

read, err := client.GetSiteConnectionStringKeyVaultReference(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSiteConnectionStringKeyVaultReferenceSlot`

```go
ctx := context.TODO()
id := webapps.NewConfigReferenceConnectionStringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "connectionStringKey")

read, err := client.GetSiteConnectionStringKeyVaultReferenceSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSiteConnectionStringKeyVaultReferences`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.GetSiteConnectionStringKeyVaultReferences(ctx, id)` can be used to do batched pagination
items, err := client.GetSiteConnectionStringKeyVaultReferencesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.GetSiteConnectionStringKeyVaultReferencesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.GetSiteConnectionStringKeyVaultReferencesSlot(ctx, id)` can be used to do batched pagination
items, err := client.GetSiteConnectionStringKeyVaultReferencesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.GetSiteContainer`

```go
ctx := context.TODO()
id := webapps.NewSitecontainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "sitecontainerName")

read, err := client.GetSiteContainer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSiteContainerSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotSitecontainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "sitecontainerName")

read, err := client.GetSiteContainerSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSiteExtension`

```go
ctx := context.TODO()
id := webapps.NewSiteExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "siteExtensionId")

read, err := client.GetSiteExtension(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSiteExtensionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotSiteExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "siteExtensionId")

read, err := client.GetSiteExtensionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSitePhpErrorLogFlag`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetSitePhpErrorLogFlag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSitePhpErrorLogFlagSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetSitePhpErrorLogFlagSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSourceControl`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetSourceControl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSourceControlSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetSourceControlSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSwiftVirtualNetworkConnection`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetSwiftVirtualNetworkConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetSwiftVirtualNetworkConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetSwiftVirtualNetworkConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetTriggeredWebJob`

```go
ctx := context.TODO()
id := webapps.NewTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "triggeredWebJobName")

read, err := client.GetTriggeredWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetTriggeredWebJobHistory`

```go
ctx := context.TODO()
id := webapps.NewHistoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "triggeredWebJobName", "historyName")

read, err := client.GetTriggeredWebJobHistory(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetTriggeredWebJobHistorySlot`

```go
ctx := context.TODO()
id := webapps.NewTriggeredWebJobHistoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "triggeredWebJobName", "historyName")

read, err := client.GetTriggeredWebJobHistorySlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetTriggeredWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "triggeredWebJobName")

read, err := client.GetTriggeredWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetVnetConnection`

```go
ctx := context.TODO()
id := webapps.NewVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName")

read, err := client.GetVnetConnection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetVnetConnectionGateway`

```go
ctx := context.TODO()
id := webapps.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName", "gatewayName")

read, err := client.GetVnetConnectionGateway(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetVnetConnectionGatewaySlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName", "gatewayName")

read, err := client.GetVnetConnectionGatewaySlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetVnetConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName")

read, err := client.GetVnetConnectionSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetWebJob`

```go
ctx := context.TODO()
id := webapps.NewWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "webJobName")

read, err := client.GetWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "webJobName")

read, err := client.GetWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetWebSiteContainerLogs`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.GetWebSiteContainerLogs(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetWebSiteContainerLogsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.GetWebSiteContainerLogsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.GetWorkflow`

```go
ctx := context.TODO()
id := webapps.NewWorkflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "workflowName")

read, err := client.GetWorkflow(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.InstallSiteExtension`

```go
ctx := context.TODO()
id := webapps.NewSiteExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "siteExtensionId")

if err := client.InstallSiteExtensionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.InstallSiteExtensionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotSiteExtensionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "siteExtensionId")

if err := client.InstallSiteExtensionSlotThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.IsCloneable`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.IsCloneable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.IsCloneableSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.IsCloneableSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListApplicationSettings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListApplicationSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListApplicationSettingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListApplicationSettingsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListAzureStorageAccounts`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListAzureStorageAccounts(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListAzureStorageAccountsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListAzureStorageAccountsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListBackupStatusSecrets`

```go
ctx := context.TODO()
id := webapps.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "backupId")

payload := webapps.BackupRequest{
	// ...
}


read, err := client.ListBackupStatusSecrets(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListBackupStatusSecretsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "backupId")

payload := webapps.BackupRequest{
	// ...
}


read, err := client.ListBackupStatusSecretsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListBackups`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListBackups(ctx, id)` can be used to do batched pagination
items, err := client.ListBackupsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListBackupsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListBackupsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListBackupsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListBasicPublishingCredentialsPolicies`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListBasicPublishingCredentialsPolicies(ctx, id)` can be used to do batched pagination
items, err := client.ListBasicPublishingCredentialsPoliciesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListBasicPublishingCredentialsPoliciesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListBasicPublishingCredentialsPoliciesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListBasicPublishingCredentialsPoliciesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, webapps.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, webapps.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListConfigurationSnapshotInfo`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListConfigurationSnapshotInfo(ctx, id)` can be used to do batched pagination
items, err := client.ListConfigurationSnapshotInfoComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListConfigurationSnapshotInfoSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListConfigurationSnapshotInfoSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListConfigurationSnapshotInfoSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListConfigurations`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListConfigurations(ctx, id)` can be used to do batched pagination
items, err := client.ListConfigurationsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListConfigurationsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListConfigurationsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListConfigurationsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListConnectionStrings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListConnectionStrings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListConnectionStringsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListConnectionStringsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListContinuousWebJobs`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListContinuousWebJobs(ctx, id)` can be used to do batched pagination
items, err := client.ListContinuousWebJobsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListContinuousWebJobsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListContinuousWebJobsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListContinuousWebJobsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListDeploymentLog`

```go
ctx := context.TODO()
id := webapps.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "deploymentName")

read, err := client.ListDeploymentLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListDeploymentLogSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "deploymentName")

read, err := client.ListDeploymentLogSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListDeployments`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListDeployments(ctx, id)` can be used to do batched pagination
items, err := client.ListDeploymentsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListDeploymentsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListDeploymentsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListDeploymentsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListDomainOwnershipIdentifiers`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListDomainOwnershipIdentifiers(ctx, id)` can be used to do batched pagination
items, err := client.ListDomainOwnershipIdentifiersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListDomainOwnershipIdentifiersSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListDomainOwnershipIdentifiersSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListDomainOwnershipIdentifiersSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListFunctionKeys`

```go
ctx := context.TODO()
id := webapps.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName")

read, err := client.ListFunctionKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListFunctionKeysSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName")

read, err := client.ListFunctionKeysSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListFunctionSecrets`

```go
ctx := context.TODO()
id := webapps.NewFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "functionName")

read, err := client.ListFunctionSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListFunctionSecretsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "functionName")

read, err := client.ListFunctionSecretsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListFunctions`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListFunctions(ctx, id)` can be used to do batched pagination
items, err := client.ListFunctionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListHostKeys`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListHostKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListHostKeysSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListHostKeysSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListHostNameBindings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListHostNameBindings(ctx, id)` can be used to do batched pagination
items, err := client.ListHostNameBindingsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListHostNameBindingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListHostNameBindingsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListHostNameBindingsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListHybridConnections`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListHybridConnections(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListHybridConnectionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListHybridConnectionsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListInstanceFunctionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListInstanceFunctionsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceFunctionsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceIdentifiers`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListInstanceIdentifiers(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceIdentifiersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceIdentifiersSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListInstanceIdentifiersSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceIdentifiersSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceProcessModules`

```go
ctx := context.TODO()
id := webapps.NewInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId", "processId")

// alternatively `client.ListInstanceProcessModules(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceProcessModulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceProcessModulesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId", "processId")

// alternatively `client.ListInstanceProcessModulesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceProcessModulesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceProcessThreads`

```go
ctx := context.TODO()
id := webapps.NewInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId", "processId")

// alternatively `client.ListInstanceProcessThreads(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceProcessThreadsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceProcessThreadsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId", "processId")

// alternatively `client.ListInstanceProcessThreadsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceProcessThreadsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceProcesses`

```go
ctx := context.TODO()
id := webapps.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "instanceId")

// alternatively `client.ListInstanceProcesses(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceProcessesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceProcessesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "instanceId")

// alternatively `client.ListInstanceProcessesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceProcessesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListInstanceWorkflowsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListInstanceWorkflowsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListInstanceWorkflowsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListMetadata`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListMetadata(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListMetadataSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListMetadataSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListNetworkFeatures`

```go
ctx := context.TODO()
id := webapps.NewNetworkFeatureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "networkFeatureName")

read, err := client.ListNetworkFeatures(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListNetworkFeaturesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotNetworkFeatureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "networkFeatureName")

read, err := client.ListNetworkFeaturesSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListPerfMonCounters`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListPerfMonCounters(ctx, id, webapps.DefaultListPerfMonCountersOperationOptions())` can be used to do batched pagination
items, err := client.ListPerfMonCountersComplete(ctx, id, webapps.DefaultListPerfMonCountersOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListPerfMonCountersSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListPerfMonCountersSlot(ctx, id, webapps.DefaultListPerfMonCountersSlotOperationOptions())` can be used to do batched pagination
items, err := client.ListPerfMonCountersSlotComplete(ctx, id, webapps.DefaultListPerfMonCountersSlotOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListPremierAddOns`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListPremierAddOns(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListPremierAddOnsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListPremierAddOnsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListProcessModules`

```go
ctx := context.TODO()
id := webapps.NewProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "processId")

// alternatively `client.ListProcessModules(ctx, id)` can be used to do batched pagination
items, err := client.ListProcessModulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListProcessModulesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "processId")

// alternatively `client.ListProcessModulesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListProcessModulesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListProcessThreads`

```go
ctx := context.TODO()
id := webapps.NewProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "processId")

// alternatively `client.ListProcessThreads(ctx, id)` can be used to do batched pagination
items, err := client.ListProcessThreadsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListProcessThreadsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotProcessID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "processId")

// alternatively `client.ListProcessThreadsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListProcessThreadsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListProcesses`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListProcesses(ctx, id)` can be used to do batched pagination
items, err := client.ListProcessesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListProcessesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListProcessesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListProcessesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListProductionSiteDeploymentStatuses`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListProductionSiteDeploymentStatuses(ctx, id)` can be used to do batched pagination
items, err := client.ListProductionSiteDeploymentStatusesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListPublicCertificates`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListPublicCertificates(ctx, id)` can be used to do batched pagination
items, err := client.ListPublicCertificatesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListPublicCertificatesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListPublicCertificatesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListPublicCertificatesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListPublishingCredentials`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

if err := client.ListPublishingCredentialsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.ListPublishingCredentialsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

if err := client.ListPublishingCredentialsSlotThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.ListPublishingProfileXmlWithSecrets`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.CsmPublishingProfileOptions{
	// ...
}


read, err := client.ListPublishingProfileXmlWithSecrets(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListPublishingProfileXmlWithSecretsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.CsmPublishingProfileOptions{
	// ...
}


read, err := client.ListPublishingProfileXmlWithSecretsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListRelayServiceConnections`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListRelayServiceConnections(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListRelayServiceConnectionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListRelayServiceConnectionsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSiteBackups`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListSiteBackups(ctx, id)` can be used to do batched pagination
items, err := client.ListSiteBackupsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSiteBackupsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListSiteBackupsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListSiteBackupsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSiteContainers`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListSiteContainers(ctx, id)` can be used to do batched pagination
items, err := client.ListSiteContainersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSiteContainersSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListSiteContainersSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListSiteContainersSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSiteExtensions`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListSiteExtensions(ctx, id)` can be used to do batched pagination
items, err := client.ListSiteExtensionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSiteExtensionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListSiteExtensionsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListSiteExtensionsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSitePushSettings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListSitePushSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSitePushSettingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListSitePushSettingsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSlotConfigurationNames`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListSlotConfigurationNames(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSlotDifferencesFromProduction`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.CsmSlotEntity{
	// ...
}


// alternatively `client.ListSlotDifferencesFromProduction(ctx, id, payload)` can be used to do batched pagination
items, err := client.ListSlotDifferencesFromProductionComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSlotDifferencesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.CsmSlotEntity{
	// ...
}


// alternatively `client.ListSlotDifferencesSlot(ctx, id, payload)` can be used to do batched pagination
items, err := client.ListSlotDifferencesSlotComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSlotSiteDeploymentStatusesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListSlotSiteDeploymentStatusesSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListSlotSiteDeploymentStatusesSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSlots`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListSlots(ctx, id)` can be used to do batched pagination
items, err := client.ListSlotsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSnapshots`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListSnapshots(ctx, id)` can be used to do batched pagination
items, err := client.ListSnapshotsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSnapshotsFromDRSecondary`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListSnapshotsFromDRSecondary(ctx, id)` can be used to do batched pagination
items, err := client.ListSnapshotsFromDRSecondaryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSnapshotsFromDRSecondarySlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListSnapshotsFromDRSecondarySlot(ctx, id)` can be used to do batched pagination
items, err := client.ListSnapshotsFromDRSecondarySlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSnapshotsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListSnapshotsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListSnapshotsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListSyncFunctionTriggers`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListSyncFunctionTriggers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSyncFunctionTriggersSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListSyncFunctionTriggersSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSyncStatus`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListSyncStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListSyncStatusSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListSyncStatusSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListTriggeredWebJobHistory`

```go
ctx := context.TODO()
id := webapps.NewTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "triggeredWebJobName")

// alternatively `client.ListTriggeredWebJobHistory(ctx, id)` can be used to do batched pagination
items, err := client.ListTriggeredWebJobHistoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListTriggeredWebJobHistorySlot`

```go
ctx := context.TODO()
id := webapps.NewSlotTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "triggeredWebJobName")

// alternatively `client.ListTriggeredWebJobHistorySlot(ctx, id)` can be used to do batched pagination
items, err := client.ListTriggeredWebJobHistorySlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListTriggeredWebJobs`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListTriggeredWebJobs(ctx, id)` can be used to do batched pagination
items, err := client.ListTriggeredWebJobsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListTriggeredWebJobsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListTriggeredWebJobsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListTriggeredWebJobsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListUsages`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListUsages(ctx, id, webapps.DefaultListUsagesOperationOptions())` can be used to do batched pagination
items, err := client.ListUsagesComplete(ctx, id, webapps.DefaultListUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListUsagesSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListUsagesSlot(ctx, id, webapps.DefaultListUsagesSlotOperationOptions())` can be used to do batched pagination
items, err := client.ListUsagesSlotComplete(ctx, id, webapps.DefaultListUsagesSlotOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListVnetConnections`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListVnetConnections(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListVnetConnectionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListVnetConnectionsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListWebJobs`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListWebJobs(ctx, id)` can be used to do batched pagination
items, err := client.ListWebJobsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListWebJobsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

// alternatively `client.ListWebJobsSlot(ctx, id)` can be used to do batched pagination
items, err := client.ListWebJobsSlotComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListWorkflows`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

// alternatively `client.ListWorkflows(ctx, id)` can be used to do batched pagination
items, err := client.ListWorkflowsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WebAppsClient.ListWorkflowsConnections`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ListWorkflowsConnections(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ListWorkflowsConnectionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ListWorkflowsConnectionsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.MigrateMySql`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.MigrateMySqlRequest{
	// ...
}


if err := client.MigrateMySqlThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.MigrateStorage`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.StorageMigrationOptions{
	// ...
}


if err := client.MigrateStorageThenPoll(ctx, id, payload, webapps.DefaultMigrateStorageOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.PutPrivateAccessVnet`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.PrivateAccess{
	// ...
}


read, err := client.PutPrivateAccessVnet(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.PutPrivateAccessVnetSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.PrivateAccess{
	// ...
}


read, err := client.PutPrivateAccessVnetSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.RecoverSiteConfigurationSnapshot`

```go
ctx := context.TODO()
id := webapps.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "snapshotId")

read, err := client.RecoverSiteConfigurationSnapshot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.RecoverSiteConfigurationSnapshotSlot`

```go
ctx := context.TODO()
id := webapps.NewWebSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "snapshotId")

read, err := client.RecoverSiteConfigurationSnapshotSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ResetProductionSlotConfig`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.ResetProductionSlotConfig(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.ResetSlotConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.ResetSlotConfigurationSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Restart`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.Restart(ctx, id, webapps.DefaultRestartOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.RestartSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.RestartSlot(ctx, id, webapps.DefaultRestartSlotOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Restore`

```go
ctx := context.TODO()
id := webapps.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "backupId")

payload := webapps.RestoreRequest{
	// ...
}


if err := client.RestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreFromBackupBlob`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.RestoreRequest{
	// ...
}


if err := client.RestoreFromBackupBlobThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreFromBackupBlobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.RestoreRequest{
	// ...
}


if err := client.RestoreFromBackupBlobSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreFromDeletedApp`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.DeletedAppRestoreRequest{
	// ...
}


if err := client.RestoreFromDeletedAppThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreFromDeletedAppSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.DeletedAppRestoreRequest{
	// ...
}


if err := client.RestoreFromDeletedAppSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "backupId")

payload := webapps.RestoreRequest{
	// ...
}


if err := client.RestoreSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreSnapshot`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SnapshotRestoreRequest{
	// ...
}


if err := client.RestoreSnapshotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RestoreSnapshotSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SnapshotRestoreRequest{
	// ...
}


if err := client.RestoreSnapshotSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.RunTriggeredWebJob`

```go
ctx := context.TODO()
id := webapps.NewTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "triggeredWebJobName")

read, err := client.RunTriggeredWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.RunTriggeredWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotTriggeredWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "triggeredWebJobName")

read, err := client.RunTriggeredWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Start`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.Start(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StartContinuousWebJob`

```go
ctx := context.TODO()
id := webapps.NewContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "continuousWebJobName")

read, err := client.StartContinuousWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StartContinuousWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "continuousWebJobName")

read, err := client.StartContinuousWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StartNetworkTrace`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

if err := client.StartNetworkTraceThenPoll(ctx, id, webapps.DefaultStartNetworkTraceOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.StartNetworkTraceSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

if err := client.StartNetworkTraceSlotThenPoll(ctx, id, webapps.DefaultStartNetworkTraceSlotOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.StartSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.StartSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StartWebSiteNetworkTrace`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.StartWebSiteNetworkTrace(ctx, id, webapps.DefaultStartWebSiteNetworkTraceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StartWebSiteNetworkTraceOperation`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

if err := client.StartWebSiteNetworkTraceOperationThenPoll(ctx, id, webapps.DefaultStartWebSiteNetworkTraceOperationOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.StartWebSiteNetworkTraceOperationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

if err := client.StartWebSiteNetworkTraceOperationSlotThenPoll(ctx, id, webapps.DefaultStartWebSiteNetworkTraceOperationSlotOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.StartWebSiteNetworkTraceSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.StartWebSiteNetworkTraceSlot(ctx, id, webapps.DefaultStartWebSiteNetworkTraceSlotOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Stop`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.Stop(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopContinuousWebJob`

```go
ctx := context.TODO()
id := webapps.NewContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "continuousWebJobName")

read, err := client.StopContinuousWebJob(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopContinuousWebJobSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotContinuousWebJobID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "continuousWebJobName")

read, err := client.StopContinuousWebJobSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopNetworkTrace`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.StopNetworkTrace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopNetworkTraceSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.StopNetworkTraceSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.StopSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopWebSiteNetworkTrace`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.StopWebSiteNetworkTrace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.StopWebSiteNetworkTraceSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.StopWebSiteNetworkTraceSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.SwapSlotSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.CsmSlotEntity{
	// ...
}


if err := client.SwapSlotSlotThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.SwapSlotWithProduction`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.CsmSlotEntity{
	// ...
}


if err := client.SwapSlotWithProductionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WebAppsClient.SyncFunctionTriggers`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.SyncFunctionTriggers(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.SyncFunctionTriggersSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.SyncFunctionTriggersSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.SyncFunctions`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.SyncFunctions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.SyncFunctionsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.SyncFunctionsSlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.SyncRepository`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

read, err := client.SyncRepository(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.SyncRepositorySlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

read, err := client.SyncRepositorySlot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SitePatchResource{
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


### Example Usage: `WebAppsClient.UpdateApplicationSettings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.StringDictionary{
	// ...
}


read, err := client.UpdateApplicationSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateApplicationSettingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.StringDictionary{
	// ...
}


read, err := client.UpdateApplicationSettingsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateAuthSettings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteAuthSettings{
	// ...
}


read, err := client.UpdateAuthSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateAuthSettingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteAuthSettings{
	// ...
}


read, err := client.UpdateAuthSettingsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateAuthSettingsV2`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteAuthSettingsV2{
	// ...
}


read, err := client.UpdateAuthSettingsV2(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateAuthSettingsV2Slot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteAuthSettingsV2{
	// ...
}


read, err := client.UpdateAuthSettingsV2Slot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateAzureStorageAccounts`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.AzureStoragePropertyDictionaryResource{
	// ...
}


read, err := client.UpdateAzureStorageAccounts(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateAzureStorageAccountsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.AzureStoragePropertyDictionaryResource{
	// ...
}


read, err := client.UpdateAzureStorageAccountsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateBackupConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.BackupRequest{
	// ...
}


read, err := client.UpdateBackupConfiguration(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateBackupConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.BackupRequest{
	// ...
}


read, err := client.UpdateBackupConfigurationSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateConfiguration`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteConfigResource{
	// ...
}


read, err := client.UpdateConfiguration(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateConfigurationSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteConfigResource{
	// ...
}


read, err := client.UpdateConfigurationSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateConnectionStrings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.ConnectionStringDictionary{
	// ...
}


read, err := client.UpdateConnectionStrings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateConnectionStringsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.ConnectionStringDictionary{
	// ...
}


read, err := client.UpdateConnectionStringsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateDiagnosticLogsConfig`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteLogsConfig{
	// ...
}


read, err := client.UpdateDiagnosticLogsConfig(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateDiagnosticLogsConfigSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteLogsConfig{
	// ...
}


read, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateDomainOwnershipIdentifier`

```go
ctx := context.TODO()
id := webapps.NewSiteDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "domainOwnershipIdentifierName")

payload := webapps.Identifier{
	// ...
}


read, err := client.UpdateDomainOwnershipIdentifier(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateDomainOwnershipIdentifierSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotDomainOwnershipIdentifierID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "domainOwnershipIdentifierName")

payload := webapps.Identifier{
	// ...
}


read, err := client.UpdateDomainOwnershipIdentifierSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateFtpAllowed`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.CsmPublishingCredentialsPoliciesEntity{
	// ...
}


read, err := client.UpdateFtpAllowed(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateFtpAllowedSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.CsmPublishingCredentialsPoliciesEntity{
	// ...
}


read, err := client.UpdateFtpAllowedSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateHybridConnection`

```go
ctx := context.TODO()
id := webapps.NewRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionNamespaceName", "relayName")

payload := webapps.HybridConnection{
	// ...
}


read, err := client.UpdateHybridConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateHybridConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionNamespaceRelayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionNamespaceName", "relayName")

payload := webapps.HybridConnection{
	// ...
}


read, err := client.UpdateHybridConnectionSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateMetadata`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.StringDictionary{
	// ...
}


read, err := client.UpdateMetadata(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateMetadataSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.StringDictionary{
	// ...
}


read, err := client.UpdateMetadataSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdatePremierAddOn`

```go
ctx := context.TODO()
id := webapps.NewPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "premierAddonName")

payload := webapps.PremierAddOnPatchResource{
	// ...
}


read, err := client.UpdatePremierAddOn(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdatePremierAddOnSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotPremierAddonID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "premierAddonName")

payload := webapps.PremierAddOnPatchResource{
	// ...
}


read, err := client.UpdatePremierAddOnSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateRelayServiceConnection`

```go
ctx := context.TODO()
id := webapps.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "hybridConnectionName")

payload := webapps.RelayServiceConnectionEntity{
	// ...
}


read, err := client.UpdateRelayServiceConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateRelayServiceConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "hybridConnectionName")

payload := webapps.RelayServiceConnectionEntity{
	// ...
}


read, err := client.UpdateRelayServiceConnectionSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateScmAllowed`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.CsmPublishingCredentialsPoliciesEntity{
	// ...
}


read, err := client.UpdateScmAllowed(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateScmAllowedSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.CsmPublishingCredentialsPoliciesEntity{
	// ...
}


read, err := client.UpdateScmAllowedSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSitePushSettings`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.PushSettings{
	// ...
}


read, err := client.UpdateSitePushSettings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSitePushSettingsSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.PushSettings{
	// ...
}


read, err := client.UpdateSitePushSettingsSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SitePatchResource{
	// ...
}


read, err := client.UpdateSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSlotConfigurationNames`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SlotConfigNamesResource{
	// ...
}


read, err := client.UpdateSlotConfigurationNames(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSourceControl`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SiteSourceControl{
	// ...
}


read, err := client.UpdateSourceControl(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSourceControlSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SiteSourceControl{
	// ...
}


read, err := client.UpdateSourceControlSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSwiftVirtualNetworkConnectionWithCheck`

```go
ctx := context.TODO()
id := commonids.NewAppServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName")

payload := webapps.SwiftVirtualNetwork{
	// ...
}


read, err := client.UpdateSwiftVirtualNetworkConnectionWithCheck(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateSwiftVirtualNetworkConnectionWithCheckSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName")

payload := webapps.SwiftVirtualNetwork{
	// ...
}


read, err := client.UpdateSwiftVirtualNetworkConnectionWithCheckSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateVnetConnection`

```go
ctx := context.TODO()
id := webapps.NewVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName")

payload := webapps.VnetInfoResource{
	// ...
}


read, err := client.UpdateVnetConnection(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateVnetConnectionGateway`

```go
ctx := context.TODO()
id := webapps.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "virtualNetworkConnectionName", "gatewayName")

payload := webapps.VnetGateway{
	// ...
}


read, err := client.UpdateVnetConnectionGateway(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateVnetConnectionGatewaySlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName", "gatewayName")

payload := webapps.VnetGateway{
	// ...
}


read, err := client.UpdateVnetConnectionGatewaySlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WebAppsClient.UpdateVnetConnectionSlot`

```go
ctx := context.TODO()
id := webapps.NewSlotVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "siteName", "slotName", "virtualNetworkConnectionName")

payload := webapps.VnetInfoResource{
	// ...
}


read, err := client.UpdateVnetConnectionSlot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
