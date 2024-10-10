
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/globalrulestack` Documentation

The `globalrulestack` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/globalrulestack"
```


### Client Initialization

```go
client := globalrulestack.NewGlobalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GlobalRulestackClient.Commit`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

if err := client.CommitThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

payload := globalrulestack.GlobalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalRulestackClient.Delete`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalRulestackClient.Get`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalRulestackClient.GetChangeLog`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

read, err := client.GetChangeLog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalRulestackClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx)` can be used to do batched pagination
items, err := client.ListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.ListAdvancedSecurityObjects`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.ListAdvancedSecurityObjects(ctx, id, globalrulestack.DefaultListAdvancedSecurityObjectsOperationOptions())` can be used to do batched pagination
items, err := client.ListAdvancedSecurityObjectsComplete(ctx, id, globalrulestack.DefaultListAdvancedSecurityObjectsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.ListAppIds`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.ListAppIds(ctx, id, globalrulestack.DefaultListAppIdsOperationOptions())` can be used to do batched pagination
items, err := client.ListAppIdsComplete(ctx, id, globalrulestack.DefaultListAppIdsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.ListCountries`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.ListCountries(ctx, id, globalrulestack.DefaultListCountriesOperationOptions())` can be used to do batched pagination
items, err := client.ListCountriesComplete(ctx, id, globalrulestack.DefaultListCountriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.ListFirewalls`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.ListFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.ListFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.ListPredefinedURLCategories`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.ListPredefinedURLCategories(ctx, id, globalrulestack.DefaultListPredefinedURLCategoriesOperationOptions())` can be used to do batched pagination
items, err := client.ListPredefinedURLCategoriesComplete(ctx, id, globalrulestack.DefaultListPredefinedURLCategoriesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.ListSecurityServices`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.ListSecurityServices(ctx, id, globalrulestack.DefaultListSecurityServicesOperationOptions())` can be used to do batched pagination
items, err := client.ListSecurityServicesComplete(ctx, id, globalrulestack.DefaultListSecurityServicesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalRulestackClient.Revert`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

read, err := client.Revert(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalRulestackClient.Update`

```go
ctx := context.TODO()
id := globalrulestack.NewGlobalRulestackID("globalRulestackName")

payload := globalrulestack.GlobalRulestackResourceUpdate{
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
