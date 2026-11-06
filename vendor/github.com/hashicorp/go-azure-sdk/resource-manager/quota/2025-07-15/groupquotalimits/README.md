
## `github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotalimits` Documentation

The `groupquotalimits` SDK allows for interaction with Azure Resource Manager `quota` (API Version `2025-07-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotalimits"
```


### Client Initialization

```go
client := groupquotalimits.NewGroupQuotaLimitsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GroupQuotaLimitsClient.List`

```go
ctx := context.TODO()
id := groupquotalimits.NewGroupQuotaLimitID("managementGroupId", "groupQuotaName", "resourceProviderName", "groupQuotaLimitName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupQuotaLimitsClient.RequestUpdate`

```go
ctx := context.TODO()
id := groupquotalimits.NewGroupQuotaLimitID("managementGroupId", "groupQuotaName", "resourceProviderName", "groupQuotaLimitName")

payload := groupquotalimits.GroupQuotaLimitList{
	// ...
}


if err := client.RequestUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
