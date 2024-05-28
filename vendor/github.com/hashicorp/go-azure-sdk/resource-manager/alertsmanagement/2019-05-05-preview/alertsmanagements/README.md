
## `github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-05-05-preview/alertsmanagements` Documentation

The `alertsmanagements` SDK allows for interaction with the Azure Resource Manager Service `alertsmanagement` (API Version `2019-05-05-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-05-05-preview/alertsmanagements"
```


### Client Initialization

```go
client := alertsmanagements.NewAlertsManagementsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AlertsManagementsClient.AlertsChangeState`

```go
ctx := context.TODO()
id := alertsmanagements.NewAlertID("12345678-1234-9876-4563-123456789012", "alertIdValue")

payload := alertsmanagements.Comments{
	// ...
}


read, err := client.AlertsChangeState(ctx, id, payload, alertsmanagements.DefaultAlertsChangeStateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementsClient.AlertsGetAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AlertsGetAll(ctx, id, alertsmanagements.DefaultAlertsGetAllOperationOptions())` can be used to do batched pagination
items, err := client.AlertsGetAllComplete(ctx, id, alertsmanagements.DefaultAlertsGetAllOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AlertsManagementsClient.AlertsGetById`

```go
ctx := context.TODO()
id := alertsmanagements.NewAlertID("12345678-1234-9876-4563-123456789012", "alertIdValue")

read, err := client.AlertsGetById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementsClient.AlertsGetHistory`

```go
ctx := context.TODO()
id := alertsmanagements.NewAlertID("12345678-1234-9876-4563-123456789012", "alertIdValue")

read, err := client.AlertsGetHistory(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementsClient.AlertsGetSummary`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.AlertsGetSummary(ctx, id, alertsmanagements.DefaultAlertsGetSummaryOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AlertsManagementsClient.AlertsMetaData`

```go
ctx := context.TODO()


read, err := client.AlertsMetaData(ctx, alertsmanagements.DefaultAlertsMetaDataOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
