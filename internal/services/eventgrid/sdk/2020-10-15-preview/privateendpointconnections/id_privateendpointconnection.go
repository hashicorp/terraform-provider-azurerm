package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PrivateEndpointConnectionId{}

// PrivateEndpointConnectionId is a struct representing the Resource ID for a Private Endpoint Connection
type PrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroupName             string
	ParentType                    ParentType
	ParentName                    string
	PrivateEndpointConnectionName string
}

// NewPrivateEndpointConnectionID returns a new PrivateEndpointConnectionId struct
func NewPrivateEndpointConnectionID(subscriptionId string, resourceGroupName string, parentType ParentType, parentName string, privateEndpointConnectionName string) PrivateEndpointConnectionId {
	return PrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		ParentType:                    parentType,
		ParentName:                    parentName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

// ParsePrivateEndpointConnectionID parses 'input' into a PrivateEndpointConnectionId
func ParsePrivateEndpointConnectionID(input string) (*PrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateEndpointConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["parentType"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'parentType' was not found in the resource id %q", input)
		}

		parentType, err := parseParentType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.ParentType = *parentType
	}

	if id.ParentName, ok = parsed.Parsed["parentName"]; !ok {
		return nil, fmt.Errorf("the segment 'parentName' was not found in the resource id %q", input)
	}

	if id.PrivateEndpointConnectionName, ok = parsed.Parsed["privateEndpointConnectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateEndpointConnectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParsePrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a PrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParsePrivateEndpointConnectionIDInsensitively(input string) (*PrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateEndpointConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if v, constFound := parsed.Parsed["parentType"]; true {
		if !constFound {
			return nil, fmt.Errorf("the segment 'parentType' was not found in the resource id %q", input)
		}

		parentType, err := parseParentType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.ParentType = *parentType
	}

	if id.ParentName, ok = parsed.Parsed["parentName"]; !ok {
		return nil, fmt.Errorf("the segment 'parentName' was not found in the resource id %q", input)
	}

	if id.PrivateEndpointConnectionName, ok = parsed.Parsed["privateEndpointConnectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateEndpointConnectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidatePrivateEndpointConnectionID checks that 'input' can be parsed as a Private Endpoint Connection ID
func ValidatePrivateEndpointConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateEndpointConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Endpoint Connection ID
func (id PrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/%s/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, string(id.ParentType), id.ParentName, id.PrivateEndpointConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Endpoint Connection ID
func (id PrivateEndpointConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.ConstantSegment("parentType", PossibleValuesForParentType(), "domains"),
		resourceids.UserSpecifiedSegment("parentName", "parentValue"),
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionValue"),
	}
}

// String returns a human-readable description of this Private Endpoint Connection ID
func (id PrivateEndpointConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Parent Type: %q", string(id.ParentType)),
		fmt.Sprintf("Parent Name: %q", id.ParentName),
		fmt.Sprintf("Private Endpoint Connection Name: %q", id.PrivateEndpointConnectionName),
	}
	return fmt.Sprintf("Private Endpoint Connection (%s)", strings.Join(components, "\n"))
}
