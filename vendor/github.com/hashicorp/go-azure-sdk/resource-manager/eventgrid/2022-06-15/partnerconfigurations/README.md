
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerconfigurations` Documentation

The `partnerconfigurations` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerconfigurations"
```


### Client Initialization

```go
client := partnerconfigurations.NewPartnerConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PartnerConfigurationsClient.AuthorizePartner`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := partnerconfigurations.Partner{
	// ...
}


read, err := client.AuthorizePartner(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerConfigurationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := partnerconfigurations.PartnerConfiguration{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerConfigurationsClient.Get`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerConfigurationsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerConfigurationsClient.ListBySubscription`

```go
ctx := context.TODO()
id := partnerconfigurations.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, partnerconfigurations.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, partnerconfigurations.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerConfigurationsClient.UnauthorizePartner`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := partnerconfigurations.Partner{
	// ...
}


read, err := client.UnauthorizePartner(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerConfigurationsClient.Update`

```go
ctx := context.TODO()
id := partnerconfigurations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

payload := partnerconfigurations.PartnerConfigurationUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
