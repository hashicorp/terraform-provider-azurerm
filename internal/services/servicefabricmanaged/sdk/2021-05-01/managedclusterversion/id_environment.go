package managedclusterversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = EnvironmentId{}

// EnvironmentId is a struct representing the Resource ID for a Environment
type EnvironmentId struct {
	SubscriptionId string
	Location       string
	Environment    Environment
}

// NewEnvironmentID returns a new EnvironmentId struct
func NewEnvironmentID(subscriptionId string, location string, environment Environment) EnvironmentId {
	return EnvironmentId{
		SubscriptionId: subscriptionId,
		Location:       location,
		Environment:    environment,
	}
}

// ParseEnvironmentID parses 'input' into a EnvironmentId
func ParseEnvironmentID(input string) (*EnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(EnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EnvironmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["environment"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'environment' was not found in the resource id %q", input)
		}

		environment, err := parseEnvironment(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.Environment = *environment
	}

	return &id, nil
}

// ParseEnvironmentIDInsensitively parses 'input' case-insensitively into a EnvironmentId
// note: this method should only be used for API response data and not user input
func ParseEnvironmentIDInsensitively(input string) (*EnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(EnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EnvironmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.Location, ok = parsed.Parsed["location"]; !ok {
		return nil, fmt.Errorf("the segment 'location' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["environment"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'environment' was not found in the resource id %q", input)
		}

		environment, err := parseEnvironment(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.Environment = *environment
	}

	return &id, nil
}

// ValidateEnvironmentID checks that 'input' can be parsed as a Environment ID
func ValidateEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Environment ID
func (id EnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ServiceFabric/locations/%s/environments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Location, string(id.Environment))
}

// Segments returns a slice of Resource ID Segments which comprise this Environment ID
func (id EnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftServiceFabric", "Microsoft.ServiceFabric", "Microsoft.ServiceFabric"),
		resourceids.StaticSegment("locations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("location", "locationValue"),
		resourceids.StaticSegment("environments", "environments", "environments"),
		resourceids.ConstantSegment("environment", PossibleValuesForEnvironment(), "Windows"),
	}
}

// String returns a human-readable description of this Environment ID
func (id EnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location: %q", id.Location),
		fmt.Sprintf("Environment: %q", string(id.Environment)),
	}
	return fmt.Sprintf("Environment (%s)", strings.Join(components, "\n"))
}
