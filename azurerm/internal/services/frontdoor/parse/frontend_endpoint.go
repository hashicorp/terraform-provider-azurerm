package parse

import "fmt"

type FrontendEndpointId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewFrontendEndpointID(id FrontDoorId, name string) FrontendEndpointId {
	return FrontendEndpointId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		FrontDoorName:  id.Name,
		Name:           name,
	}
}

func (id FrontendEndpointId) ID(_ string) string {
	base := NewFrontDoorID(id.SubscriptionId, id.ResourceGroup, id.FrontDoorName).ID("")
	return fmt.Sprintf("%s/frontendEndpoints/%s", base, id.Name)
}

func FrontendEndpointID(input string) (*FrontendEndpointId, error) {
	return parseFrontendEndpointID(input, false)
}

func FrontendEndpointIDForImport(input string) (*FrontendEndpointId, error) {
	return parseFrontendEndpointID(input, true)
}

func parseFrontendEndpointID(input string, caseSensitive bool) (*FrontendEndpointId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Frontend Endpoint ID %q: %+v", input, err)
	}

	endpointId := FrontendEndpointId{
		SubscriptionId: frontDoorId.SubscriptionId,
		ResourceGroup:  frontDoorId.ResourceGroup,
		FrontDoorName:  frontDoorId.Name,
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
