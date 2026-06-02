
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/checknameavailabilitydisasterrecoveryconfigs` Documentation

The `checknameavailabilitydisasterrecoveryconfigs` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/checknameavailabilitydisasterrecoveryconfigs"
```


### Client Initialization

```go
client := checknameavailabilitydisasterrecoveryconfigs.NewCheckNameAvailabilityDisasterRecoveryConfigsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CheckNameAvailabilityDisasterRecoveryConfigsClient.DisasterRecoveryConfigsCheckNameAvailability`

```go
ctx := context.TODO()
id := checknameavailabilitydisasterrecoveryconfigs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityParameter{
	// ...
}


read, err := client.DisasterRecoveryConfigsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
