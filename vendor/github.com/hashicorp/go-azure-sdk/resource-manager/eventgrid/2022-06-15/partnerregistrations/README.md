
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations` Documentation

The `partnerregistrations` SDK allows for interaction with the Azure Resource Manager Service `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations"
```


### Client Initialization

```go
client := partnerregistrations.NewPartnerRegistrationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PartnerRegistrationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := partnerregistrations.NewPartnerRegistrationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerRegistrationValue")

payload := partnerregistrations.PartnerRegistration{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerRegistrationsClient.Delete`

```go
ctx := context.TODO()
id := partnerregistrations.NewPartnerRegistrationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerRegistrationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PartnerRegistrationsClient.Get`

```go
ctx := context.TODO()
id := partnerregistrations.NewPartnerRegistrationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerRegistrationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PartnerRegistrationsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := partnerregistrations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, partnerregistrations.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, partnerregistrations.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerRegistrationsClient.ListBySubscription`

```go
ctx := context.TODO()
id := partnerregistrations.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, partnerregistrations.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, partnerregistrations.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PartnerRegistrationsClient.Update`

```go
ctx := context.TODO()
id := partnerregistrations.NewPartnerRegistrationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "partnerRegistrationValue")

payload := partnerregistrations.PartnerRegistrationUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
