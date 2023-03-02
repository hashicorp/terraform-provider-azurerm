
## `github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/referencedatasets` Documentation

The `referencedatasets` SDK allows for interaction with the Azure Resource Manager Service `timeseriesinsights` (API Version `2020-05-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/referencedatasets"
```


### Client Initialization

```go
client := referencedatasets.NewReferenceDataSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReferenceDataSetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := referencedatasets.NewReferenceDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "referenceDataSetValue")

payload := referencedatasets.ReferenceDataSetCreateOrUpdateParameters{
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


### Example Usage: `ReferenceDataSetsClient.Delete`

```go
ctx := context.TODO()
id := referencedatasets.NewReferenceDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "referenceDataSetValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReferenceDataSetsClient.Get`

```go
ctx := context.TODO()
id := referencedatasets.NewReferenceDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "referenceDataSetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReferenceDataSetsClient.ListByEnvironment`

```go
ctx := context.TODO()
id := referencedatasets.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue")

read, err := client.ListByEnvironment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ReferenceDataSetsClient.Update`

```go
ctx := context.TODO()
id := referencedatasets.NewReferenceDataSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "referenceDataSetValue")

payload := referencedatasets.ReferenceDataSetUpdateParameters{
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
