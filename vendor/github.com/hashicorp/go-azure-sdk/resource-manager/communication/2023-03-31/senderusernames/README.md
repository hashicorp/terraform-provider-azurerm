
## `github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/senderusernames` Documentation

The `senderusernames` SDK allows for interaction with Azure Resource Manager `communication` (API Version `2023-03-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/senderusernames"
```


### Client Initialization

```go
client := senderusernames.NewSenderUsernamesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SenderUsernamesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := senderusernames.NewSenderUsernameID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "senderUsernameName")

payload := senderusernames.SenderUsernameResource{
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


### Example Usage: `SenderUsernamesClient.Delete`

```go
ctx := context.TODO()
id := senderusernames.NewSenderUsernameID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "senderUsernameName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SenderUsernamesClient.Get`

```go
ctx := context.TODO()
id := senderusernames.NewSenderUsernameID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName", "senderUsernameName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SenderUsernamesClient.ListByDomains`

```go
ctx := context.TODO()
id := senderusernames.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

// alternatively `client.ListByDomains(ctx, id)` can be used to do batched pagination
items, err := client.ListByDomainsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
