
## `github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/linkedservices` Documentation

The `linkedservices` SDK allows for interaction with <unknown source data type> `synapse` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/linkedservices"
```


### Client Initialization

```go
client := linkedservices.NewLinkedServicesClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinkedServicesClient.LinkedServiceCreateOrUpdateLinkedService`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("linkedServiceName")

payload := linkedservices.LinkedServiceResource{
	// ...
}


if err := client.LinkedServiceCreateOrUpdateLinkedServiceThenPoll(ctx, id, payload, linkedservices.DefaultLinkedServiceCreateOrUpdateLinkedServiceOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `LinkedServicesClient.LinkedServiceDeleteLinkedService`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("linkedServiceName")

if err := client.LinkedServiceDeleteLinkedServiceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LinkedServicesClient.LinkedServiceGetLinkedService`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("linkedServiceName")

read, err := client.LinkedServiceGetLinkedService(ctx, id, linkedservices.DefaultLinkedServiceGetLinkedServiceOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedServicesClient.LinkedServiceGetLinkedServicesByWorkspace`

```go
ctx := context.TODO()


// alternatively `client.LinkedServiceGetLinkedServicesByWorkspace(ctx)` can be used to do batched pagination
items, err := client.LinkedServiceGetLinkedServicesByWorkspaceComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LinkedServicesClient.LinkedServiceRenameLinkedService`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("linkedServiceName")

payload := linkedservices.ArtifactRenameRequest{
	// ...
}


if err := client.LinkedServiceRenameLinkedServiceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
