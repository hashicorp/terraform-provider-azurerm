
## `github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors` Documentation

The `azuretrafficcollectors` SDK allows for interaction with the Azure Resource Manager Service `networkfunction` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
```


### Client Initialization

```go
client := azuretrafficcollectors.NewAzureTrafficCollectorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AzureTrafficCollectorsClient.ByResourceGroupList`

```go
ctx := context.TODO()
id := azuretrafficcollectors.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ByResourceGroupList(ctx, id)` can be used to do batched pagination
items, err := client.ByResourceGroupListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AzureTrafficCollectorsClient.BySubscriptionList`

```go
ctx := context.TODO()
id := azuretrafficcollectors.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.BySubscriptionList(ctx, id)` can be used to do batched pagination
items, err := client.BySubscriptionListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AzureTrafficCollectorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := azuretrafficcollectors.NewAzureTrafficCollectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorValue")

payload := azuretrafficcollectors.AzureTrafficCollector{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AzureTrafficCollectorsClient.Delete`

```go
ctx := context.TODO()
id := azuretrafficcollectors.NewAzureTrafficCollectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AzureTrafficCollectorsClient.Get`

```go
ctx := context.TODO()
id := azuretrafficcollectors.NewAzureTrafficCollectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AzureTrafficCollectorsClient.UpdateTags`

```go
ctx := context.TODO()
id := azuretrafficcollectors.NewAzureTrafficCollectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "azureTrafficCollectorValue")

payload := azuretrafficcollectors.TagsObject{
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
