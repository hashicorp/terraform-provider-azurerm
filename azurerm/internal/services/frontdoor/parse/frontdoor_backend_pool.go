package parse

import "fmt"

type FrontDoorBackendPoolId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewFrontDoorBackendPoolID(id FrontDoorId, name string) FrontDoorBackendPoolId {
	return FrontDoorBackendPoolId{
		ResourceGroup: id.ResourceGroup,
		FrontDoorName: id.Name,
		Name:          name,
	}
}

func (id FrontDoorBackendPoolId) ID(subscriptionId string) string {
	base := NewFrontDoorID(id.ResourceGroup, id.Name).ID(subscriptionId)
	return fmt.Sprintf("%s/backendPools/%s", base, id.Name)
}

func FrontDoorBackendPoolID(input string) (*FrontDoorBackendPoolId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor Backend Pool ID %q: %+v", input, err)
	}

	endpointId := FrontDoorBackendPoolId{
		ResourceGroup: frontDoorId.ResourceGroup,
		FrontDoorName: frontDoorId.Name,
	}

	// TODO: handle this being case-insensitive
	// https://github.com/Azure/azure-sdk-for-go/issues/6762
	if endpointId.Name, err = id.PopSegment("backendPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return nil, nil
}
