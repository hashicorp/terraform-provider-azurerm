package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AppServiceIPRestrictionId struct {
	SubscriptionId    string
	ResourceGroup     string
	SiteName          string
	IpRestrictionName string
}

func NewAppServiceIPRestrictionID(subscriptionId, resourceGroup, siteName, ipRestrictionName string) AppServiceIPRestrictionId {
	return AppServiceIPRestrictionId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		SiteName:          siteName,
		IpRestrictionName: ipRestrictionName,
	}
}

func (id AppServiceIPRestrictionId) String() string {
	segments := []string{
		fmt.Sprintf("Ip Restriction Name %q", id.IpRestrictionName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Service I P Restriction", segmentsStr)
}

func (id AppServiceIPRestrictionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/ipRestriction/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.IpRestrictionName)
}

// AppServiceIPRestrictionID parses a AppServiceIPRestriction ID into an AppServiceIPRestrictionId struct
func AppServiceIPRestrictionID(input string) (*AppServiceIPRestrictionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AppServiceIPRestrictionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.IpRestrictionName, err = id.PopSegment("ipRestriction"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
