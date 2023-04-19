
## `github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource` Documentation

The `resource` SDK allows for interaction with the Azure Resource Manager Service `healthcareapis` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
```


### Client Initialization

```go
client := resource.NewResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceClient.ServicesCreateOrUpdate`

```go
ctx := context.TODO()
id := resource.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

payload := resource.ServicesDescription{
	// ...
}


if err := client.ServicesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ResourceClient.ServicesDelete`

```go
ctx := context.TODO()
id := resource.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

if err := client.ServicesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ResourceClient.ServicesGet`

```go
ctx := context.TODO()
id := resource.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.ServicesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ResourceClient.ServicesUpdate`

```go
ctx := context.TODO()
id := resource.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

payload := resource.ServicesPatchDescription{
	// ...
}


if err := client.ServicesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
