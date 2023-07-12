
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/redis` Documentation

The `redis` SDK allows for interaction with the Azure Resource Manager Service `redis` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/redis"
```


### Client Initialization

```go
client := redis.NewRedisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := redis.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := redis.CheckNameAvailabilityParameters{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.Create`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.RedisCreateParameters{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisClient.Delete`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisClient.ExportData`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.ExportRDBParameters{
	// ...
}


if err := client.ExportDataThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisClient.FirewallRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := redis.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "firewallRuleValue")

payload := redis.RedisFirewallRule{
	// ...
}


read, err := client.FirewallRulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.FirewallRulesDelete`

```go
ctx := context.TODO()
id := redis.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "firewallRuleValue")

read, err := client.FirewallRulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.FirewallRulesGet`

```go
ctx := context.TODO()
id := redis.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "firewallRuleValue")

read, err := client.FirewallRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.FirewallRulesList`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

// alternatively `client.FirewallRulesList(ctx, id)` can be used to do batched pagination
items, err := client.FirewallRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisClient.ForceReboot`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.RedisRebootParameters{
	// ...
}


read, err := client.ForceReboot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.Get`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.ImportData`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.ImportRDBParameters{
	// ...
}


if err := client.ImportDataThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisClient.LinkedServerCreate`

```go
ctx := context.TODO()
id := redis.NewLinkedServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "linkedServerValue")

payload := redis.RedisLinkedServerCreateParameters{
	// ...
}


if err := client.LinkedServerCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisClient.LinkedServerDelete`

```go
ctx := context.TODO()
id := redis.NewLinkedServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "linkedServerValue")

if err := client.LinkedServerDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisClient.LinkedServerGet`

```go
ctx := context.TODO()
id := redis.NewLinkedServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue", "linkedServerValue")

read, err := client.LinkedServerGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.LinkedServerList`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

// alternatively `client.LinkedServerList(ctx, id)` can be used to do batched pagination
items, err := client.LinkedServerListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := redis.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisClient.ListBySubscription`

```go
ctx := context.TODO()
id := redis.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisClient.ListKeys`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.ListUpgradeNotifications`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

// alternatively `client.ListUpgradeNotifications(ctx, id, redis.DefaultListUpgradeNotificationsOperationOptions())` can be used to do batched pagination
items, err := client.ListUpgradeNotificationsComplete(ctx, id, redis.DefaultListUpgradeNotificationsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisClient.PatchSchedulesCreateOrUpdate`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.RedisPatchSchedule{
	// ...
}


read, err := client.PatchSchedulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.PatchSchedulesDelete`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

read, err := client.PatchSchedulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.PatchSchedulesGet`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

read, err := client.PatchSchedulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.PatchSchedulesListByRedisResource`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

// alternatively `client.PatchSchedulesListByRedisResource(ctx, id)` can be used to do batched pagination
items, err := client.PatchSchedulesListByRedisResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisClient.RegenerateKey`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.RedisRegenerateKeyParameters{
	// ...
}


read, err := client.RegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.Update`

```go
ctx := context.TODO()
id := redis.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisValue")

payload := redis.RedisUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
