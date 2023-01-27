package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type MedTechServiceId struct {
	SubscriptionId   string
	ResourceGroup    string
	WorkspaceName    string
	IotConnectorName string
}

func NewMedTechServiceID(subscriptionId, resourceGroup, workspaceName, iotConnectorName string) MedTechServiceId {
	return MedTechServiceId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		WorkspaceName:    workspaceName,
		IotConnectorName: iotConnectorName,
	}
}

func (id MedTechServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Iot Connector Name %q", id.IotConnectorName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Med Tech Service", segmentsStr)
}

func (id MedTechServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/iotConnectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName)
}

// MedTechServiceID parses a MedTechService ID into an MedTechServiceId struct
func MedTechServiceID(input string) (*MedTechServiceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MedTechServiceId{
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
	if resourceId.IotConnectorName, err = id.PopSegment("iotConnectors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// MedTechServiceIDInsensitively parses an MedTechService ID into an MedTechServiceId struct, insensitively
// This should only be used to parse an ID for rewriting, the MedTechServiceID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func MedTechServiceIDInsensitively(input string) (*MedTechServiceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MedTechServiceId{
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

	// find the correct casing for the 'iotConnectors' segment
	iotConnectorsKey := "iotConnectors"
	for key := range id.Path {
		if strings.EqualFold(key, iotConnectorsKey) {
			iotConnectorsKey = key
			break
		}
	}
	if resourceId.IotConnectorName, err = id.PopSegment(iotConnectorsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
