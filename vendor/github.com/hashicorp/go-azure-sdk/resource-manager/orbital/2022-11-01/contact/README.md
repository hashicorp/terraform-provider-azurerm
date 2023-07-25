
## `github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contact` Documentation

The `contact` SDK allows for interaction with the Azure Resource Manager Service `orbital` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contact"
```


### Client Initialization

```go
client := contact.NewContactClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContactClient.Create`

```go
ctx := context.TODO()
id := contact.NewContactID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftValue", "contactValue")

payload := contact.Contact{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContactClient.Delete`

```go
ctx := context.TODO()
id := contact.NewContactID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftValue", "contactValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContactClient.Get`

```go
ctx := context.TODO()
id := contact.NewContactID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftValue", "contactValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContactClient.List`

```go
ctx := context.TODO()
id := contact.NewSpacecraftID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContactClient.SpacecraftsListAvailableContacts`

```go
ctx := context.TODO()
id := contact.NewSpacecraftID("12345678-1234-9876-4563-123456789012", "example-resource-group", "spacecraftValue")

payload := contact.ContactParameters{
	// ...
}


// alternatively `client.SpacecraftsListAvailableContacts(ctx, id, payload)` can be used to do batched pagination
items, err := client.SpacecraftsListAvailableContactsComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
