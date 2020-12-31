package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FailoverGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerName     string
	Name           string
}

func NewFailoverGroupID(subscriptionId, resourceGroup, serverName, name string) FailoverGroupId {
	return FailoverGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerName:     serverName,
		Name:           name,
	}
}

func (id FailoverGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Failover Group", segmentsStr)
}

func (id FailoverGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/failoverGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.Name)
}

// FailoverGroupID parses a FailoverGroup ID into an FailoverGroupId struct
func FailoverGroupID(input string) (*FailoverGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FailoverGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("failoverGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
