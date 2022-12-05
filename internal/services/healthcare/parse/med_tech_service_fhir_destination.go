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
	IotConnectorName    string
	FhirDestinationName string
}

func NewMedTechServiceFhirDestinationID(subscriptionId, resourceGroup, workspaceName, iotConnectorName, fhirDestinationName string) MedTechServiceFhirDestinationId {
	return MedTechServiceFhirDestinationId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		WorkspaceName:       workspaceName,
		IotConnectorName:    iotConnectorName,
		FhirDestinationName: fhirDestinationName,
	}
}

func (id MedTechServiceFhirDestinationId) String() string {
	segments := []string{
		fmt.Sprintf("Fhir Destination Name %q", id.FhirDestinationName),
		fmt.Sprintf("Iot Connector Name %q", id.IotConnectorName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Med Tech Service Fhir Destination", segmentsStr)
}

func (id MedTechServiceFhirDestinationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthcareApis/workspaces/%s/iotConnectors/%s/fhirDestinations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName)
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
	if resourceId.IotConnectorName, err = id.PopSegment("iotConnectors"); err != nil {
		return nil, err
	}
	if resourceId.FhirDestinationName, err = id.PopSegment("fhirDestinations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// MedTechServiceFhirDestinationIDInsensitively parses an MedTechServiceFhirDestination ID into an MedTechServiceFhirDestinationId struct, insensitively
// This should only be used to parse an ID for rewriting, the MedTechServiceFhirDestinationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func MedTechServiceFhirDestinationIDInsensitively(input string) (*MedTechServiceFhirDestinationId, error) {
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

	// find the correct casing for the 'fhirDestinations' segment
	fhirDestinationsKey := "fhirDestinations"
	for key := range id.Path {
		if strings.EqualFold(key, fhirDestinationsKey) {
			fhirDestinationsKey = key
			break
		}
	}
	if resourceId.FhirDestinationName, err = id.PopSegment(fhirDestinationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
