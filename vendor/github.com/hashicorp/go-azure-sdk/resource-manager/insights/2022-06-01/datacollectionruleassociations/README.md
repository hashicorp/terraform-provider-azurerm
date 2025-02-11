
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionruleassociations` Documentation

The `datacollectionruleassociations` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2022-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionruleassociations"
```


### Client Initialization

```go
client := datacollectionruleassociations.NewDataCollectionRuleAssociationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataCollectionRuleAssociationsClient.Create`

```go
ctx := context.TODO()
id := datacollectionruleassociations.NewScopedDataCollectionRuleAssociationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "dataCollectionRuleAssociationName")

payload := datacollectionruleassociations.DataCollectionRuleAssociationProxyOnlyResource{
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


### Example Usage: `DataCollectionRuleAssociationsClient.Delete`

```go
ctx := context.TODO()
id := datacollectionruleassociations.NewScopedDataCollectionRuleAssociationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "dataCollectionRuleAssociationName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionRuleAssociationsClient.Get`

```go
ctx := context.TODO()
id := datacollectionruleassociations.NewScopedDataCollectionRuleAssociationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "dataCollectionRuleAssociationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataCollectionRuleAssociationsClient.ListByDataCollectionEndpoint`

```go
ctx := context.TODO()
id := datacollectionruleassociations.NewDataCollectionEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionEndpointName")

// alternatively `client.ListByDataCollectionEndpoint(ctx, id)` can be used to do batched pagination
items, err := client.ListByDataCollectionEndpointComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataCollectionRuleAssociationsClient.ListByResource`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListByResource(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataCollectionRuleAssociationsClient.ListByRule`

```go
ctx := context.TODO()
id := datacollectionruleassociations.NewDataCollectionRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dataCollectionRuleName")

// alternatively `client.ListByRule(ctx, id)` can be used to do batched pagination
items, err := client.ListByRuleComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
