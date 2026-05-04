
## `github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links` Documentation

The `links` SDK allows for interaction with Azure Resource Manager `servicelinker` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
```


### Client Initialization

```go
client := links.NewLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinksClient.LinkerDelete`

```go
ctx := context.TODO()
id := links.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerName")

if err := client.LinkerDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LinksClient.LinkerListConfigurations`

```go
ctx := context.TODO()
id := links.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerName")

read, err := client.LinkerListConfigurations(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinksClient.LinkerUpdate`

```go
ctx := context.TODO()
id := links.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerName")

payload := links.LinkerPatch{
	// ...
}


if err := client.LinkerUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LinksClient.LinkerValidate`

```go
ctx := context.TODO()
id := links.NewScopedLinkerID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "linkerName")

if err := client.LinkerValidateThenPoll(ctx, id); err != nil {
	// handle the error
}
```
