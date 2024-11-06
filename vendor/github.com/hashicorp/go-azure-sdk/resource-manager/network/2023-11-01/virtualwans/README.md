
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans` Documentation

The `virtualwans` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
```


### Client Initialization

```go
client := virtualwans.NewVirtualWANsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualWANsClient.ConfigurationPolicyGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName", "configurationPolicyGroupName")

payload := virtualwans.VpnServerConfigurationPolicyGroup{
	// ...
}


if err := client.ConfigurationPolicyGroupsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.ConfigurationPolicyGroupsDelete`

```go
ctx := context.TODO()
id := virtualwans.NewConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName", "configurationPolicyGroupName")

if err := client.ConfigurationPolicyGroupsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.ConfigurationPolicyGroupsGet`

```go
ctx := context.TODO()
id := virtualwans.NewConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName", "configurationPolicyGroupName")

read, err := client.ConfigurationPolicyGroupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.ConfigurationPolicyGroupsListByVpnServerConfiguration`

```go
ctx := context.TODO()
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName")

// alternatively `client.ConfigurationPolicyGroupsListByVpnServerConfiguration(ctx, id)` can be used to do batched pagination
items, err := client.ConfigurationPolicyGroupsListByVpnServerConfigurationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.Generatevirtualwanvpnserverconfigurationvpnprofile`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

payload := virtualwans.VirtualWanVpnProfileParameters{
	// ...
}


if err := client.GeneratevirtualwanvpnserverconfigurationvpnprofileThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubRouteTablesCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "hubRouteTableName")

payload := virtualwans.HubRouteTable{
	// ...
}


if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubRouteTablesDelete`

```go
ctx := context.TODO()
id := virtualwans.NewHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "hubRouteTableName")

if err := client.HubRouteTablesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubRouteTablesGet`

```go
ctx := context.TODO()
id := virtualwans.NewHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "hubRouteTableName")

read, err := client.HubRouteTablesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.HubRouteTablesList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.HubRouteTablesList(ctx, id)` can be used to do batched pagination
items, err := client.HubRouteTablesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.HubVirtualNetworkConnectionsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewHubVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "hubVirtualNetworkConnectionName")

payload := virtualwans.HubVirtualNetworkConnection{
	// ...
}


if err := client.HubVirtualNetworkConnectionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubVirtualNetworkConnectionsDelete`

```go
ctx := context.TODO()
id := virtualwans.NewHubVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "hubVirtualNetworkConnectionName")

if err := client.HubVirtualNetworkConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubVirtualNetworkConnectionsGet`

```go
ctx := context.TODO()
id := virtualwans.NewHubVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "hubVirtualNetworkConnectionName")

read, err := client.HubVirtualNetworkConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.HubVirtualNetworkConnectionsList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.HubVirtualNetworkConnectionsList(ctx, id)` can be used to do batched pagination
items, err := client.HubVirtualNetworkConnectionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.NatRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "natRuleName")

payload := virtualwans.VpnGatewayNatRule{
	// ...
}


if err := client.NatRulesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.NatRulesDelete`

```go
ctx := context.TODO()
id := virtualwans.NewNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "natRuleName")

if err := client.NatRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.NatRulesGet`

```go
ctx := context.TODO()
id := virtualwans.NewNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "natRuleName")

read, err := client.NatRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.NatRulesListByVpnGateway`

```go
ctx := context.TODO()
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

// alternatively `client.NatRulesListByVpnGateway(ctx, id)` can be used to do batched pagination
items, err := client.NatRulesListByVpnGatewayComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.NetworkVirtualApplianceConnectionsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewNetworkVirtualApplianceConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "networkVirtualApplianceConnectionName")

payload := virtualwans.NetworkVirtualApplianceConnection{
	// ...
}


if err := client.NetworkVirtualApplianceConnectionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.NetworkVirtualApplianceConnectionsDelete`

```go
ctx := context.TODO()
id := virtualwans.NewNetworkVirtualApplianceConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "networkVirtualApplianceConnectionName")

if err := client.NetworkVirtualApplianceConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.NetworkVirtualApplianceConnectionsGet`

```go
ctx := context.TODO()
id := virtualwans.NewNetworkVirtualApplianceConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "networkVirtualApplianceConnectionName")

read, err := client.NetworkVirtualApplianceConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.NetworkVirtualApplianceConnectionsList`

```go
ctx := context.TODO()
id := virtualwans.NewNetworkVirtualApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName")

