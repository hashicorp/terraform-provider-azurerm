
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/securitypartnerproviders` Documentation

The `securitypartnerproviders` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/securitypartnerproviders"
```


### Client Initialization

```go
client := securitypartnerproviders.NewSecurityPartnerProvidersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityPartnerProvidersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securitypartnerproviders.NewSecurityPartnerProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "securityPartnerProviderValue")

payload := securitypartnerproviders.SecurityPartnerProvider{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityPartnerProvidersClient.Delete`

```go
ctx := context.TODO()
id := securitypartnerproviders.NewSecurityPartnerProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "securityPartnerProviderValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityPartnerProvidersClient.Get`

```go
ctx := context.TODO()
id := securitypartnerproviders.NewSecurityPartnerProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "securityPartnerProviderValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityPartnerProvidersClient.List`

```go
ctx := context.TODO()
id := securitypartnerproviders.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecurityPartnerProvidersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := securitypartnerproviders.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SecurityPartnerProvidersClient.UpdateTags`

```go
ctx := context.TODO()
id := securitypartnerproviders.NewSecurityPartnerProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "securityPartnerProviderValue")

payload := securitypartnerproviders.TagsObject{
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
