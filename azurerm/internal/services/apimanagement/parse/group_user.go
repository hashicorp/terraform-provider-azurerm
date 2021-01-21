package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type GroupUserId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	GroupName      string
	UserName       string
}

func NewGroupUserID(subscriptionId, resourceGroup, serviceName, groupName, userName string) GroupUserId {
	return GroupUserId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		GroupName:      groupName,
		UserName:       userName,
	}
}

func (id GroupUserId) String() string {
	segments := []string{
		fmt.Sprintf("User Name %q", id.UserName),
		fmt.Sprintf("Group Name %q", id.GroupName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Group User", segmentsStr)
}

func (id GroupUserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/groups/%s/users/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.GroupName, id.UserName)
}

// GroupUserID parses a GroupUser ID into an GroupUserId struct
func GroupUserID(input string) (*GroupUserId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := GroupUserId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.GroupName, err = id.PopSegment("groups"); err != nil {
		return nil, err
	}
	if resourceId.UserName, err = id.PopSegment("users"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
