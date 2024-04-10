
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules` Documentation

The `datacollectionrules` SDK allows for interaction with the Azure Resource Manager Service `insights` (API Version `2022-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
```


### Client Initialization

```go
client := datacollectionrules.NewDataCollectionRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataCollectionRulesClient.Create`

```go
ctx := context.TODO()
id := datacollectionrules.NewDataCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionRuleValue")

payload := datacollectionrules.DataCollectionRuleResource{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionRulesClient.Delete`

```go
ctx := context.TODO()
id := datacollectionrules.NewDataCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionRuleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionRulesClient.Get`

```go
ctx := context.TODO()
id := datacollectionrules.NewDataCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionRulesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataCollectionRulesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataCollectionRulesClient.Update`

```go
ctx := context.TODO()
id := datacollectionrules.NewDataCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionRuleValue")

payload := datacollectionrules.ResourceForUpdate{
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
