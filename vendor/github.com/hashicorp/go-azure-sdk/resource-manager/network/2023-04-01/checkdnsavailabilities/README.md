
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/checkdnsavailabilities` Documentation

The `checkdnsavailabilities` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/checkdnsavailabilities"
```


### Client Initialization

```go
client := checkdnsavailabilities.NewCheckDnsAvailabilitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CheckDnsAvailabilitiesClient.CheckDnsNameAvailability`

```go
ctx := context.TODO()
id := checkdnsavailabilities.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.CheckDnsNameAvailability(ctx, id, checkdnsavailabilities.DefaultCheckDnsNameAvailabilityOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
