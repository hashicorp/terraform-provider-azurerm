
## `github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices` Documentation

The `fhirservices` SDK allows for interaction with the Azure Resource Manager Service `healthcareapis` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
```


### Client Initialization

```go
client := fhirservices.NewFhirServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FhirServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fhirservices.NewFhirServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "fhirServiceValue")

payload := fhirservices.FhirService{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FhirServicesClient.Delete`

```go
ctx := context.TODO()
id := fhirservices.NewFhirServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "fhirServiceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FhirServicesClient.Get`

```go
ctx := context.TODO()
id := fhirservices.NewFhirServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "fhirServiceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FhirServicesClient.ListByWorkspace`

```go
ctx := context.TODO()
id := fhirservices.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FhirServicesClient.Update`

```go
ctx := context.TODO()
id := fhirservices.NewFhirServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "fhirServiceValue")

payload := fhirservices.FhirServicePatchResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
