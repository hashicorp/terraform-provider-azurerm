package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DomainServiceTrustId struct {
	SubscriptionId    string
	ResourceGroup     string
	DomainServiceName string
	TrustName         string
}

func NewDomainServiceTrustID(subscriptionId, resourceGroup, domainServiceName, trustName string) DomainServiceTrustId {
	return DomainServiceTrustId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		DomainServiceName: domainServiceName,
		TrustName:         trustName,
	}
}

func (id DomainServiceTrustId) String() string {
	segments := []string{
		fmt.Sprintf("Trust Name %q", id.TrustName),
		fmt.Sprintf("Domain Service Name %q", id.DomainServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Domain Service Trust", segmentsStr)
}

func (id DomainServiceTrustId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AAD/domainServices/%s/trusts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, id.TrustName)
}

// DomainServiceTrustID parses a DomainServiceTrust ID into an DomainServiceTrustId struct
func DomainServiceTrustID(input string) (*DomainServiceTrustId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DomainServiceTrustId{
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
	if resourceId.TrustName, err = id.PopSegment("trusts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
