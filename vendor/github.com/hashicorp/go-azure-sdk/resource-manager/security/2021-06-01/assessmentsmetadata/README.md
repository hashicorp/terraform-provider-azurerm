
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata` Documentation

The `assessmentsmetadata` SDK allows for interaction with Azure Resource Manager `security` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
```


### Client Initialization

```go
client := assessmentsmetadata.NewAssessmentsMetadataClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssessmentsMetadataClient.CreateInSubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewProviderAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "assessmentMetadataName")

payload := assessmentsmetadata.SecurityAssessmentMetadataResponse{
	// ...
}


read, err := client.CreateInSubscription(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.DeleteInSubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewProviderAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "assessmentMetadataName")

read, err := client.DeleteInSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.Get`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewAssessmentMetadataID("assessmentMetadataName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.GetInSubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewProviderAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "assessmentMetadataName")

read, err := client.GetInSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx)` can be used to do batched pagination
items, err := client.ListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AssessmentsMetadataClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
