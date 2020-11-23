package parse

import "fmt"

type LoadBalancingId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewLoadBalancingID(id FrontDoorId, name string) LoadBalancingId {
	return LoadBalancingId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		FrontDoorName:  id.Name,
		Name:           name,
	}
}

func (id LoadBalancingId) ID(_ string) string {
	base := NewFrontDoorID(id.SubscriptionId, id.ResourceGroup, id.FrontDoorName).ID("")
	return fmt.Sprintf("%s/loadBalancingSettings/%s", base, id.Name)
}

func LoadBalancingID(input string) (*LoadBalancingId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor Load Balancing ID %q: %+v", input, err)
	}

	loadBalancingId := LoadBalancingId{
		SubscriptionId: frontDoorId.SubscriptionId,
		ResourceGroup:  frontDoorId.ResourceGroup,
		FrontDoorName:  frontDoorId.Name,
	}

	// https://github.com/Azure/azure-sdk-for-go/issues/6762
	// note: the ordering is important since the defined case (we want to error with) is loadBalancingSettings
	if loadBalancingId.Name, err = id.PopSegment("LoadBalancingSettings"); err != nil {
		if loadBalancingId.Name, err = id.PopSegment("loadBalancingSettings"); err != nil {
			return nil, err
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &loadBalancingId, nil
}
