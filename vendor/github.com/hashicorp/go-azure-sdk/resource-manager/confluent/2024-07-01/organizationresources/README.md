
## `github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources` Documentation

The `organizationresources` SDK allows for interaction with Azure Resource Manager `confluent` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources"
```


### Client Initialization

```go
client := organizationresources.NewOrganizationResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OrganizationResourcesClient.AccessCreateRoleBinding`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.AccessCreateRoleBindingRequestModel{
	// ...
}


read, err := client.AccessCreateRoleBinding(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessInviteUser`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.AccessInviteUserAccountModel{
	// ...
}


read, err := client.AccessInviteUser(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListClusters`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListClusters(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListEnvironments`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListEnvironments(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListInvitations`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListInvitations(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListRoleBindingNameList`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListRoleBindingNameList(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListRoleBindings`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListRoleBindings(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListServiceAccounts`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListServiceAccounts(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.AccessListUsers`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.AccessListUsers(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationCreate`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.OrganizationResource{
	// ...
}


if err := client.OrganizationCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationDelete`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

if err := client.OrganizationDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationGet`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

read, err := client.OrganizationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.OrganizationListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.OrganizationListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.OrganizationListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.OrganizationListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationListRegions`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.ListAccessRequestModel{
	// ...
}


read, err := client.OrganizationListRegions(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OrganizationResourcesClient.OrganizationUpdate`

```go
ctx := context.TODO()
id := organizationresources.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

payload := organizationresources.OrganizationResourceUpdate{
	// ...
}


read, err := client.OrganizationUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
