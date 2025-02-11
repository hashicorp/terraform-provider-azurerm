
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appservicecertificateorders` Documentation

The `appservicecertificateorders` SDK allows for interaction with Azure Resource Manager `web` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appservicecertificateorders"
```


### Client Initialization

```go
client := appservicecertificateorders.NewAppServiceCertificateOrdersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AppServiceCertificateOrdersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

payload := appservicecertificateorders.AppServiceCertificateOrder{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceCertificateOrdersClient.CreateOrUpdateCertificate`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName", "certificateName")

payload := appservicecertificateorders.AppServiceCertificateResource{
	// ...
}


if err := client.CreateOrUpdateCertificateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AppServiceCertificateOrdersClient.Delete`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.DeleteCertificate`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName", "certificateName")

read, err := client.DeleteCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.Get`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.GetCertificate`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName", "certificateName")

read, err := client.GetCertificate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceCertificateOrdersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceCertificateOrdersClient.ListCertificates`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

// alternatively `client.ListCertificates(ctx, id)` can be used to do batched pagination
items, err := client.ListCertificatesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AppServiceCertificateOrdersClient.Reissue`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

payload := appservicecertificateorders.ReissueCertificateOrderRequest{
	// ...
}


read, err := client.Reissue(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.Renew`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

payload := appservicecertificateorders.RenewCertificateOrderRequest{
	// ...
}


read, err := client.Renew(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.ResendEmail`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

read, err := client.ResendEmail(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.ResendRequestEmails`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

payload := appservicecertificateorders.NameIdentifier{
	// ...
}


read, err := client.ResendRequestEmails(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.RetrieveCertificateActions`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

read, err := client.RetrieveCertificateActions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.RetrieveCertificateEmailHistory`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

read, err := client.RetrieveCertificateEmailHistory(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.RetrieveSiteSeal`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

payload := appservicecertificateorders.SiteSealRequest{
	// ...
}


read, err := client.RetrieveSiteSeal(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.Update`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

payload := appservicecertificateorders.AppServiceCertificateOrderPatchResource{
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


### Example Usage: `AppServiceCertificateOrdersClient.UpdateCertificate`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName", "certificateName")

payload := appservicecertificateorders.AppServiceCertificatePatchResource{
	// ...
}


read, err := client.UpdateCertificate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.ValidatePurchaseInformation`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := appservicecertificateorders.AppServiceCertificateOrder{
	// ...
}


read, err := client.ValidatePurchaseInformation(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AppServiceCertificateOrdersClient.VerifyDomainOwnership`

```go
ctx := context.TODO()
id := appservicecertificateorders.NewCertificateOrderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "certificateOrderName")

read, err := client.VerifyDomainOwnership(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
