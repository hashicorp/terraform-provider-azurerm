
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials` Documentation

The `federatedidentitycredentials` SDK allows for interaction with Azure Resource Manager `managedidentity` (API Version `2024-11-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
```


### Client Initialization

```go
client := federatedidentitycredentials.NewFederatedIdentityCredentialsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FederatedIdentityCredentialsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := federatedidentitycredentials.NewFederatedIdentityCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName", "federatedIdentityCredentialName")

payload := federatedidentitycredentials.FederatedIdentityCredential{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FederatedIdentityCredentialsClient.Delete`

```go
ctx := context.TODO()
id := federatedidentitycredentials.NewFederatedIdentityCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName", "federatedIdentityCredentialName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FederatedIdentityCredentialsClient.Get`

```go
ctx := context.TODO()
id := federatedidentitycredentials.NewFederatedIdentityCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName", "federatedIdentityCredentialName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FederatedIdentityCredentialsClient.List`

```go
ctx := context.TODO()
id := commonids.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName")

// alternatively `client.List(ctx, id, federatedidentitycredentials.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, federatedidentitycredentials.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
