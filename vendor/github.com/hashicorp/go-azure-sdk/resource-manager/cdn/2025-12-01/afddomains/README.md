
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains` Documentation

The `afddomains` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
```


### Client Initialization

```go
client := afddomains.NewAFDDomainsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AFDDomainsClient.AFDCustomDomainsCreate`

```go
ctx := context.TODO()
id := afddomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

payload := afddomains.AFDDomain{
	// ...
}


if err := client.AFDCustomDomainsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AFDDomainsClient.AFDCustomDomainsDelete`

```go
ctx := context.TODO()
id := afddomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

if err := client.AFDCustomDomainsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDDomainsClient.AFDCustomDomainsGet`

```go
ctx := context.TODO()
id := afddomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

read, err := client.AFDCustomDomainsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AFDDomainsClient.AFDCustomDomainsListByProfile`

```go
ctx := context.TODO()
id := afddomains.NewProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName")

// alternatively `client.AFDCustomDomainsListByProfile(ctx, id)` can be used to do batched pagination
items, err := client.AFDCustomDomainsListByProfileComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AFDDomainsClient.AFDCustomDomainsRefreshValidationToken`

```go
ctx := context.TODO()
id := afddomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

if err := client.AFDCustomDomainsRefreshValidationTokenThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AFDDomainsClient.AFDCustomDomainsUpdate`

```go
ctx := context.TODO()
id := afddomains.NewCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "customDomainName")

payload := afddomains.AFDDomainUpdateParameters{
	// ...
}


if err := client.AFDCustomDomainsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
