
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectglobalrulestack` Documentation

The `certificateobjectglobalrulestack` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2022-08-29`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/certificateobjectglobalrulestack"
```


### Client Initialization

```go
client := certificateobjectglobalrulestack.NewCertificateObjectGlobalRulestackClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateObjectGlobalRulestackClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestack.NewCertificateID("globalRulestackValue", "certificateValue")

payload := certificateobjectglobalrulestack.CertificateObjectGlobalRulestackResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectGlobalRulestackClient.Delete`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestack.NewCertificateID("globalRulestackValue", "certificateValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectGlobalRulestackClient.Get`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestack.NewCertificateID("globalRulestackValue", "certificateValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateObjectGlobalRulestackClient.List`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestack.NewGlobalRulestackID("globalRulestackValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
