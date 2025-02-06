
## `github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets` Documentation

The `budgets` SDK allows for interaction with Azure Resource Manager `consumption` (API Version `2019-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/consumption/2019-10-01/budgets"
```


### Client Initialization

```go
client := budgets.NewBudgetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BudgetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := budgets.NewScopedBudgetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "budgetName")

payload := budgets.Budget{
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


### Example Usage: `BudgetsClient.Delete`

```go
ctx := context.TODO()
id := budgets.NewScopedBudgetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "budgetName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BudgetsClient.Get`

```go
ctx := context.TODO()
id := budgets.NewScopedBudgetID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "budgetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BudgetsClient.List`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
