
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/policy` Documentation

The `policy` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/policy"
```


### Client Initialization

```go
client := policy.NewPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := policy.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

payload := policy.PolicyContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, policy.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyClient.Delete`

```go
ctx := context.TODO()
id := policy.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.Delete(ctx, id, policy.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyClient.Get`

```go
ctx := context.TODO()
id := policy.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.Get(ctx, id, policy.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyClient.GetEntityTag`

```go
ctx := context.TODO()
id := policy.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyClient.ListByService`

```go
ctx := context.TODO()
id := policy.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

read, err := client.ListByService(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
