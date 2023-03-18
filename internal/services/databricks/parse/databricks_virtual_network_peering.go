package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DatabricksVirtualNetworkPeeringId struct {
	SubscriptionId            string
	ResourceGroup             string
	WorkspaceName             string
	VirtualNetworkPeeringName string
}

func NewDatabricksVirtualNetworkPeeringID(subscriptionId, resourceGroup, workspaceName, virtualNetworkPeeringName string) DatabricksVirtualNetworkPeeringId {
	return DatabricksVirtualNetworkPeeringId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		WorkspaceName:             workspaceName,
		VirtualNetworkPeeringName: virtualNetworkPeeringName,
	}
}

func (id DatabricksVirtualNetworkPeeringId) String() string {
	segments := []string{
		fmt.Sprintf("Virtual Network Peering Name %q", id.VirtualNetworkPeeringName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Databricks Virtual Network Peering", segmentsStr)
}

func (id DatabricksVirtualNetworkPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Databricks/workspaces/%s/virtualNetworkPeerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.VirtualNetworkPeeringName)
}

// DatabricksVirtualNetworkPeeringID parses a DatabricksVirtualNetworkPeering ID into an DatabricksVirtualNetworkPeeringId struct
func DatabricksVirtualNetworkPeeringID(input string) (*DatabricksVirtualNetworkPeeringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatabricksVirtualNetworkPeeringId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.VirtualNetworkPeeringName, err = id.PopSegment("virtualNetworkPeerings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// DatabricksVirtualNetworkPeeringIDInsensitively parses an DatabricksVirtualNetworkPeering ID into an DatabricksVirtualNetworkPeeringId struct, insensitively
// This should only be used to parse an ID for rewriting, the DatabricksVirtualNetworkPeeringID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func DatabricksVirtualNetworkPeeringIDInsensitively(input string) (*DatabricksVirtualNetworkPeeringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatabricksVirtualNetworkPeeringId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'workspaces' segment
	workspacesKey := "workspaces"
	for key := range id.Path {
		if strings.EqualFold(key, workspacesKey) {
			workspacesKey = key
			break
		}
	}
	if resourceId.WorkspaceName, err = id.PopSegment(workspacesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'virtualNetworkPeerings' segment
	virtualNetworkPeeringsKey := "virtualNetworkPeerings"
	for key := range id.Path {
		if strings.EqualFold(key, virtualNetworkPeeringsKey) {
			virtualNetworkPeeringsKey = key
			break
		}
	}
	if resourceId.VirtualNetworkPeeringName, err = id.PopSegment(virtualNetworkPeeringsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
