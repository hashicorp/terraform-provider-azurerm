package parse

import "fmt"

type RoutingRuleId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewRoutingRuleID(id FrontDoorId, name string) RoutingRuleId {
	return RoutingRuleId{
		ResourceGroup: id.ResourceGroup,
		FrontDoorName: id.Name,
		Name:          name,
	}
}

func (id RoutingRuleId) ID(subscriptionId string) string {
	base := NewFrontDoorID(id.ResourceGroup, id.FrontDoorName).ID(subscriptionId)
	return fmt.Sprintf("%s/routingRules/%s", base, id.Name)
}

func RoutingRuleID(input string) (*RoutingRuleId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Routing Rule ID %q: %+v", input, err)
	}

	poolId := RoutingRuleId{
		ResourceGroup: frontDoorId.ResourceGroup,
		FrontDoorName: frontDoorId.Name,
	}

	// API is broken - https://github.com/Azure/azure-sdk-for-go/issues/6762
	// note: the ordering is important since the defined case (we want to error with) is routingRules
	if poolId.Name, err = id.PopSegment("RoutingRules"); err != nil {
		if poolId.Name, err = id.PopSegment("routingRules"); err != nil {
			return nil, err
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &poolId, nil
}
