
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectglobalrulestackresources` Documentation

The `certificateobjectglobalrulestackresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectglobalrulestackresources"
```


### Client Initialization

```go
client := certificateobjectglobalrulestackresources.NewCertificateObjectGlobalRulestackResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CertificateObjectGlobalRulestackResourcesClient.CertificateObjectGlobalRulestackCreateOrUpdate`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestackresources.NewCertificateID("globalRulestackName", "certificateName")

payload := certificateobjectglobalrulestackresources.CertificateObjectGlobalRulestackResource{
	// ...
}


if err := client.CertificateObjectGlobalRulestackCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectGlobalRulestackResourcesClient.CertificateObjectGlobalRulestackDelete`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestackresources.NewCertificateID("globalRulestackName", "certificateName")

if err := client.CertificateObjectGlobalRulestackDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CertificateObjectGlobalRulestackResourcesClient.CertificateObjectGlobalRulestackGet`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestackresources.NewCertificateID("globalRulestackName", "certificateName")

read, err := client.CertificateObjectGlobalRulestackGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CertificateObjectGlobalRulestackResourcesClient.CertificateObjectGlobalRulestackList`

```go
ctx := context.TODO()
id := certificateobjectglobalrulestackresources.NewGlobalRulestackID("globalRulestackName")

// alternatively `client.CertificateObjectGlobalRulestackList(ctx, id)` can be used to do batched pagination
items, err := client.CertificateObjectGlobalRulestackListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
