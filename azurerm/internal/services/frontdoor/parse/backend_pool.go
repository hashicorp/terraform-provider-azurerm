package parse

import "fmt"

type BackendPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewBackendPoolID(id FrontDoorId, name string) BackendPoolId {
	return BackendPoolId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		FrontDoorName:  id.Name,
		Name:           name,
	}
}

func (id BackendPoolId) ID(_ string) string {
	base := NewFrontDoorID(id.SubscriptionId, id.ResourceGroup, id.FrontDoorName).ID("")
	return fmt.Sprintf("%s/backendPools/%s", base, id.Name)
}

func BackendPoolID(input string) (*BackendPoolId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Backend Pool ID %q: %+v", input, err)
	}

	poolId := BackendPoolId{
		SubscriptionId: frontDoorId.SubscriptionId,
		ResourceGroup:  frontDoorId.ResourceGroup,
		FrontDoorName:  frontDoorId.Name,
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
