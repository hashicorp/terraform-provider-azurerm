
## `github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants` Documentation

The `tenants` SDK allows for interaction with the Azure Resource Manager Service `aadb2c` (API Version `2021-04-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
```


### Client Initialization

```go
client := tenants.NewTenantsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TenantsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := tenants.NewSubscriptionID()

payload := tenants.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantsClient.Create`

```go
ctx := context.TODO()
id := tenants.NewB2CDirectoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "directoryValue")

payload := tenants.CreateTenant{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TenantsClient.Delete`

```go
ctx := context.TODO()
id := tenants.NewB2CDirectoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "directoryValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TenantsClient.Get`

```go
ctx := context.TODO()
id := tenants.NewB2CDirectoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "directoryValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := tenants.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TenantsClient.ListBySubscription`

```go
ctx := context.TODO()
id := tenants.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TenantsClient.Update`

```go
ctx := context.TODO()
id := tenants.NewB2CDirectoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "directoryValue")

payload := tenants.UpdateTenant{
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
