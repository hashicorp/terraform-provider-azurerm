
## `github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/smtpusernames` Documentation

The `smtpusernames` SDK allows for interaction with Azure Resource Manager `communication` (API Version `2026-03-18`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/smtpusernames"
```


### Client Initialization

```go
client := smtpusernames.NewSmtpUsernamesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SmtpUsernamesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := smtpusernames.NewSmtpUsernameID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceName", "smtpUsernameName")

payload := smtpusernames.SmtpUsernameResource{
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


### Example Usage: `SmtpUsernamesClient.Delete`

```go
ctx := context.TODO()
id := smtpusernames.NewSmtpUsernameID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceName", "smtpUsernameName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SmtpUsernamesClient.Get`

```go
ctx := context.TODO()
id := smtpusernames.NewSmtpUsernameID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceName", "smtpUsernameName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SmtpUsernamesClient.List`

```go
ctx := context.TODO()
id := smtpusernames.NewCommunicationServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationServiceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
