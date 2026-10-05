
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servervalidateestimatehighavailability` Documentation

The `servervalidateestimatehighavailability` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servervalidateestimatehighavailability"
```


### Client Initialization

```go
client := servervalidateestimatehighavailability.NewServerValidateEstimateHighAvailabilityClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerValidateEstimateHighAvailabilityClient.ServersValidateEstimateHighAvailability`

```go
ctx := context.TODO()
id := servervalidateestimatehighavailability.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := servervalidateestimatehighavailability.HighAvailabilityValidationEstimation{
	// ...
}


read, err := client.ServersValidateEstimateHighAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
