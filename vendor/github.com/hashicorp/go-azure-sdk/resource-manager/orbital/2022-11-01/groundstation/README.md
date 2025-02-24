
## `github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/groundstation` Documentation

The `groundstation` SDK allows for interaction with Azure Resource Manager `orbital` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/groundstation"
```


### Client Initialization

```go
client := groundstation.NewGroundStationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GroundStationClient.AvailableGroundStationsListByCapability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AvailableGroundStationsListByCapability(ctx, id, groundstation.DefaultAvailableGroundStationsListByCapabilityOperationOptions())` can be used to do batched pagination
items, err := client.AvailableGroundStationsListByCapabilityComplete(ctx, id, groundstation.DefaultAvailableGroundStationsListByCapabilityOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
