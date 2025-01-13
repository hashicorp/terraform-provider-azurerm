
## `github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles` Documentation

The `configurationprofiles` SDK allows for interaction with Azure Resource Manager `automanage` (API Version `2022-05-04`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
```


### Client Initialization

```go
client := configurationprofiles.NewConfigurationProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfigurationProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := configurationprofiles.NewConfigurationProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationProfileName")

payload := configurationprofiles.ConfigurationProfile{
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


### Example Usage: `ConfigurationProfilesClient.Delete`

```go
ctx := context.TODO()
id := configurationprofiles.NewConfigurationProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationProfileName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfilesClient.Get`

```go
ctx := context.TODO()
id := configurationprofiles.NewConfigurationProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationProfileName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfilesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfilesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfigurationProfilesClient.Update`

```go
ctx := context.TODO()
id := configurationprofiles.NewConfigurationProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "configurationProfileName")

payload := configurationprofiles.ConfigurationProfileUpdate{
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
