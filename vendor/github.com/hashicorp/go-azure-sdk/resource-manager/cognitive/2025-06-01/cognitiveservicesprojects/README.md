
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects` Documentation

The `cognitiveservicesprojects` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
```


### Client Initialization

```go
client := cognitiveservicesprojects.NewCognitiveServicesProjectsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CognitiveServicesProjectsClient.ProjectsCreate`

```go
ctx := context.TODO()
id := cognitiveservicesprojects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName")

payload := cognitiveservicesprojects.Project{
	// ...
}


if err := client.ProjectsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CognitiveServicesProjectsClient.ProjectsDelete`

```go
ctx := context.TODO()
id := cognitiveservicesprojects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName")

if err := client.ProjectsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CognitiveServicesProjectsClient.ProjectsGet`

```go
ctx := context.TODO()
id := cognitiveservicesprojects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName")

read, err := client.ProjectsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesProjectsClient.ProjectsList`

```go
ctx := context.TODO()
id := cognitiveservicesprojects.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.ProjectsList(ctx, id)` can be used to do batched pagination
items, err := client.ProjectsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CognitiveServicesProjectsClient.ProjectsUpdate`

```go
ctx := context.TODO()
id := cognitiveservicesprojects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName")

payload := cognitiveservicesprojects.Project{
	// ...
}


if err := client.ProjectsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
