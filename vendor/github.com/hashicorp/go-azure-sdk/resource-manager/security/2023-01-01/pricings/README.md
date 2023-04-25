
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings` Documentation

The `pricings` SDK allows for interaction with the Azure Resource Manager Service `security` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
```


### Client Initialization

```go
client := pricings.NewPricingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PricingsClient.Get`

```go
ctx := context.TODO()
id := pricings.NewPricingID("12345678-1234-9876-4563-123456789012", "pricingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PricingsClient.List`

```go
ctx := context.TODO()
id := pricings.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PricingsClient.Update`

```go
ctx := context.TODO()
id := pricings.NewPricingID("12345678-1234-9876-4563-123456789012", "pricingValue")

payload := pricings.Pricing{
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
