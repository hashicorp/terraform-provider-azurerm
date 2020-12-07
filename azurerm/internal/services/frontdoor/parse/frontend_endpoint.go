package parse

import "fmt"

type FrontendEndpointId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewFrontendEndpointID(subscriptionId, resourceGroup, frontDoorName, name string) FrontendEndpointId {
	return FrontendEndpointId{
		ResourceGroup: resourceGroup,
		FrontDoorName: frontDoorName,
		Name:          name,
	}
}

func (id FrontendEndpointId) ID(subscriptionId string) string {
	base := NewFrontDoorID(subscriptionId, id.ResourceGroup, id.FrontDoorName).ID(subscriptionId)
	return fmt.Sprintf("%s/frontendEndpoints/%s", base, id.Name)
}

func FrontendEndpointIDInsensitively(input string) (*FrontendEndpointId, error) {
	return parseFrontendEndpointID(input, false)
}

func FrontendEndpointID(input string) (*FrontendEndpointId, error) {
	return parseFrontendEndpointID(input, true)
}

func parseFrontendEndpointID(input string, caseSensitive bool) (*FrontendEndpointId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Frontend Endpoint ID %q: %+v", input, err)
	}

	endpointId := FrontendEndpointId{
		ResourceGroup: frontDoorId.ResourceGroup,
		FrontDoorName: frontDoorId.Name,
	}

	// The Azure API (per the ARM Spec/chatting with the ARM Team) should be following Postel's Law;
	// where ID's are insensitive for Requests but sensitive in responses - but it's not.
	//
	// For us this means ID's should be sensitive at import time - but we have to work around these
	// API bugs for the moment.
	if caseSensitive {
		if endpointId.Name, err = id.PopSegment("frontendEndpoints"); err != nil {
			return nil, err
		}
	} else {
		// https://github.com/Azure/azure-sdk-for-go/issues/6762
		// note: the ordering is important since the defined case (we want to error with) is frontendEndpoints
		if endpointId.Name, err = id.PopSegment("FrontendEndpoints"); err != nil {
			if endpointId.Name, err = id.PopSegment("frontendEndpoints"); err != nil {
				return nil, err
			}
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &endpointId, nil
}
