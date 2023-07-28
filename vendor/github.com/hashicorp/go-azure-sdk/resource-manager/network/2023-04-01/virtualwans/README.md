
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans` Documentation

The `virtualwans` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans"
```


### Client Initialization

```go
client := virtualwans.NewVirtualWANsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualWANsClient.ConfigurationPolicyGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue", "configurationPolicyGroupValue")

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
id := virtualwans.NewConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue", "configurationPolicyGroupValue")

if err := client.ConfigurationPolicyGroupsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.ConfigurationPolicyGroupsGet`

```go
ctx := context.TODO()
id := virtualwans.NewConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue", "configurationPolicyGroupValue")

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
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

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
id := virtualwans.NewHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "hubRouteTableValue")

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
id := virtualwans.NewHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "hubRouteTableValue")

if err := client.HubRouteTablesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubRouteTablesGet`

```go
ctx := context.TODO()
id := virtualwans.NewHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "hubRouteTableValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewHubVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "hubVirtualNetworkConnectionValue")

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
id := virtualwans.NewHubVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "hubVirtualNetworkConnectionValue")

if err := client.HubVirtualNetworkConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.HubVirtualNetworkConnectionsGet`

```go
ctx := context.TODO()
id := virtualwans.NewHubVirtualNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "hubVirtualNetworkConnectionValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "natRuleValue")

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
id := virtualwans.NewNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "natRuleValue")

if err := client.NatRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.NatRulesGet`

```go
ctx := context.TODO()
id := virtualwans.NewNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "natRuleValue")

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
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue")

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
id := virtualwans.NewNetworkVirtualApplianceConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceValue", "networkVirtualApplianceConnectionValue")

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
id := virtualwans.NewNetworkVirtualApplianceConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceValue", "networkVirtualApplianceConnectionValue")

if err := client.NetworkVirtualApplianceConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.NetworkVirtualApplianceConnectionsGet`

```go
ctx := context.TODO()
id := virtualwans.NewNetworkVirtualApplianceConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceValue", "networkVirtualApplianceConnectionValue")

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
id := virtualwans.NewNetworkVirtualApplianceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceValue")

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
id := virtualwans.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayValue")

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
id := virtualwans.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayValue")

if err := client.P2sVpnGatewaysDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.P2sVpnGatewaysGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayValue")

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
id := virtualwans.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := virtualwans.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := virtualwans.NewRouteMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routeMapValue")

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
id := virtualwans.NewRouteMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routeMapValue")

if err := client.RouteMapsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.RouteMapsGet`

```go
ctx := context.TODO()
id := virtualwans.NewRouteMapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routeMapValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewRoutingIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routingIntentValue")

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
id := virtualwans.NewRoutingIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routingIntentValue")

if err := client.RoutingIntentDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.RoutingIntentGet`

```go
ctx := context.TODO()
id := virtualwans.NewRoutingIntentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routingIntentValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

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
id := virtualwans.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "bgpConnectionValue")

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
id := virtualwans.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "bgpConnectionValue")

if err := client.VirtualHubBgpConnectionDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "bgpConnectionValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "bgpConnectionValue")

if err := client.VirtualHubBgpConnectionsListAdvertisedRoutesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubBgpConnectionsListLearnedRoutes`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubBGPConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "bgpConnectionValue")

if err := client.VirtualHubBgpConnectionsListLearnedRoutesThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubIPConfigurationCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "ipConfigurationValue")

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
id := virtualwans.NewVirtualHubIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "ipConfigurationValue")

if err := client.VirtualHubIPConfigurationDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubIPConfigurationGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "ipConfigurationValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routeTableValue")

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
id := virtualwans.NewVirtualHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routeTableValue")

if err := client.VirtualHubRouteTableV2sDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubRouteTableV2sGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue", "routeTableValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

if err := client.VirtualHubsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualHubsGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := virtualwans.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := virtualwans.NewVirtualHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualHubValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

if err := client.VirtualWansDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VirtualWansGet`

```go
ctx := context.TODO()
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

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
id := virtualwans.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := virtualwans.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := virtualwans.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue")

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
id := virtualwans.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue")

if err := client.VpnConnectionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnConnectionsGet`

```go
ctx := context.TODO()
id := virtualwans.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue")

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
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue")

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
id := virtualwans.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue")

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
id := virtualwans.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue")

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
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue")

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
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue")

if err := client.VpnGatewaysDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnGatewaysGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue")

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
id := virtualwans.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := virtualwans.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := virtualwans.NewVpnLinkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue", "vpnLinkConnectionValue")

if err := client.VpnLinkConnectionsGetIkeSasThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnLinkConnectionsListByVpnConnection`

```go
ctx := context.TODO()
id := virtualwans.NewVPNConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

if err := client.VpnServerConfigurationsAssociatedWithVirtualWanListThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue")

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
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue")

if err := client.VpnServerConfigurationsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnServerConfigurationsGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationValue")

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
id := virtualwans.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := virtualwans.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := virtualwans.NewVpnLinkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayValue", "vpnConnectionValue", "vpnLinkConnectionValue")

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
id := virtualwans.NewVpnSiteLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteValue", "vpnSiteLinkValue")

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
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteValue")

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
id := virtualwans.NewVirtualWANID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualWanValue")

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
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteValue")

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
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteValue")

if err := client.VpnSitesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualWANsClient.VpnSitesGet`

```go
ctx := context.TODO()
id := virtualwans.NewVpnSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnSiteValue")

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
id := virtualwans.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := virtualwans.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.VpnSitesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.VpnSitesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
