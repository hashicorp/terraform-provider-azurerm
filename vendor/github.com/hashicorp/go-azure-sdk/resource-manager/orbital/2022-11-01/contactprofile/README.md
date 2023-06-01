
## `github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile` Documentation

The `contactprofile` SDK allows for interaction with the Azure Resource Manager Service `orbital` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile"
```


### Client Initialization

```go
client := contactprofile.NewContactProfileClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContactProfileClient.ContactProfilesCreateOrUpdate`

```go
ctx := context.TODO()
id := contactprofile.NewContactProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "contactProfileValue")

payload := contactprofile.ContactProfile{
	// ...
}


if err := client.ContactProfilesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContactProfileClient.ContactProfilesDelete`

```go
ctx := context.TODO()
id := contactprofile.NewContactProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "contactProfileValue")

if err := client.ContactProfilesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContactProfileClient.ContactProfilesGet`

```go
ctx := context.TODO()
id := contactprofile.NewContactProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "contactProfileValue")

read, err := client.ContactProfilesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContactProfileClient.ContactProfilesList`

```go
ctx := context.TODO()
id := contactprofile.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ContactProfilesList(ctx, id)` can be used to do batched pagination
items, err := client.ContactProfilesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContactProfileClient.ContactProfilesListBySubscription`

```go
ctx := context.TODO()
id := contactprofile.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ContactProfilesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ContactProfilesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContactProfileClient.ContactProfilesUpdateTags`

```go
ctx := context.TODO()
id := contactprofile.NewContactProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "contactProfileValue")

payload := contactprofile.TagsObject{
	// ...
}


if err := client.ContactProfilesUpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
