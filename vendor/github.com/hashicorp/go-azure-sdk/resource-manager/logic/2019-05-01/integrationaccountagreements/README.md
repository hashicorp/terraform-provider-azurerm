
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements` Documentation

The `integrationaccountagreements` SDK allows for interaction with the Azure Resource Manager Service `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
```


### Client Initialization

```go
client := integrationaccountagreements.NewIntegrationAccountAgreementsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountAgreementsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountagreements.NewAgreementID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "agreementValue")

payload := integrationaccountagreements.IntegrationAccountAgreement{
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


### Example Usage: `IntegrationAccountAgreementsClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountagreements.NewAgreementID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "agreementValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountAgreementsClient.Get`

```go
ctx := context.TODO()
id := integrationaccountagreements.NewAgreementID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "agreementValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountAgreementsClient.List`

```go
ctx := context.TODO()
id := integrationaccountagreements.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

// alternatively `client.List(ctx, id, integrationaccountagreements.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, integrationaccountagreements.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationAccountAgreementsClient.ListContentCallbackUrl`

```go
ctx := context.TODO()
id := integrationaccountagreements.NewAgreementID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "agreementValue")

payload := integrationaccountagreements.GetCallbackUrlParameters{
	// ...
}


read, err := client.ListContentCallbackUrl(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
