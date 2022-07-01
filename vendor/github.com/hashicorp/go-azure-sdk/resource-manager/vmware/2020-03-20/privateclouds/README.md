
## `github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/privateclouds` Documentation

The `privateclouds` SDK allows for interaction with the Azure Resource Manager Service `vmware` (API Version `2020-03-20`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2020-03-20/privateclouds"
```


### Client Initialization

```go
client := privateclouds.NewPrivateCloudsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateCloudsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privateclouds.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

payload := privateclouds.PrivateCloud{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateCloudsClient.Delete`

```go
ctx := context.TODO()
id := privateclouds.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateCloudsClient.Get`

```go
ctx := context.TODO()
id := privateclouds.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateCloudsClient.List`

```go
ctx := context.TODO()
id := privateclouds.NewResourceGroupID()

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateCloudsClient.ListAdminCredentials`

```go
ctx := context.TODO()
id := privateclouds.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

read, err := client.ListAdminCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateCloudsClient.ListInSubscription`

```go
ctx := context.TODO()
id := privateclouds.NewSubscriptionID()

// alternatively `client.ListInSubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListInSubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateCloudsClient.Update`

```go
ctx := context.TODO()
id := privateclouds.NewPrivateCloudID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateCloudValue")

payload := privateclouds.PrivateCloudUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
