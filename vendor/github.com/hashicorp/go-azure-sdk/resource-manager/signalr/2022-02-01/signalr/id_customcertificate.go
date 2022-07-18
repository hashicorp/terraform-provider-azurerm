package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CustomCertificateId{}

// CustomCertificateId is a struct representing the Resource ID for a Custom Certificate
type CustomCertificateId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
	CertificateName   string
}

// NewCustomCertificateID returns a new CustomCertificateId struct
func NewCustomCertificateID(subscriptionId string, resourceGroupName string, resourceName string, certificateName string) CustomCertificateId {
	return CustomCertificateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceName:      resourceName,
		CertificateName:   certificateName,
	}
}

// ParseCustomCertificateID parses 'input' into a CustomCertificateId
func ParseCustomCertificateID(input string) (*CustomCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, fmt.Errorf("the segment 'certificateName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCustomCertificateIDInsensitively parses 'input' case-insensitively into a CustomCertificateId
// note: this method should only be used for API response data and not user input
func ParseCustomCertificateIDInsensitively(input string) (*CustomCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, fmt.Errorf("the segment 'certificateName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCustomCertificateID checks that 'input' can be parsed as a Custom Certificate ID
func ValidateCustomCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Custom Certificate ID
func (id CustomCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/signalR/%s/customCertificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Custom Certificate ID
func (id CustomCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticSignalR", "signalR", "signalR"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticCustomCertificates", "customCertificates", "customCertificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Custom Certificate ID
func (id CustomCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Custom Certificate (%s)", strings.Join(components, "\n"))
}
