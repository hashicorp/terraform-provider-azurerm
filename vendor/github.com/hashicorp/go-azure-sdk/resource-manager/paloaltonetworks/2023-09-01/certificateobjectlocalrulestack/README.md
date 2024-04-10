
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/certificateobjectlocalrulestack` Documentation

The `certificateobjectlocalrulestack` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2023-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/certificateobjectlocalrulestack"
```


### Client Initialization

```go
client := certificateobjectlocalrulestack.NewCertificateObjectLocalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateObjectLocalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestack.NewLocalRulestackCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "certificateValue")

payload := certificateobjectlocalrulestack.CertificateObjectLocalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectLocalRulestackClient.Delete`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestack.NewLocalRulestackCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "certificateValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectLocalRulestackClient.Get`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestack.NewLocalRulestackCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue", "certificateValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateObjectLocalRulestackClient.ListByLocalRulestacks`

```go
ctx := context.TODO()
id := certificateobjectlocalrulestack.NewLocalRulestackID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localRulestackValue")

// alternatively `client.ListByLocalRulestacks(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocalRulestacksComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
