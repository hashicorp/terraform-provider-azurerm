package parse

import "fmt"

type BackendPoolId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewBackendPoolID(subscriptionId, resourceGroup, frontDoorName, name string) BackendPoolId {
	return BackendPoolId{
		ResourceGroup: resourceGroup,
		FrontDoorName: frontDoorName,
		Name:          name,
	}
}

func (id BackendPoolId) ID(subscriptionId string) string {
	base := NewFrontDoorID(subscriptionId, id.ResourceGroup, id.FrontDoorName).ID(subscriptionId)
	return fmt.Sprintf("%s/backendPools/%s", base, id.Name)
}

func BackendPoolIDInsensitively(input string) (*BackendPoolId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Backend Pool ID %q: %+v", input, err)
	}

	poolId := BackendPoolId{
		ResourceGroup: frontDoorId.ResourceGroup,
		FrontDoorName: frontDoorId.Name,
	}

	// API is broken - https://github.com/Azure/azure-sdk-for-go/issues/6762
	// note: the ordering is important since the defined case (we want to error with) is backendPools
	if poolId.Name, err = id.PopSegment("backendpools"); err != nil {
		if poolId.Name, err = id.PopSegment("BackendPools"); err != nil {
			if poolId.Name, err = id.PopSegment("backendPools"); err != nil {
				return nil, err
			}
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &poolId, nil
}
