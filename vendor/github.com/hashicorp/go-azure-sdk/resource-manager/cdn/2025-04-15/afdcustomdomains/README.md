
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains` Documentation

The `afdcustomdomains` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-04-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
```


### Client Initialization

```go
client := afdcustomdomains.NewAFDCustomDomainsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AFDCustomDomainsClient.Create`

```go
ctx := context.TODO()
id := afdcustomdomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

payload := afdcustomdomains.AFDDomain{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDCustomDomainsClient.Delete`

```go
ctx := context.TODO()
id := afdcustomdomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDCustomDomainsClient.Get`

```go
ctx := context.TODO()
id := afdcustomdomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AFDCustomDomainsClient.ListByProfile`

```go
ctx := context.TODO()
id := afdcustomdomains.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.ListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDCustomDomainsClient.RefreshValidationToken`

```go
ctx := context.TODO()
id := afdcustomdomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

if err := client.RefreshValidationTokenThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDCustomDomainsClient.Update`

```go
ctx := context.TODO()
id := afdcustomdomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

payload := afdcustomdomains.AFDDomainUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
