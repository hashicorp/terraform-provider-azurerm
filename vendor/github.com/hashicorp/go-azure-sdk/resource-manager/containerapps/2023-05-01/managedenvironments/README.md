
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironments` Documentation

The `managedenvironments` SDK allows for interaction with the Azure Resource Manager Service `containerapps` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironments"
```


### Client Initialization

```go
client := managedenvironments.NewManagedEnvironmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedEnvironmentsClient.CertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := managedenvironments.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "certificateValue")

payload := managedenvironments.Certificate{
	// ...
}


read, err := client.CertificatesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.CertificatesDelete`

```go
ctx := context.TODO()
id := managedenvironments.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "certificateValue")

read, err := client.CertificatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.CertificatesGet`

```go
ctx := context.TODO()
id := managedenvironments.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "certificateValue")

read, err := client.CertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.CertificatesList`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

// alternatively `client.CertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.CertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedEnvironmentsClient.CertificatesUpdate`

```go
ctx := context.TODO()
id := managedenvironments.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "certificateValue")

payload := managedenvironments.CertificatePatch{
	// ...
}


read, err := client.CertificatesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

payload := managedenvironments.ManagedEnvironment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedEnvironmentsClient.Delete`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedEnvironmentsClient.DiagnosticsGetRoot`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

read, err := client.DiagnosticsGetRoot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.Get`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.GetAuthToken`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

read, err := client.GetAuthToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := managedenvironments.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedEnvironmentsClient.ListBySubscription`

```go
ctx := context.TODO()
id := managedenvironments.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedEnvironmentsClient.ListWorkloadProfileStates`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

// alternatively `client.ListWorkloadProfileStates(ctx, id)` can be used to do batched pagination
items, err := client.ListWorkloadProfileStatesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedCertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "managedCertificateValue")

payload := managedenvironments.ManagedCertificate{
	// ...
}


if err := client.ManagedCertificatesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedCertificatesDelete`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "managedCertificateValue")

read, err := client.ManagedCertificatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedCertificatesGet`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "managedCertificateValue")

read, err := client.ManagedCertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedCertificatesList`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

// alternatively `client.ManagedCertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.ManagedCertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedCertificatesUpdate`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "managedCertificateValue")

payload := managedenvironments.ManagedCertificatePatch{
	// ...
}


read, err := client.ManagedCertificatesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedEnvironmentDiagnosticsGetDetector`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentDetectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue", "detectorValue")

read, err := client.ManagedEnvironmentDiagnosticsGetDetector(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.ManagedEnvironmentDiagnosticsListDetectors`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

read, err := client.ManagedEnvironmentDiagnosticsListDetectors(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.NamespacesCheckNameAvailability`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

payload := managedenvironments.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.NamespacesCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedEnvironmentsClient.Update`

```go
ctx := context.TODO()
id := managedenvironments.NewManagedEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedEnvironmentValue")

payload := managedenvironments.ManagedEnvironment{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
