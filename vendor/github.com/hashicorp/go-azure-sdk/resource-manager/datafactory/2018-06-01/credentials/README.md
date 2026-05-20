
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials` Documentation

The `credentials` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
```


### Client Initialization

```go
client := credentials.NewCredentialsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CredentialsClient.CredentialOperationsCreateOrUpdate`

```go
ctx := context.TODO()
id := credentials.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "credentialName")

payload := credentials.CredentialResource{
	// ...
}


read, err := client.CredentialOperationsCreateOrUpdate(ctx, id, payload, credentials.DefaultCredentialOperationsCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CredentialsClient.CredentialOperationsDelete`

```go
ctx := context.TODO()
id := credentials.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "credentialName")

read, err := client.CredentialOperationsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CredentialsClient.CredentialOperationsGet`

```go
ctx := context.TODO()
id := credentials.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "credentialName")

read, err := client.CredentialOperationsGet(ctx, id, credentials.DefaultCredentialOperationsGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CredentialsClient.CredentialOperationsListByFactory`

```go
ctx := context.TODO()
id := credentials.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName")

// alternatively `client.CredentialOperationsListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.CredentialOperationsListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
