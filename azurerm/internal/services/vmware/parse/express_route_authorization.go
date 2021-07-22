package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteAuthorizationId struct {
	SubscriptionId    string
	ResourceGroup     string
	PrivateCloudName  string
	AuthorizationName string
}

func NewExpressRouteAuthorizationID(subscriptionId, resourceGroup, privateCloudName, authorizationName string) ExpressRouteAuthorizationId {
	return ExpressRouteAuthorizationId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		PrivateCloudName:  privateCloudName,
		AuthorizationName: authorizationName,
	}
}

func (id ExpressRouteAuthorizationId) String() string {
	segments := []string{
		fmt.Sprintf("Authorization Name %q", id.AuthorizationName),
		fmt.Sprintf("Private Cloud Name %q", id.PrivateCloudName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Authorization", segmentsStr)
}

func (id ExpressRouteAuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateCloudName, id.AuthorizationName)
}

// ExpressRouteAuthorizationID parses a ExpressRouteAuthorization ID into an ExpressRouteAuthorizationId struct
func ExpressRouteAuthorizationID(input string) (*ExpressRouteAuthorizationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExpressRouteAuthorizationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateCloudName, err = id.PopSegment("privateClouds"); err != nil {
		return nil, err
	}
	if resourceId.AuthorizationName, err = id.PopSegment("authorizations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
