package checkfrontdoornameavailability

import "strings"

type Availability string

const (
	AvailabilityAvailable   Availability = "Available"
	AvailabilityUnavailable Availability = "Unavailable"
)

func PossibleValuesForAvailability() []string {
	return []string{
		string(AvailabilityAvailable),
		string(AvailabilityUnavailable),
	}
}

func parseAvailability(input string) (*Availability, error) {
	vals := map[string]Availability{
		"available":   AvailabilityAvailable,
		"unavailable": AvailabilityUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Availability(input)
	return &out, nil
}

type ResourceType string

const (
	ResourceTypeMicrosoftPointNetworkFrontDoors                  ResourceType = "Microsoft.Network/frontDoors"
	ResourceTypeMicrosoftPointNetworkFrontDoorsFrontendEndpoints ResourceType = "Microsoft.Network/frontDoors/frontendEndpoints"
)

func PossibleValuesForResourceType() []string {
	return []string{
		string(ResourceTypeMicrosoftPointNetworkFrontDoors),
		string(ResourceTypeMicrosoftPointNetworkFrontDoorsFrontendEndpoints),
	}
}

func parseResourceType(input string) (*ResourceType, error) {
	vals := map[string]ResourceType{
		"microsoft.network/frontdoors":                   ResourceTypeMicrosoftPointNetworkFrontDoors,
		"microsoft.network/frontdoors/frontendendpoints": ResourceTypeMicrosoftPointNetworkFrontDoorsFrontendEndpoints,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceType(input)
	return &out, nil
}
