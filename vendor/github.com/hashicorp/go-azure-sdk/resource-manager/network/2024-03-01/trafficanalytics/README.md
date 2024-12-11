
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/trafficanalytics` Documentation

The `trafficanalytics` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/trafficanalytics"
```


### Client Initialization

```go
client := trafficanalytics.NewTrafficAnalyticsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TrafficAnalyticsClient.NetworkWatchersGetFlowLogStatus`

```go
ctx := context.TODO()
id := trafficanalytics.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := trafficanalytics.FlowLogStatusParameters{
	// ...
}


if err := client.NetworkWatchersGetFlowLogStatusThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TrafficAnalyticsClient.NetworkWatchersSetFlowLogConfiguration`

```go
ctx := context.TODO()
id := trafficanalytics.NewNetworkWatcherID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkWatcherName")

payload := trafficanalytics.FlowLogInformation{
	// ...
}


if err := client.NetworkWatchersSetFlowLogConfigurationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
