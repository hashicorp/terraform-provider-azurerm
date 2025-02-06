
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/credential` Documentation

The `credential` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/credential"
```


### Client Initialization

```go
client := credential.NewCredentialClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CredentialClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := credential.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "credentialName")

payload := credential.CredentialCreateOrUpdateParameters{
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


### Example Usage: `CredentialClient.Delete`

```go
ctx := context.TODO()
id := credential.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "credentialName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CredentialClient.Get`

```go
ctx := context.TODO()
id := credential.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "credentialName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CredentialClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := credential.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CredentialClient.Update`

```go
ctx := context.TODO()
id := credential.NewCredentialID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "credentialName")

payload := credential.CredentialUpdateParameters{
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
