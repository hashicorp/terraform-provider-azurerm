package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type MedTechServiceFhirDestinationId struct {
	SubscriptionId      string
	ResourceGroup       string
	WorkspaceName       string
	IotconnectorName    string
	FhirdestinationName string
}

func NewMedTechServiceFhirDestinationID(subscriptionId, resourceGroup, workspaceName, iotconnectorName, fhirdestinationName string) MedTechServiceFhirDestinationId {
	return MedTechServiceFhirDestinationId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		WorkspaceName:       workspaceName,
		IotconnectorName:    iotconnectorName,
		FhirdestinationName: fhirdestinationName,
	}
}

func (id MedTechServiceFhirDestinationId) String() string {
	segments := []string{
		fmt.Sprintf("Fhirdestination Name %q", id.FhirdestinationName),
		fmt.Sprintf("Iotconnector Name %q", id.IotconnectorName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Med Tech Service Fhir Destination", segmentsStr)
}

func (id MedTechServiceFhirDestinationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/iotconnectors/%s/fhirdestinations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.IotconnectorName, id.FhirdestinationName)
}

// MedTechServiceFhirDestinationID parses a MedTechServiceFhirDestination ID into an MedTechServiceFhirDestinationId struct
func MedTechServiceFhirDestinationID(input string) (*MedTechServiceFhirDestinationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MedTechServiceFhirDestinationId{
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
	if resourceId.FhirdestinationName, err = id.PopSegment("fhirdestinations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
