package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DomainServiceResourceId struct {
	SubscriptionId    string
	ResourceGroup     string
	DomainServiceName string
}

func NewDomainServiceResourceID(subscriptionId, resourceGroup, domainServiceName string) DomainServiceResourceId {
	return DomainServiceResourceId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		DomainServiceName: domainServiceName,
	}
}

func (id DomainServiceResourceId) String() string {
	segments := []string{
		fmt.Sprintf("Domain Service Name %q", id.DomainServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Domain Service Resource", segmentsStr)
}

func (id DomainServiceResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AAD/domainServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)
}

// DomainServiceResourceID parses a DomainServiceResource ID into an DomainServiceResourceId struct
func DomainServiceResourceID(input string) (*DomainServiceResourceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DomainServiceResourceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DomainServiceName, err = id.PopSegment("domainServices"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
