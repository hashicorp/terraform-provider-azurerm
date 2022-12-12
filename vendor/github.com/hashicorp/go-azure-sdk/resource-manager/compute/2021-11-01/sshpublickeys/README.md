
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/sshpublickeys` Documentation

The `sshpublickeys` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/sshpublickeys"
```


### Client Initialization

```go
client := sshpublickeys.NewSshPublicKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SshPublicKeysClient.Create`

```go
ctx := context.TODO()
id := sshpublickeys.NewSshPublicKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sshPublicKeyValue")

payload := sshpublickeys.SshPublicKeyResource{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SshPublicKeysClient.Delete`

```go
ctx := context.TODO()
id := sshpublickeys.NewSshPublicKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sshPublicKeyValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SshPublicKeysClient.GenerateKeyPair`

```go
ctx := context.TODO()
id := sshpublickeys.NewSshPublicKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sshPublicKeyValue")

read, err := client.GenerateKeyPair(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SshPublicKeysClient.Get`

```go
ctx := context.TODO()
id := sshpublickeys.NewSshPublicKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sshPublicKeyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SshPublicKeysClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := sshpublickeys.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SshPublicKeysClient.ListBySubscription`

```go
ctx := context.TODO()
id := sshpublickeys.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SshPublicKeysClient.Update`

```go
ctx := context.TODO()
id := sshpublickeys.NewSshPublicKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sshPublicKeyValue")

payload := sshpublickeys.SshPublicKeyUpdateResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
