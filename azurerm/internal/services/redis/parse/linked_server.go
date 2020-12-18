package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LinkedServerId struct {
	SubscriptionId string
	ResourceGroup  string
	RediName       string
	Name           string
}

func NewLinkedServerID(subscriptionId, resourceGroup, rediName, name string) LinkedServerId {
	return LinkedServerId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RediName:       rediName,
		Name:           name,
	}
}

func (id LinkedServerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Redi Name %q", id.RediName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Linked Server", segmentsStr)
}

func (id LinkedServerId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/Redis/%s/linkedServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RediName, id.Name)
}

// LinkedServerID parses a LinkedServer ID into an LinkedServerId struct
func LinkedServerID(input string) (*LinkedServerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LinkedServerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RediName, err = id.PopSegment("Redis"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("linkedServers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
