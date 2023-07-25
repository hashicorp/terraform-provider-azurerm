// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
)

// NOTE: this workaround client exists to allow the removal of 'healthProbes' during update

type CdnFrontDoorOriginGroupsWorkaroundClient struct {
	sdkClient *cdn.AFDOriginGroupsClient
}

func NewCdnFrontDoorOriginGroupsWorkaroundClient(client *cdn.AFDOriginGroupsClient) CdnFrontDoorOriginGroupsWorkaroundClient {
	return CdnFrontDoorOriginGroupsWorkaroundClient{
		sdkClient: client,
	}
}

func (c *CdnFrontDoorOriginGroupsWorkaroundClient) Update(ctx context.Context, resourceGroupName string, profileName string, originGroupName string, afdOriginGroupUpdateProperties AFDOriginGroupUpdateParameters) (result cdn.AFDOriginGroupsUpdateFuture, err error) {
	req, err := c.UpdatePreparer(ctx, resourceGroupName, profileName, originGroupName, afdOriginGroupUpdateProperties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.OriginGroupsClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = c.sdkClient.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.OriginGroupsClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client CdnFrontDoorOriginGroupsWorkaroundClient) UpdatePreparer(ctx context.Context, resourceGroupName string, profileName string, originGroupName string, afdOriginGroupUpdateProperties AFDOriginGroupUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"originGroupName":   autorest.Encode("path", originGroupName),
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.sdkClient.SubscriptionID),
	}

	const APIVersion = "2021-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.sdkClient.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/originGroups/{originGroupName}", pathParameters),
		autorest.WithJSON(afdOriginGroupUpdateProperties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// OriginGroupUpdateParameters origin group properties needed for origin group creation or update.
type AFDOriginGroupUpdateParameters struct {
	*AFDOriginGroupUpdatePropertiesParameters `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for OriginGroupUpdateParameters.
func (ogup AFDOriginGroupUpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ogup.AFDOriginGroupUpdatePropertiesParameters != nil {
		objectMap["properties"] = ogup.AFDOriginGroupUpdatePropertiesParameters
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for AFDOriginGroupUpdateParameters struct.
func (ogup *AFDOriginGroupUpdateParameters) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k { // nolint gocritic
		case "properties":
			if v != nil {
				var afdOriginGroupUpdatePropertiesParameters AFDOriginGroupUpdatePropertiesParameters
				err = json.Unmarshal(*v, &afdOriginGroupUpdatePropertiesParameters)
				if err != nil {
					return err
				}
				ogup.AFDOriginGroupUpdatePropertiesParameters = &afdOriginGroupUpdatePropertiesParameters
			}
		}
	}

	return nil
}

// AFDOriginGroupUpdatePropertiesParameters the JSON object that contains the properties of the origin group.
type AFDOriginGroupUpdatePropertiesParameters struct {
	// ProfileName - READ-ONLY; The name of the profile which holds the origin group.
	ProfileName *string `json:"profileName,omitempty"`
	// LoadBalancingSettings - Load balancing settings for a backend pool
	LoadBalancingSettings *cdn.LoadBalancingSettingsParameters `json:"loadBalancingSettings,omitempty"`
	// HealthProbeSettings - Health probe settings to the origin that is used to determine the health of the origin.
	HealthProbeSettings *cdn.HealthProbeParameters `json:"healthProbeSettings,omitempty"`
	// TrafficRestorationTimeToHealedOrNewEndpointsInMinutes - Time in minutes to shift the traffic to the endpoint gradually when an unhealthy endpoint comes healthy or a new endpoint is added. Default is 10 mins. This property is currently not supported.
	TrafficRestorationTimeToHealedOrNewEndpointsInMinutes *int32 `json:"trafficRestorationTimeToHealedOrNewEndpointsInMinutes,omitempty"`
	// SessionAffinityState - Whether to allow session affinity on this host. Valid options are 'Enabled' or 'Disabled'. Possible values include: 'EnabledStateEnabled', 'EnabledStateDisabled'
	SessionAffinityState cdn.EnabledState `json:"sessionAffinityState,omitempty"`
}

// MarshalJSON is the custom marshaler for AFDOriginGroupUpdatePropertiesParameters.
func (ogupp AFDOriginGroupUpdatePropertiesParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if ogupp.LoadBalancingSettings != nil {
		objectMap["loadBalancingSettings"] = ogupp.LoadBalancingSettings
	}

	objectMap["healthProbeSettings"] = ogupp.HealthProbeSettings

	if ogupp.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes != nil {
		objectMap["trafficRestorationTimeToHealedOrNewEndpointsInMinutes"] = ogupp.TrafficRestorationTimeToHealedOrNewEndpointsInMinutes
	}
	if ogupp.SessionAffinityState != "" {
		objectMap["sessionAffinityState"] = ogupp.SessionAffinityState
	}
	return json.Marshal(objectMap)
}
