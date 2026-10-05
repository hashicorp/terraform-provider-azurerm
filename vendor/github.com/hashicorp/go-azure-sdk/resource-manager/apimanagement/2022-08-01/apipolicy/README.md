
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy` Documentation

The `apipolicy` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
```


### Client Initialization

```go
client := apipolicy.NewApiPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiPolicyClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

payload := apipolicy.PolicyContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apipolicy.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.Delete`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

read, err := client.Delete(ctx, id, apipolicy.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.Get`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

read, err := client.Get(ctx, id, apipolicy.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.GetEntityTag`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.ListByApi`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.ListByApi(ctx, id)` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
