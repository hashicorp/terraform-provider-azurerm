package parse

import "fmt"

type RoutingRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewRoutingRuleID(id FrontDoorId, name string) RoutingRuleId {
	return RoutingRuleId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		FrontDoorName:  id.Name,
		Name:           name,
	}
}

func (id RoutingRuleId) ID(_ string) string {
	base := NewFrontDoorID(id.SubscriptionId, id.ResourceGroup, id.FrontDoorName).ID("")
	return fmt.Sprintf("%s/routingRules/%s", base, id.Name)
}

func RoutingRuleID(input string) (*RoutingRuleId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Routing Rule ID %q: %+v", input, err)
	}

	poolId := RoutingRuleId{
		SubscriptionId: frontDoorId.SubscriptionId,
		ResourceGroup:  frontDoorId.ResourceGroup,
		FrontDoorName:  frontDoorId.Name,
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
