
## `github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices` Documentation

The `dicomservices` SDK allows for interaction with Azure Resource Manager `healthcareapis` (API Version `2024-03-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
```


### Client Initialization

```go
client := dicomservices.NewDicomServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DicomServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dicomservices.NewDicomServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dicomServiceName")

payload := dicomservices.DicomService{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DicomServicesClient.Delete`

```go
ctx := context.TODO()
id := dicomservices.NewDicomServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dicomServiceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DicomServicesClient.Get`

```go
ctx := context.TODO()
id := dicomservices.NewDicomServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dicomServiceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DicomServicesClient.ListByWorkspace`

```go
ctx := context.TODO()
id := dicomservices.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DicomServicesClient.Update`

```go
ctx := context.TODO()
id := dicomservices.NewDicomServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "dicomServiceName")

payload := dicomservices.DicomServicePatchResource{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
