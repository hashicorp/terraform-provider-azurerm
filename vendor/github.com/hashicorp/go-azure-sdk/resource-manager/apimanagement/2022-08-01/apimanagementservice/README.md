
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice` Documentation

The `apimanagementservice` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
```


### Client Initialization

```go
client := apimanagementservice.NewApiManagementServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiManagementServiceClient.ApplyNetworkConfigurationUpdates`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := apimanagementservice.ApiManagementServiceApplyNetworkConfigurationParameters{
	// ...
}


if err := client.ApplyNetworkConfigurationUpdatesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiManagementServiceClient.Backup`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := apimanagementservice.ApiManagementServiceBackupRestoreParameters{
	// ...
}


if err := client.BackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiManagementServiceClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := apimanagementservice.ApiManagementServiceCheckNameAvailabilityParameters{
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


### Example Usage: `ApiManagementServiceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := apimanagementservice.ApiManagementServiceResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiManagementServiceClient.Delete`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApiManagementServiceClient.Get`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiManagementServiceClient.GetDomainOwnershipIdentifier`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.GetDomainOwnershipIdentifier(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiManagementServiceClient.GetSsoToken`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

read, err := client.GetSsoToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiManagementServiceClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiManagementServiceClient.ListByResourceGroup`

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


### Example Usage: `ApiManagementServiceClient.MigrateToStv2`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

if err := client.MigrateToStv2ThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApiManagementServiceClient.Restore`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := apimanagementservice.ApiManagementServiceBackupRestoreParameters{
	// ...
}


if err := client.RestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiManagementServiceClient.Update`

```go
ctx := context.TODO()
id := apimanagementservice.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := apimanagementservice.ApiManagementServiceUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
