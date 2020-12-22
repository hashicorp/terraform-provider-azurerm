package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudCustomDomainId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	AppName        string
	DomainName     string
}

func NewSpringCloudCustomDomainID(subscriptionId, resourceGroup, springName, appName, domainName string) SpringCloudCustomDomainId {
	return SpringCloudCustomDomainId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		AppName:        appName,
		DomainName:     domainName,
	}
}

func (id SpringCloudCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Domain Name %q", id.DomainName),
		fmt.Sprintf("App Name %q", id.AppName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Custom Domain", segmentsStr)
}

func (id SpringCloudCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/apps/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName, id.DomainName)
}

// SpringCloudCustomDomainID parses a SpringCloudCustomDomain ID into an SpringCloudCustomDomainId struct
func SpringCloudCustomDomainID(input string) (*SpringCloudCustomDomainId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudCustomDomainId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}
	if resourceId.AppName, err = id.PopSegment("apps"); err != nil {
		return nil, err
	}
	if resourceId.DomainName, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
