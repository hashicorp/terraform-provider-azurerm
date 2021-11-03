package checkfrontdoornameavailabilitywithsubscription

type Availability string

const (
	AvailabilityAvailable   Availability = "Available"
	AvailabilityUnavailable Availability = "Unavailable"
)

type ResourceType string

const (
	ResourceTypeMicrosoftPointNetworkFrontDoors                  ResourceType = "Microsoft.Network/frontDoors"
	ResourceTypeMicrosoftPointNetworkFrontDoorsFrontendEndpoints ResourceType = "Microsoft.Network/frontDoors/frontendEndpoints"
)
