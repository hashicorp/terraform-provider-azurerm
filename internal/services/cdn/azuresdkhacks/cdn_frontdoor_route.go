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

// NOTE: this workaround client exists to allow the removal of 'cacheConfiguration' items during update
type CdnFrontDoorRoutesWorkaroundClient struct {
	sdkClient *cdn.RoutesClient
}

func NewCdnFrontDoorRoutesWorkaroundClient(client *cdn.RoutesClient) CdnFrontDoorRoutesWorkaroundClient {
	return CdnFrontDoorRoutesWorkaroundClient{
		sdkClient: client,
	}
}

func (c *CdnFrontDoorRoutesWorkaroundClient) Update(ctx context.Context, resourceGroupName string, profileName string, endpointName string, routeName string, routeUpdateProperties RouteUpdateParameters) (result cdn.RoutesUpdateFuture, err error) {
	req, err := c.UpdatePreparer(ctx, resourceGroupName, profileName, endpointName, routeName, routeUpdateProperties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RoutesClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = c.sdkClient.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cdn.RoutesClient", "Update", result.Response(), "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client CdnFrontDoorRoutesWorkaroundClient) UpdatePreparer(ctx context.Context, resourceGroupName string, profileName string, endpointName string, routeName string, routeUpdateProperties RouteUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"endpointName":      autorest.Encode("path", endpointName),
		"profileName":       autorest.Encode("path", profileName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"routeName":         autorest.Encode("path", routeName),
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
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/afdEndpoints/{endpointName}/routes/{routeName}", pathParameters),
		autorest.WithJSON(routeUpdateProperties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RouteUpdateParameters the domain JSON object required for domain creation or update.
type RouteUpdateParameters struct {
	*RouteUpdatePropertiesParameters `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for RouteUpdateParameters.
func (rup RouteUpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if rup.RouteUpdatePropertiesParameters != nil {
		objectMap["properties"] = rup.RouteUpdatePropertiesParameters
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for RouteUpdateParameters struct.
func (rup *RouteUpdateParameters) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k { // nolint gocritic
		case "properties":
			if v != nil {
				var routeUpdatePropertiesParameters RouteUpdatePropertiesParameters
				err = json.Unmarshal(*v, &routeUpdatePropertiesParameters)
				if err != nil {
					return err
				}

				rup.RouteUpdatePropertiesParameters = &routeUpdatePropertiesParameters
			}
		}
	}

	return nil
}

// RouteUpdatePropertiesParameters the JSON object that contains the properties of the domain to create.
type RouteUpdatePropertiesParameters struct {
	// EndpointName - READ-ONLY; The name of the endpoint which holds the route.
	EndpointName *string `json:"endpointName,omitempty"`
	// CustomDomains - Domains referenced by this endpoint.
	CustomDomains *[]cdn.ActivatedResourceReference `json:"customDomains,omitempty"`
	// OriginGroup - A reference to the origin group.
	OriginGroup *cdn.ResourceReference `json:"originGroup,omitempty"`
	// OriginPath - A directory path on the origin that AzureFrontDoor can use to retrieve content from, e.g. contoso.cloudapp.net/originpath.
	OriginPath *string `json:"originPath,omitempty"`
	// RuleSets - rule sets referenced by this endpoint.
	RuleSets *[]cdn.ResourceReference `json:"ruleSets,omitempty"`
	// SupportedProtocols - List of supported protocols for this route.
	SupportedProtocols *[]cdn.AFDEndpointProtocols `json:"supportedProtocols,omitempty"`
	// PatternsToMatch - The route patterns of the rule.
	PatternsToMatch *[]string `json:"patternsToMatch,omitempty"`
	// CacheConfiguration - The caching configuration for this route. To disable caching, do not provide a cacheConfiguration object.
	CacheConfiguration *cdn.AfdRouteCacheConfiguration `json:"cacheConfiguration,omitempty"`
	// ForwardingProtocol - Protocol this rule will use when forwarding traffic to backends. Possible values include: 'ForwardingProtocolHTTPOnly', 'ForwardingProtocolHTTPSOnly', 'ForwardingProtocolMatchRequest'
	ForwardingProtocol cdn.ForwardingProtocol `json:"forwardingProtocol,omitempty"`
	// LinkToDefaultDomain - whether this route will be linked to the default endpoint domain. Possible values include: 'LinkToDefaultDomainEnabled', 'LinkToDefaultDomainDisabled'
	LinkToDefaultDomain cdn.LinkToDefaultDomain `json:"linkToDefaultDomain,omitempty"`
	// HTTPSRedirect - Whether to automatically redirect HTTP traffic to HTTPS traffic. Note that this is a easy way to set up this rule and it will be the first rule that gets executed. Possible values include: 'HTTPSRedirectEnabled', 'HTTPSRedirectDisabled'
	HTTPSRedirect cdn.HTTPSRedirect `json:"httpsRedirect,omitempty"`
	// EnabledState - Whether to enable use of this rule. Permitted values are 'Enabled' or 'Disabled'. Possible values include: 'EnabledStateEnabled', 'EnabledStateDisabled'
	EnabledState cdn.EnabledState `json:"enabledState,omitempty"`
}

// MarshalJSON is the custom marshaler for RouteUpdatePropertiesParameters.
func (rupp RouteUpdatePropertiesParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})

	objectMap["customDomains"] = rupp.CustomDomains

	if rupp.OriginGroup != nil {
		objectMap["originGroup"] = rupp.OriginGroup
	}

	// OriginPath must be set to nil to be removed
	objectMap["originPath"] = rupp.OriginPath

	if rupp.RuleSets != nil {
		objectMap["ruleSets"] = rupp.RuleSets
	}
	if rupp.SupportedProtocols != nil {
		objectMap["supportedProtocols"] = rupp.SupportedProtocols
	}
	if rupp.PatternsToMatch != nil {
		objectMap["patternsToMatch"] = rupp.PatternsToMatch
	}

	objectMap["cacheConfiguration"] = rupp.CacheConfiguration

	if rupp.ForwardingProtocol != "" {
		objectMap["forwardingProtocol"] = rupp.ForwardingProtocol
	}
	if rupp.LinkToDefaultDomain != "" {
		objectMap["linkToDefaultDomain"] = rupp.LinkToDefaultDomain
	}
	if rupp.HTTPSRedirect != "" {
		objectMap["httpsRedirect"] = rupp.HTTPSRedirect
	}
	if rupp.EnabledState != "" {
		objectMap["enabledState"] = rupp.EnabledState
	}
	return json.Marshal(objectMap)
}
