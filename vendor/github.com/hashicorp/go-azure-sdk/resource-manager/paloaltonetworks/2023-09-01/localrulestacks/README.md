
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/localrulestacks` Documentation

The `localrulestacks` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/localrulestacks"
```


### Client Initialization

```go
client := localrulestacks.NewLocalRulestacksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalRulestacksClient.Commit`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

if err := client.CommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulestacksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

payload := localrulestacks.LocalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulestacksClient.Delete`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalRulestacksClient.Get`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestacksClient.GetChangeLog`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.GetChangeLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestacksClient.GetSupportInfo`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.GetSupportInfo(ctx, id, localrulestacks.DefaultGetSupportInfoOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestacksClient.ListAdvancedSecurityObjects`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListAdvancedSecurityObjects(ctx, id, localrulestacks.DefaultListAdvancedSecurityObjectsOperationOptions())` can be used to do batched pagination
items, err := client.ListAdvancedSecurityObjectsComplete(ctx, id, localrulestacks.DefaultListAdvancedSecurityObjectsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestacksClient.ListAppIds`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListAppIds(ctx, id, localrulestacks.DefaultListAppIdsOperationOptions())` can be used to do batched pagination
items, err := client.ListAppIdsComplete(ctx, id, localrulestacks.DefaultListAppIdsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestacksClient.ListByResourceGroup`

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


### Example Usage: `LocalRulestacksClient.ListBySubscription`

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


### Example Usage: `LocalRulestacksClient.ListCountries`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListCountries(ctx, id, localrulestacks.DefaultListCountriesOperationOptions())` can be used to do batched pagination
items, err := client.ListCountriesComplete(ctx, id, localrulestacks.DefaultListCountriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestacksClient.ListFirewalls`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.ListFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestacksClient.ListPredefinedURLCategories`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListPredefinedURLCategories(ctx, id, localrulestacks.DefaultListPredefinedURLCategoriesOperationOptions())` can be used to do batched pagination
items, err := client.ListPredefinedURLCategoriesComplete(ctx, id, localrulestacks.DefaultListPredefinedURLCategoriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestacksClient.ListSecurityServices`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

// alternatively `client.ListSecurityServices(ctx, id, localrulestacks.DefaultListSecurityServicesOperationOptions())` can be used to do batched pagination
items, err := client.ListSecurityServicesComplete(ctx, id, localrulestacks.DefaultListSecurityServicesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalRulestacksClient.Revert`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

read, err := client.Revert(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalRulestacksClient.Update`

```go
ctx := context.TODO()
id := localrulestacks.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackName")

payload := localrulestacks.LocalRulestackResourceUpdate{
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
