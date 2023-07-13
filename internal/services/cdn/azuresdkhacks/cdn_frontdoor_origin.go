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

// NOTE: this workaround client exists to allow the removal of 'originHostHeader' during update

type CdnFrontDoorOriginWorkaroundClient struct {
	sdkClient *cdn.AFDOriginsClient
}

func NewCdnFrontDoorOriginsWorkaroundClient(client *cdn.AFDOriginsClient) CdnFrontDoorOriginWorkaroundClient {
	return CdnFrontDoorOriginWorkaroundClient{
		sdkClient: client,
	}
}

func (c *CdnFrontDoorOriginWorkaroundClient) Update(ctx context.Context, resourceGroupName string, profileName string, originGroupName string, originName string, originUpdateProperties AFDOriginUpdateParameters) (result cdn.AFDOriginsUpdateFuture, err error) {
	req, err := c.UpdatePreparer(ctx, resourceGroupName, profileName, originGroupName, originName, originUpdateProperties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = c.sdkClient.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.OriginsClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client CdnFrontDoorOriginWorkaroundClient) UpdatePreparer(ctx context.Context, resourceGroupName string, profileName string, originGroupName string, originName string, originUpdateProperties AFDOriginUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"originGroupName":   autorest.Encode("path", originGroupName),
		"originName":        autorest.Encode("path", originName),
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/originGroups/{originGroupName}/origins/{originName}", pathParameters),
		autorest.WithJSON(originUpdateProperties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// AFDOriginUpdateParameters afdOrigin properties needed for origin update.
type AFDOriginUpdateParameters struct {
	*AFDOriginUpdatePropertiesParameters `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for AFDOriginUpdateParameters.
func (aoup AFDOriginUpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if aoup.AFDOriginUpdatePropertiesParameters != nil {
		objectMap["properties"] = aoup.AFDOriginUpdatePropertiesParameters
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for AFDOriginUpdateParameters struct.
func (aoup *AFDOriginUpdateParameters) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k { // nolint gocritic
		case "properties":
			if v != nil {
				var aFDOriginUpdatePropertiesParameters AFDOriginUpdatePropertiesParameters
				err = json.Unmarshal(*v, &aFDOriginUpdatePropertiesParameters)
				if err != nil {
					return err
				}
				aoup.AFDOriginUpdatePropertiesParameters = &aFDOriginUpdatePropertiesParameters
			}
		}
	}

	return nil
}

// AFDOriginUpdatePropertiesParameters the JSON object that contains the properties of the origin.
type AFDOriginUpdatePropertiesParameters struct {
	// OriginGroupName - READ-ONLY; The name of the origin group which contains this origin.
	OriginGroupName *string `json:"originGroupName,omitempty"`
	// AzureOrigin - Resource reference to the Azure origin resource.
	AzureOrigin *cdn.ResourceReference `json:"azureOrigin,omitempty"`
	// HostName - The address of the origin. Domain names, IPv4 addresses, and IPv6 addresses are supported.This should be unique across all origins in an endpoint.
	HostName *string `json:"hostName,omitempty"`
	// HTTPPort - The value of the HTTP port. Must be between 1 and 65535.
	HTTPPort *int32 `json:"httpPort,omitempty"`
	// HTTPSPort - The value of the HTTPS port. Must be between 1 and 65535.
	HTTPSPort *int32 `json:"httpsPort,omitempty"`
	// OriginHostHeader - The host header value sent to the origin with each request. If you leave this blank, the request hostname determines this value. Azure CDN origins, such as Web Apps, Blob Storage, and Cloud Services require this host header value to match the origin hostname by default. This overrides the host header defined at Endpoint
	OriginHostHeader *string `json:"originHostHeader,omitempty"`
	// Priority - Priority of origin in given origin group for load balancing. Higher priorities will not be used for load balancing if any lower priority origin is healthy.Must be between 1 and 5
	Priority *int32 `json:"priority,omitempty"`
	// Weight - Weight of the origin in given origin group for load balancing. Must be between 1 and 1000
	Weight *int32 `json:"weight,omitempty"`
	// SharedPrivateLinkResource - The properties of the private link resource for private origin.
	SharedPrivateLinkResource *cdn.SharedPrivateLinkResourceProperties `json:"sharedPrivateLinkResource,omitempty"`
	// EnabledState - Whether to enable health probes to be made against backends defined under backendPools. Health probes can only be disabled if there is a single enabled backend in single enabled backend pool. Possible values include: 'EnabledStateEnabled', 'EnabledStateDisabled'
	EnabledState cdn.EnabledState `json:"enabledState,omitempty"`
	// EnforceCertificateNameCheck - Whether to enable certificate name check at origin level
	EnforceCertificateNameCheck *bool `json:"enforceCertificateNameCheck,omitempty"`
}

// MarshalJSON is the custom marshaler for AFDOriginUpdatePropertiesParameters.
func (aoupp AFDOriginUpdatePropertiesParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if aoupp.AzureOrigin != nil {
		objectMap["azureOrigin"] = aoupp.AzureOrigin
	}
	if aoupp.HostName != nil {
		objectMap["hostName"] = aoupp.HostName
	}
	if aoupp.HTTPPort != nil {
		objectMap["httpPort"] = aoupp.HTTPPort
	}
	if aoupp.HTTPSPort != nil {
		objectMap["httpsPort"] = aoupp.HTTPSPort
	}

	objectMap["originHostHeader"] = aoupp.OriginHostHeader

	if aoupp.Priority != nil {
		objectMap["priority"] = aoupp.Priority
	}
	if aoupp.Weight != nil {
		objectMap["weight"] = aoupp.Weight
	}
	if aoupp.SharedPrivateLinkResource != nil {
		objectMap["sharedPrivateLinkResource"] = aoupp.SharedPrivateLinkResource
	}
	if aoupp.EnabledState != "" {
		objectMap["enabledState"] = aoupp.EnabledState
	}
	if aoupp.EnforceCertificateNameCheck != nil {
		objectMap["enforceCertificateNameCheck"] = aoupp.EnforceCertificateNameCheck
	}
	return json.Marshal(objectMap)
}
