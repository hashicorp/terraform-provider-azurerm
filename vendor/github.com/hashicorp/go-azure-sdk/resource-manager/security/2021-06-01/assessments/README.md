
## `github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessments` Documentation

The `assessments` SDK allows for interaction with the Azure Resource Manager Service `security` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessments"
```


### Client Initialization

```go
client := assessments.NewAssessmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssessmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := assessments.NewScopedAssessmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "assessmentValue")

payload := assessments.SecurityAssessment{
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


### Example Usage: `AssessmentsClient.Delete`

```go
ctx := context.TODO()
id := assessments.NewScopedAssessmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "assessmentValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsClient.Get`

```go
ctx := context.TODO()
id := assessments.NewScopedAssessmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "assessmentValue")

read, err := client.Get(ctx, id, assessments.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssessmentsClient.List`

```go
ctx := context.TODO()
id := assessments.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
