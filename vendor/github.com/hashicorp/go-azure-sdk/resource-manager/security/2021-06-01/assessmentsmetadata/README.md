
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata` Documentation

The `assessmentsmetadata` SDK allows for interaction with the Azure Resource Manager Service `security` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
```


### Client Initialization

```go
client := assessmentsmetadata.NewAssessmentsMetadataClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssessmentsMetadataClient.AssessmentsMetadataCreateInSubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewProviderAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "assessmentMetadataValue")

payload := assessmentsmetadata.SecurityAssessmentMetadataResponse{
	// ...
}


read, err := client.AssessmentsMetadataCreateInSubscription(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.AssessmentsMetadataDeleteInSubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewProviderAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "assessmentMetadataValue")

read, err := client.AssessmentsMetadataDeleteInSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.AssessmentsMetadataGet`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewAssessmentMetadataID("assessmentMetadataValue")

read, err := client.AssessmentsMetadataGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.AssessmentsMetadataGetInSubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewProviderAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "assessmentMetadataValue")

read, err := client.AssessmentsMetadataGetInSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsMetadataClient.AssessmentsMetadataList`

```go
ctx := context.TODO()


// alternatively `client.AssessmentsMetadataList(ctx)` can be used to do batched pagination
items, err := client.AssessmentsMetadataListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AssessmentsMetadataClient.AssessmentsMetadataListBySubscription`

```go
ctx := context.TODO()
id := assessmentsmetadata.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AssessmentsMetadataListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.AssessmentsMetadataListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
