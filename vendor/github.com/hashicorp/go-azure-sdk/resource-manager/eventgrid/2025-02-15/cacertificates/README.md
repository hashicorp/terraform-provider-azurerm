
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/cacertificates` Documentation

The `cacertificates` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2025-02-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/cacertificates"
```


### Client Initialization

```go
client := cacertificates.NewCaCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CaCertificatesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := cacertificates.NewCaCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "caCertificateName")

payload := cacertificates.CaCertificate{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CaCertificatesClient.Delete`

```go
ctx := context.TODO()
id := cacertificates.NewCaCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "caCertificateName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CaCertificatesClient.Get`

```go
ctx := context.TODO()
id := cacertificates.NewCaCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "caCertificateName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CaCertificatesClient.ListByNamespace`

```go
ctx := context.TODO()
id := cacertificates.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, cacertificates.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, cacertificates.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
