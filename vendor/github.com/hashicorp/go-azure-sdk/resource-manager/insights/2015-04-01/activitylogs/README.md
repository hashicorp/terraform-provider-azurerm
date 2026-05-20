
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2015-04-01/activitylogs` Documentation

The `activitylogs` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2015-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2015-04-01/activitylogs"
```


### Client Initialization

```go
client := activitylogs.NewActivityLogsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ActivityLogsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, activitylogs.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, activitylogs.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
