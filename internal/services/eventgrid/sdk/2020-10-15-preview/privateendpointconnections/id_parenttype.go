package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ParentTypeId{}

// ParentTypeId is a struct representing the Resource ID for a Parent Type
type ParentTypeId struct {
	SubscriptionId    string
	ResourceGroupName string
	ParentType        ParentType
	ParentName        string
}

// NewParentTypeID returns a new ParentTypeId struct
func NewParentTypeID(subscriptionId string, resourceGroupName string, parentType ParentType, parentName string) ParentTypeId {
	return ParentTypeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ParentType:        parentType,
		ParentName:        parentName,
	}
}

// ParseParentTypeID parses 'input' into a ParentTypeId
func ParseParentTypeID(input string) (*ParentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ParentTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ParentTypeId{}

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

	return &id, nil
}

// ParseParentTypeIDInsensitively parses 'input' case-insensitively into a ParentTypeId
// note: this method should only be used for API response data and not user input
func ParseParentTypeIDInsensitively(input string) (*ParentTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ParentTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ParentTypeId{}

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

	return &id, nil
}

// ValidateParentTypeID checks that 'input' can be parsed as a Parent Type ID
func ValidateParentTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseParentTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Parent Type ID
func (id ParentTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, string(id.ParentType), id.ParentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Parent Type ID
func (id ParentTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.ConstantSegment("parentType", PossibleValuesForParentType(), "domains"),
		resourceids.UserSpecifiedSegment("parentName", "parentValue"),
	}
}

// String returns a human-readable description of this Parent Type ID
func (id ParentTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Parent Type: %q", string(id.ParentType)),
		fmt.Sprintf("Parent Name: %q", id.ParentName),
	}
	return fmt.Sprintf("Parent Type (%s)", strings.Join(components, "\n"))
}
