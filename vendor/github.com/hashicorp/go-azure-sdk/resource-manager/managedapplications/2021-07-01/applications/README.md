
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications` Documentation

The `applications` SDK allows for interaction with the Azure Resource Manager Service `managedapplications` (API Version `2021-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
```


### Client Initialization

```go
client := applications.NewApplicationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

payload := applications.Application{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.CreateOrUpdateById`

```go
ctx := context.TODO()
id := applications.NewApplicationIdID("applicationIdValue")

payload := applications.Application{
	// ...
}


if err := client.CreateOrUpdateByIdThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.Delete`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.DeleteById`

```go
ctx := context.TODO()
id := applications.NewApplicationIdID("applicationIdValue")

if err := client.DeleteByIdThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.Get`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.GetById`

```go
ctx := context.TODO()
id := applications.NewApplicationIdID("applicationIdValue")

read, err := client.GetById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.ListAllowedUpgradePlans`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

read, err := client.ListAllowedUpgradePlans(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := applications.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationsClient.ListBySubscription`

```go
ctx := context.TODO()
id := applications.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationsClient.ListTokens`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

payload := applications.ListTokenRequest{
	// ...
}


read, err := client.ListTokens(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationsClient.RefreshPermissions`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

if err := client.RefreshPermissionsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.Update`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

payload := applications.ApplicationPatchable{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.UpdateAccess`

```go
ctx := context.TODO()
id := applications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationValue")

payload := applications.UpdateAccessDefinition{
	// ...
}


if err := client.UpdateAccessThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationsClient.UpdateById`

```go
ctx := context.TODO()
id := applications.NewApplicationIdID("applicationIdValue")

payload := applications.ApplicationPatchable{
	// ...
}


if err := client.UpdateByIdThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
