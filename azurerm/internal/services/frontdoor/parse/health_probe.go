package parse

import "fmt"

type HealthProbeId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewHealthProbeID(subscriptionId, resourceGroup, frontDoorName, name string) HealthProbeId {
	return HealthProbeId{
		ResourceGroup: resourceGroup,
		FrontDoorName: frontDoorName,
		Name:          name,
	}
}

func (id HealthProbeId) ID(subscriptionId string) string {
	base := NewFrontDoorID(subscriptionId, id.ResourceGroup, id.FrontDoorName).ID(subscriptionId)
	return fmt.Sprintf("%s/healthProbeSettings/%s", base, id.Name)
}

func HealthProbeID(input string) (*HealthProbeId, error) {
	frontDoorId, id, err := parseFrontDoorChildResourceId(input)
	if err != nil {
		return nil, fmt.Errorf("parsing FrontDoor Health Probe ID %q: %+v", input, err)
	}

	probeId := HealthProbeId{
		ResourceGroup: frontDoorId.ResourceGroup,
		FrontDoorName: frontDoorId.Name,
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
