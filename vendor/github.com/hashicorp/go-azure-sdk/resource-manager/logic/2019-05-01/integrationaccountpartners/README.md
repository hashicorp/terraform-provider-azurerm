
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountpartners` Documentation

The `integrationaccountpartners` SDK allows for interaction with Azure Resource Manager `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountpartners"
```


### Client Initialization

```go
client := integrationaccountpartners.NewIntegrationAccountPartnersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountPartnersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountpartners.NewPartnerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "partnerName")

payload := integrationaccountpartners.IntegrationAccountPartner{
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


### Example Usage: `IntegrationAccountPartnersClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountpartners.NewPartnerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "partnerName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountPartnersClient.Get`

```go
ctx := context.TODO()
id := integrationaccountpartners.NewPartnerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "partnerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountPartnersClient.List`

```go
ctx := context.TODO()
id := integrationaccountpartners.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName")

// alternatively `client.List(ctx, id, integrationaccountpartners.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, integrationaccountpartners.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationAccountPartnersClient.ListContentCallbackURL`

```go
ctx := context.TODO()
id := integrationaccountpartners.NewPartnerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountName", "partnerName")

payload := integrationaccountpartners.GetCallbackURLParameters{
	// ...
}


read, err := client.ListContentCallbackURL(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
