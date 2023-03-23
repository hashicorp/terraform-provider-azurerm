
## `github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders` Documentation

The `attestationproviders` SDK allows for interaction with the Azure Resource Manager Service `attestation` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
```


### Client Initialization

```go
client := attestationproviders.NewAttestationProvidersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AttestationProvidersClient.Create`

```go
ctx := context.TODO()
id := attestationproviders.NewAttestationProvidersID("12345678-1234-9876-4563-123456789012", "example-resource-group", "attestationProviderValue")

payload := attestationproviders.AttestationServiceCreationParams{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.Delete`

```go
ctx := context.TODO()
id := attestationproviders.NewAttestationProvidersID("12345678-1234-9876-4563-123456789012", "example-resource-group", "attestationProviderValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.Get`

```go
ctx := context.TODO()
id := attestationproviders.NewAttestationProvidersID("12345678-1234-9876-4563-123456789012", "example-resource-group", "attestationProviderValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.GetDefaultByLocation`

```go
ctx := context.TODO()
id := attestationproviders.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.GetDefaultByLocation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.List`

```go
ctx := context.TODO()
id := attestationproviders.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := attestationproviders.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.ListDefault`

```go
ctx := context.TODO()
id := attestationproviders.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListDefault(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AttestationProvidersClient.Update`

```go
ctx := context.TODO()
id := attestationproviders.NewAttestationProvidersID("12345678-1234-9876-4563-123456789012", "example-resource-group", "attestationProviderValue")

payload := attestationproviders.AttestationServicePatchParams{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
