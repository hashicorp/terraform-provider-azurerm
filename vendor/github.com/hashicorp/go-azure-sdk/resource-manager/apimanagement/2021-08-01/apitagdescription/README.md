
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apitagdescription` Documentation

The `apitagdescription` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apitagdescription"
```


### Client Initialization

```go
client := apitagdescription.NewApiTagDescriptionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiTagDescriptionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apitagdescription.NewTagDescriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagDescriptionIdValue")

payload := apitagdescription.TagDescriptionCreateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apitagdescription.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagDescriptionClient.Delete`

```go
ctx := context.TODO()
id := apitagdescription.NewTagDescriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagDescriptionIdValue")

read, err := client.Delete(ctx, id, apitagdescription.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagDescriptionClient.Get`

```go
ctx := context.TODO()
id := apitagdescription.NewTagDescriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagDescriptionIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagDescriptionClient.GetEntityTag`

```go
ctx := context.TODO()
id := apitagdescription.NewTagDescriptionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagDescriptionIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagDescriptionClient.ListByService`

```go
ctx := context.TODO()
id := apitagdescription.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

// alternatively `client.ListByService(ctx, id, apitagdescription.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, apitagdescription.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
