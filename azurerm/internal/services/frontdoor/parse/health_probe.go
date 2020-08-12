package parse

import "fmt"

type HealthProbeId struct {
	ResourceGroup string
	FrontDoorName string
	Name          string
}

func NewHealthProbeID(id FrontDoorId, name string) HealthProbeId {
	return HealthProbeId{
		ResourceGroup: id.ResourceGroup,
		FrontDoorName: id.Name,
		Name:          name,
	}
}

func (id HealthProbeId) ID(subscriptionId string) string {
	base := NewFrontDoorID(id.ResourceGroup, id.Name).ID(subscriptionId)
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

	// TODO: handle this being case-insensitive
	// https://github.com/Azure/azure-sdk-for-go/issues/6762
	if probeId.Name, err = id.PopSegment("healthProbeSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &probeId, nil
}
