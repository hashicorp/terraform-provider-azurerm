package parse

import "fmt"

type FrontendEndpointId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewFrontendEndpointID(id FrontDoorId, name string) FrontendEndpointId {
	return FrontendEndpointId{
		ResourceGroup: id.ResourceGroup,
		FrontDoorName: id.Name,
		Name:          name,
	}
}

func (id FrontendEndpointId) ID(subscriptionId string) string {
	base := NewFrontDoorID(id.ResourceGroup, id.Name).ID(subscriptionId)
	return fmt.Sprintf("%s/frontendEndpoints/%s", base, id.Name)
}

func FrontDoorFrontendEndpointID(input string) (*FrontendEndpointId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Frontend Endpoint ID %q: %+v", input, err)
	}

	endpointId := FrontendEndpointId{
		ResourceGroup: frontDoorId.ResourceGroup,
		FrontDoorName: frontDoorId.Name,
	}

	// TODO: handle this being case-insensitive
	// https://github.com/Azure/azure-sdk-for-go/issues/6762
	if endpointId.Name, err = id.PopSegment("frontendEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &endpointId, nil
}
