
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/commitmentplans` Documentation

The `commitmentplans` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/commitmentplans"
```


### Client Initialization

```go
client := commitmentplans.NewCommitmentPlansClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CommitmentPlansClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountCommitmentPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "commitmentPlanName")

payload := commitmentplans.CommitmentPlan{
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


### Example Usage: `CommitmentPlansClient.CreateOrUpdateAssociation`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "commitmentPlanName", "accountAssociationName")

payload := commitmentplans.CommitmentPlanAccountAssociation{
	// ...
}


if err := client.CreateOrUpdateAssociationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CommitmentPlansClient.Delete`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountCommitmentPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "commitmentPlanName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CommitmentPlansClient.DeleteAssociation`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "commitmentPlanName", "accountAssociationName")

if err := client.DeleteAssociationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CommitmentPlansClient.Get`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountCommitmentPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "commitmentPlanName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommitmentPlansClient.GetAssociation`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountAssociationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "commitmentPlanName", "accountAssociationName")

read, err := client.GetAssociation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CommitmentPlansClient.List`

```go
ctx := context.TODO()
id := commitmentplans.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CommitmentPlansClient.ListAssociations`

```go
ctx := context.TODO()
id := commitmentplans.NewCommitmentPlanID("12345678-1234-9876-4563-123456789012", "example-resource-group", "commitmentPlanName")

// alternatively `client.ListAssociations(ctx, id)` can be used to do batched pagination
items, err := client.ListAssociationsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
