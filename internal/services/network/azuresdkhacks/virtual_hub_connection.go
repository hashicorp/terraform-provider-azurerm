package azuresdkhacks

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

// workaround for https://github.com/Azure/azure-rest-api-specs/issues/23920
// TODO: check if it could be removed in 4.0

type HubVirtualNetworkConnectionsClient struct {
	OriginalClient *network.HubVirtualNetworkConnectionsClient
}

func (client HubVirtualNetworkConnectionsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, virtualHubName string, connectionName string, hubVirtualNetworkConnectionParameters HubVirtualNetworkConnection) (result network.HubVirtualNetworkConnectionsCreateOrUpdateFuture, err error) {
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, virtualHubName, connectionName, hubVirtualNetworkConnectionParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.HubVirtualNetworkConnectionsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.OriginalClient.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.HubVirtualNetworkConnectionsClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}

	return
}

func (client HubVirtualNetworkConnectionsClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, virtualHubName string, connectionName string, hubVirtualNetworkConnectionParameters HubVirtualNetworkConnection) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"connectionName":    autorest.Encode("path", connectionName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.OriginalClient.SubscriptionID),
		"virtualHubName":    autorest.Encode("path", virtualHubName),
	}

	const APIVersion = "2022-07-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	hubVirtualNetworkConnectionParameters.Etag = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.OriginalClient.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualHubs/{virtualHubName}/hubVirtualNetworkConnections/{connectionName}", pathParameters),
		autorest.WithJSON(hubVirtualNetworkConnectionParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

type StaticRoutesConfig struct {
	PropagateStaticRoutes          *bool                                  `json:"propagateStaticRoutes,omitempty"`
	VnetLocalRouteOverrideCriteria network.VnetLocalRouteOverrideCriteria `json:"vnetLocalRouteOverrideCriteria,omitempty"`
}

func (src StaticRoutesConfig) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if src.VnetLocalRouteOverrideCriteria != "" {
		objectMap["vnetLocalRouteOverrideCriteria"] = src.VnetLocalRouteOverrideCriteria
	}
	if src.PropagateStaticRoutes != nil {
		objectMap["propagateStaticRoutes"] = src.PropagateStaticRoutes
	}
	return json.Marshal(objectMap)
}

type VnetRoute struct {
	StaticRoutesConfig *StaticRoutesConfig    `json:"staticRoutesConfig,omitempty"`
	StaticRoutes       *[]network.StaticRoute `json:"staticRoutes,omitempty"`
	BgpConnections     *[]network.SubResource `json:"bgpConnections,omitempty"`
}

func (vr VnetRoute) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if vr.StaticRoutesConfig != nil {
		objectMap["staticRoutesConfig"] = vr.StaticRoutesConfig
	}
	if vr.StaticRoutes != nil {
		objectMap["staticRoutes"] = vr.StaticRoutes
	}
	return json.Marshal(objectMap)
}

type RoutingConfiguration struct {
	AssociatedRouteTable  *network.SubResource          `json:"associatedRouteTable,omitempty"`
	PropagatedRouteTables *network.PropagatedRouteTable `json:"propagatedRouteTables,omitempty"`
	VnetRoutes            *VnetRoute                    `json:"vnetRoutes,omitempty"`
	InboundRouteMap       *network.SubResource          `json:"inboundRouteMap,omitempty"`
	OutboundRouteMap      *network.SubResource          `json:"outboundRouteMap,omitempty"`
}

type HubVirtualNetworkConnectionProperties struct {
	RemoteVirtualNetwork                *network.SubResource      `json:"remoteVirtualNetwork,omitempty"`
	AllowHubToRemoteVnetTransit         *bool                     `json:"allowHubToRemoteVnetTransit,omitempty"`
	AllowRemoteVnetToUseHubVnetGateways *bool                     `json:"allowRemoteVnetToUseHubVnetGateways,omitempty"`
	EnableInternetSecurity              *bool                     `json:"enableInternetSecurity,omitempty"`
	RoutingConfiguration                *RoutingConfiguration     `json:"routingConfiguration,omitempty"`
	ProvisioningState                   network.ProvisioningState `json:"provisioningState,omitempty"`
}

func (hvncp HubVirtualNetworkConnectionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if hvncp.RemoteVirtualNetwork != nil {
		objectMap["remoteVirtualNetwork"] = hvncp.RemoteVirtualNetwork
	}
	if hvncp.AllowHubToRemoteVnetTransit != nil {
		objectMap["allowHubToRemoteVnetTransit"] = hvncp.AllowHubToRemoteVnetTransit
	}
	if hvncp.AllowRemoteVnetToUseHubVnetGateways != nil {
		objectMap["allowRemoteVnetToUseHubVnetGateways"] = hvncp.AllowRemoteVnetToUseHubVnetGateways
	}
	if hvncp.EnableInternetSecurity != nil {
		objectMap["enableInternetSecurity"] = hvncp.EnableInternetSecurity
	}
	if hvncp.RoutingConfiguration != nil {
		objectMap["routingConfiguration"] = hvncp.RoutingConfiguration
	}
	return json.Marshal(objectMap)
}

type HubVirtualNetworkConnection struct {
	autorest.Response                      `json:"-"`
	*HubVirtualNetworkConnectionProperties `json:"properties,omitempty"`
	Name                                   *string `json:"name,omitempty"`
	Etag                                   *string `json:"etag,omitempty"`
	ID                                     *string `json:"id,omitempty"`
}

func (hvnc HubVirtualNetworkConnection) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if hvnc.HubVirtualNetworkConnectionProperties != nil {
		objectMap["properties"] = hvnc.HubVirtualNetworkConnectionProperties
	}
	if hvnc.Name != nil {
		objectMap["name"] = hvnc.Name
	}
	if hvnc.ID != nil {
		objectMap["id"] = hvnc.ID
	}
	return json.Marshal(objectMap)
}
