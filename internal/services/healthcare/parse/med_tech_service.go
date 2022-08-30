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
	IotconnectorName string
}

func NewMedTechServiceID(subscriptionId, resourceGroup, workspaceName, iotconnectorName string) MedTechServiceId {
	return MedTechServiceId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		WorkspaceName:    workspaceName,
		IotconnectorName: iotconnectorName,
	}
}

func (id MedTechServiceId) String() string {
	segments := []string{
		fmt.Sprintf("Iotconnector Name %q", id.IotconnectorName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Med Tech Service", segmentsStr)
}

func (id MedTechServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/iotconnectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.IotconnectorName)
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
	if resourceId.IotconnectorName, err = id.PopSegment("iotconnectors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
