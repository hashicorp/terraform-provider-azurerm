
## `github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains` Documentation

The `domains` SDK allows for interaction with Azure Resource Manager `communication` (API Version `2023-03-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
```


### Client Initialization

```go
client := domains.NewDomainsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DomainsClient.CancelVerification`

```go
ctx := context.TODO()
id := domains.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

payload := domains.VerificationParameter{
	// ...
}


if err := client.CancelVerificationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DomainsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := domains.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

payload := domains.DomainResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DomainsClient.Delete`

```go
ctx := context.TODO()
id := domains.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DomainsClient.Get`

```go
ctx := context.TODO()
id := domains.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DomainsClient.InitiateVerification`

```go
ctx := context.TODO()
id := domains.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

payload := domains.VerificationParameter{
	// ...
}


if err := client.InitiateVerificationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DomainsClient.ListByEmailServiceResource`

```go
ctx := context.TODO()
id := domains.NewEmailServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName")

// alternatively `client.ListByEmailServiceResource(ctx, id)` can be used to do batched pagination
items, err := client.ListByEmailServiceResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DomainsClient.Update`

```go
ctx := context.TODO()
id := domains.NewDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "emailServiceName", "domainName")

payload := domains.UpdateDomainRequestParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
