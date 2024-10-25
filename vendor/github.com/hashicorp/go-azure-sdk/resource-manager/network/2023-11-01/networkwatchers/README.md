
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkwatchers` Documentation

The `networkwatchers` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkwatchers"
```


### Client Initialization

```go
client := networkwatchers.NewNetworkWatchersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkWatchersClient.CheckConnectivity`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.ConnectivityParameters{
	// ...
}


if err := client.CheckConnectivityThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.NetworkWatcher{
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


### Example Usage: `NetworkWatchersClient.Delete`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.Get`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkWatchersClient.GetAzureReachabilityReport`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.AzureReachabilityReportParameters{
	// ...
}


if err := client.GetAzureReachabilityReportThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.GetFlowLogStatus`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.FlowLogStatusParameters{
	// ...
}


if err := client.GetFlowLogStatusThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.GetNetworkConfigurationDiagnostic`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.NetworkConfigurationDiagnosticParameters{
	// ...
}


if err := client.GetNetworkConfigurationDiagnosticThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.GetNextHop`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.NextHopParameters{
	// ...
}


if err := client.GetNextHopThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.GetTopology`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.TopologyParameters{
	// ...
}


read, err := client.GetTopology(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkWatchersClient.GetTroubleshooting`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.TroubleshootingParameters{
	// ...
}


if err := client.GetTroubleshootingThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.GetTroubleshootingResult`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.QueryTroubleshootingParameters{
	// ...
}


if err := client.GetTroubleshootingResultThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.GetVMSecurityRules`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.SecurityGroupViewParameters{
	// ...
}


if err := client.GetVMSecurityRulesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkWatchersClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListAll(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkWatchersClient.ListAvailableProviders`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.AvailableProvidersListParameters{
	// ...
}


if err := client.ListAvailableProvidersThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.SetFlowLogConfiguration`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.FlowLogInformation{
	// ...
}


if err := client.SetFlowLogConfigurationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkWatchersClient.UpdateTags`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkWatchersClient.VerifyIPFlow`

```go
ctx := context.TODO()
id := networkwatchers.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := networkwatchers.VerificationIPFlowParameters{
	// ...
}


if err := client.VerifyIPFlowThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
