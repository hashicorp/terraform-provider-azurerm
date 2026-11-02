
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains` Documentation

The `customdomains` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
```


### Client Initialization

```go
client := customdomains.NewCustomDomainsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CustomDomainsClient.Create`

```go
ctx := context.TODO()
id := customdomains.NewEndpointCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "customDomainName")

payload := customdomains.CustomDomainParameters{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CustomDomainsClient.Delete`

```go
ctx := context.TODO()
id := customdomains.NewEndpointCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "customDomainName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CustomDomainsClient.DisableCustomHTTPS`

```go
ctx := context.TODO()
id := customdomains.NewEndpointCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "customDomainName")

if err := client.DisableCustomHTTPSThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CustomDomainsClient.EnableCustomHTTPS`

```go
ctx := context.TODO()
id := customdomains.NewEndpointCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "customDomainName")

payload := customdomains.CustomDomainHTTPSParameters{
	// ...
}


if err := client.EnableCustomHTTPSThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CustomDomainsClient.Get`

```go
ctx := context.TODO()
id := customdomains.NewEndpointCustomDomainID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "customDomainName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CustomDomainsClient.ListByEndpoint`

```go
ctx := context.TODO()
id := customdomains.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

// alternatively `client.ListByEndpoint(ctx, id)` can be used to do batched pagination
items, err := client.ListByEndpointComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