// alternatively `client.NetworkVirtualApplianceConnectionsList(ctx, id)` can be used to do batched pagination
items, err := client.NetworkVirtualApplianceConnectionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.P2sVpnGatewaysCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

payload := virtualwans.P2SVpnGateway{
	// ...
}


if err := client.P2sVpnGatewaysCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.P2sVpnGatewaysDelete`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

if err := client.P2sVpnGatewaysDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.P2sVpnGatewaysGet`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

read, err := client.P2sVpnGatewaysGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.P2sVpnGatewaysList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.P2sVpnGatewaysList(ctx, id)` can be used to do batched pagination
items, err := client.P2sVpnGatewaysListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.P2sVpnGatewaysListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.P2sVpnGatewaysListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.P2sVpnGatewaysListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.RouteMapsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewRouteMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routeMapName")

payload := virtualwans.RouteMap{
	// ...
}


if err := client.RouteMapsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.RouteMapsDelete`

```go
ctx := context.TODO()
id := virtualwans.NewRouteMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routeMapName")

if err := client.RouteMapsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.RouteMapsGet`

```go
ctx := context.TODO()
id := virtualwans.NewRouteMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routeMapName")

read, err := client.RouteMapsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.RouteMapsList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.RouteMapsList(ctx, id)` can be used to do batched pagination
items, err := client.RouteMapsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.RoutingIntentCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewRoutingIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routingIntentName")

payload := virtualwans.RoutingIntent{
	// ...
}


if err := client.RoutingIntentCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.RoutingIntentDelete`

```go
ctx := context.TODO()
id := virtualwans.NewRoutingIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routingIntentName")

if err := client.RoutingIntentDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.RoutingIntentGet`

```go
ctx := context.TODO()
id := virtualwans.NewRoutingIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routingIntentName")

read, err := client.RoutingIntentGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.RoutingIntentList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.RoutingIntentList(ctx, id)` can be used to do batched pagination
items, err := client.RoutingIntentListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.SupportedSecurityProviders`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

read, err := client.SupportedSecurityProviders(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.UpdateTags`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

payload := virtualwans.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "bgpConnectionName")

payload := virtualwans.BgpConnection{
	// ...
}


if err := client.VirtualHubBgpConnectionCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionDelete`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "bgpConnectionName")

if err := client.VirtualHubBgpConnectionDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionGet`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "bgpConnectionName")

read, err := client.VirtualHubBgpConnectionGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionsList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.VirtualHubBgpConnectionsList(ctx, id)` can be used to do batched pagination
items, err := client.VirtualHubBgpConnectionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionsListAdvertisedRoutes`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "bgpConnectionName")

if err := client.VirtualHubBgpConnectionsListAdvertisedRoutesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionsListLearnedRoutes`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "bgpConnectionName")

if err := client.VirtualHubBgpConnectionsListLearnedRoutesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubIPConfigurationCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "ipConfigurationName")

payload := virtualwans.HubIPConfiguration{
	// ...
}


if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubIPConfigurationDelete`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "ipConfigurationName")

if err := client.VirtualHubIPConfigurationDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubIPConfigurationGet`

```go
ctx := context.TODO()
id := commonids.NewVirtualHubIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "ipConfigurationName")

read, err := client.VirtualHubIPConfigurationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualHubIPConfigurationList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.VirtualHubIPConfigurationList(ctx, id)` can be used to do batched pagination
items, err := client.VirtualHubIPConfigurationListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VirtualHubRouteTableV2sCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routeTableName")

payload := virtualwans.VirtualHubRouteTableV2{
	// ...
}


if err := client.VirtualHubRouteTableV2sCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubRouteTableV2sDelete`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routeTableName")

if err := client.VirtualHubRouteTableV2sDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubRouteTableV2sGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName", "routeTableName")

read, err := client.VirtualHubRouteTableV2sGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualHubRouteTableV2sList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

// alternatively `client.VirtualHubRouteTableV2sList(ctx, id)` can be used to do batched pagination
items, err := client.VirtualHubRouteTableV2sListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

payload := virtualwans.VirtualHub{
	// ...
}


if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsDelete`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

if err := client.VirtualHubsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

read, err := client.VirtualHubsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsGetEffectiveVirtualHubRoutes`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

payload := virtualwans.EffectiveRoutesParameters{
	// ...
}


if err := client.VirtualHubsGetEffectiveVirtualHubRoutesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsGetInboundRoutes`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

payload := virtualwans.GetInboundRoutesParameters{
	// ...
}


if err := client.VirtualHubsGetInboundRoutesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsGetOutboundRoutes`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

payload := virtualwans.GetOutboundRoutesParameters{
	// ...
}


if err := client.VirtualHubsGetOutboundRoutesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.VirtualHubsList(ctx, id)` can be used to do batched pagination
items, err := client.VirtualHubsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.VirtualHubsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.VirtualHubsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsUpdateTags`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubName")

payload := virtualwans.TagsObject{
	// ...
}


read, err := client.VirtualHubsUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualWansCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

payload := virtualwans.VirtualWAN{
	// ...
}


if err := client.VirtualWansCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualWansDelete`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

if err := client.VirtualWansDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualWansGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

read, err := client.VirtualWansGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VirtualWansList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.VirtualWansList(ctx, id)` can be used to do batched pagination
items, err := client.VirtualWansListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VirtualWansListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.VirtualWansListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.VirtualWansListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName")

payload := virtualwans.VpnConnection{
	// ...
}


if err := client.VpnConnectionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsDelete`

```go
ctx := context.TODO()
id := commonids.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName")

if err := client.VpnConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsGet`

```go
ctx := context.TODO()
id := commonids.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName")

read, err := client.VpnConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsListByVpnGateway`

```go
ctx := context.TODO()
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

// alternatively `client.VpnConnectionsListByVpnGateway(ctx, id)` can be used to do batched pagination
items, err := client.VpnConnectionsListByVpnGatewayComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsStartPacketCapture`

```go
ctx := context.TODO()
id := commonids.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName")

payload := virtualwans.VpnConnectionPacketCaptureStartParameters{
	// ...
}


if err := client.VpnConnectionsStartPacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsStopPacketCapture`

```go
ctx := context.TODO()
id := commonids.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName")

payload := virtualwans.VpnConnectionPacketCaptureStopParameters{
	// ...
}


if err := client.VpnConnectionsStopPacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnGatewaysCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

payload := virtualwans.VpnGateway{
	// ...
}


if err := client.VpnGatewaysCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnGatewaysDelete`

```go
ctx := context.TODO()
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

if err := client.VpnGatewaysDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnGatewaysGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

read, err := client.VpnGatewaysGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VpnGatewaysList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.VpnGatewaysList(ctx, id)` can be used to do batched pagination
items, err := client.VpnGatewaysListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnGatewaysListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.VpnGatewaysListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.VpnGatewaysListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnLinkConnectionsGetIkeSas`

```go
ctx := context.TODO()
id := virtualwans.NewVpnLinkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName", "vpnLinkConnectionName")

if err := client.VpnLinkConnectionsGetIkeSasThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnLinkConnectionsListByVpnConnection`

```go
ctx := context.TODO()
id := commonids.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName")

// alternatively `client.VpnLinkConnectionsListByVpnConnection(ctx, id)` can be used to do batched pagination
items, err := client.VpnLinkConnectionsListByVpnConnectionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsAssociatedWithVirtualWanList`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

if err := client.VpnServerConfigurationsAssociatedWithVirtualWanListThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName")

payload := virtualwans.VpnServerConfiguration{
	// ...
}


if err := client.VpnServerConfigurationsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsDelete`

```go
ctx := context.TODO()
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName")

if err := client.VpnServerConfigurationsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName")

read, err := client.VpnServerConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.VpnServerConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.VpnServerConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.VpnServerConfigurationsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.VpnServerConfigurationsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnSiteLinkConnectionsGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnLinkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName", "vpnConnectionName", "vpnLinkConnectionName")

read, err := client.VpnSiteLinkConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VpnSiteLinksGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnSiteLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteName", "vpnSiteLinkName")

read, err := client.VpnSiteLinksGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VpnSiteLinksListByVpnSite`

```go
ctx := context.TODO()
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteName")

// alternatively `client.VpnSiteLinksListByVpnSite(ctx, id)` can be used to do batched pagination
items, err := client.VpnSiteLinksListByVpnSiteComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnSitesConfigurationDownload`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanName")

payload := virtualwans.GetVpnSitesConfigurationRequest{
	// ...
}


if err := client.VpnSitesConfigurationDownloadThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnSitesCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteName")

payload := virtualwans.VpnSite{
	// ...
}


if err := client.VpnSitesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnSitesDelete`

```go
ctx := context.TODO()
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteName")

if err := client.VpnSitesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnSitesGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteName")

read, err := client.VpnSitesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualWANsClient.VpnSitesList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.VpnSitesList(ctx, id)` can be used to do batched pagination
items, err := client.VpnSitesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualWANsClient.VpnSitesListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.VpnSitesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.VpnSitesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
