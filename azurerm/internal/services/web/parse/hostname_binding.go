package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HostnameBindingId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
	Name           string
}

func NewHostnameBindingID(subscriptionId, resourceGroup, siteName, name string) HostnameBindingId {
	return HostnameBindingId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
		Name:           name,
	}
}

func (id HostnameBindingId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id HostnameBindingId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hostNameBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.Name)
}

// HostnameBindingID parses a HostnameBinding ID into an HostnameBindingId struct
func HostnameBindingID(input string) (*HostnameBindingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HostnameBindingId{
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
	if resourceId.Name, err = id.PopSegment("hostNameBindings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
