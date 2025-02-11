
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site` Documentation

The `site` SDK allows for interaction with Azure Resource Manager `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
```


### Client Initialization

```go
client := site.NewSiteClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SiteClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := site.NewSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "siteName")

payload := site.Site{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SiteClient.Delete`

```go
ctx := context.TODO()
id := site.NewSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "siteName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SiteClient.Get`

```go
ctx := context.TODO()
id := site.NewSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "siteName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SiteClient.UpdateTags`

```go
ctx := context.TODO()
id := site.NewSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "siteName")

payload := site.TagsObject{
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
