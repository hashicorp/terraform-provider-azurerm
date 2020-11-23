package parse

import "fmt"

type HealthProbeId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewHealthProbeID(id FrontDoorId, name string) HealthProbeId {
	return HealthProbeId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		FrontDoorName:  id.Name,
		Name:           name,
	}
}

func (id HealthProbeId) ID(_ string) string {
	base := NewFrontDoorID(id.SubscriptionId, id.ResourceGroup, id.FrontDoorName).ID("")
	return fmt.Sprintf("%s/healthProbeSettings/%s", base, id.Name)
}

func HealthProbeID(input string) (*HealthProbeId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor Health Probe ID %q: %+v", input, err)
	}

	probeId := HealthProbeId{
		SubscriptionId: frontDoorId.SubscriptionId,
		ResourceGroup:  frontDoorId.ResourceGroup,
		FrontDoorName:  frontDoorId.Name,
	}

	// https://github.com/Azure/azure-sdk-for-go/issues/6762
	// note: the ordering is important since the defined case (we want to error with) is healthProbeSettings
	if probeId.Name, err = id.PopSegment("HealthProbeSettings"); err != nil {
		if probeId.Name, err = id.PopSegment("healthProbeSettings"); err != nil {
			return nil, err
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &probeId, nil
}
