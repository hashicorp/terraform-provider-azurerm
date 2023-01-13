package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type NetworkManagerManagementGroupConnectionId struct {
	ManagementGroupName          string
	NetworkManagerConnectionName string
}

func NewNetworkManagerManagementGroupConnectionID(managementGroupName, networkManagerConnectionName string) NetworkManagerManagementGroupConnectionId {
	return NetworkManagerManagementGroupConnectionId{
		ManagementGroupName:          managementGroupName,
		NetworkManagerConnectionName: networkManagerConnectionName,
	}
}

func (id NetworkManagerManagementGroupConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Network Manager Connection Name %q", id.NetworkManagerConnectionName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Management Group Connection", segmentsStr)
}

func (id NetworkManagerManagementGroupConnectionId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Network/networkManagerConnections/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.NetworkManagerConnectionName)
}

// NetworkManagerManagementGroupConnectionID parses a NetworkManagerManagementGroupConnection ID into an NetworkManagerManagementGroupConnectionId struct
func NetworkManagerManagementGroupConnectionID(input string) (*NetworkManagerManagementGroupConnectionId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerManagementGroupConnectionId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
		return nil, err
	}
	if resourceId.NetworkManagerConnectionName, err = id.PopSegment("networkManagerConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
