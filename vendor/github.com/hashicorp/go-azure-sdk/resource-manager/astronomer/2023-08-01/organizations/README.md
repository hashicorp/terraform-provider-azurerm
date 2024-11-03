
## `github.com/hashicorp/go-azure-sdk/resource-manager/astronomer/2023-08-01/organizations` Documentation

The `organizations` SDK allows for interaction with Azure Resource Manager `astronomer` (API Version `2023-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/astronomer/2023-08-01/organizations"
```


### Client Initialization

```go
client := organizations.NewOrganizationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OrganizationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := organizations.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizations.OrganizationResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OrganizationsClient.Delete`

```go
ctx := context.TODO()
id := organizations.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OrganizationsClient.Get`

```go
ctx := context.TODO()
id := organizations.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationsClient.ListByResourceGroup`

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


### Example Usage: `OrganizationsClient.ListBySubscription`

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


### Example Usage: `OrganizationsClient.Update`

```go
ctx := context.TODO()
id := organizations.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizations.OrganizationResourceUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
