package secrets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SecretId{}

// SecretId is a struct representing the Resource ID for a Secret
type SecretId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	SecretName        string
}

// NewSecretID returns a new SecretId struct
func NewSecretID(subscriptionId string, resourceGroupName string, profileName string, secretName string) SecretId {
	return SecretId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		SecretName:        secretName,
	}
}

// ParseSecretID parses 'input' into a SecretId
func ParseSecretID(input string) (*SecretId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecretId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecretId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.SecretName, ok = parsed.Parsed["secretName"]; !ok {
		return nil, fmt.Errorf("the segment 'secretName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSecretIDInsensitively parses 'input' case-insensitively into a SecretId
// note: this method should only be used for API response data and not user input
func ParseSecretIDInsensitively(input string) (*SecretId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecretId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecretId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.SecretName, ok = parsed.Parsed["secretName"]; !ok {
		return nil, fmt.Errorf("the segment 'secretName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSecretID checks that 'input' can be parsed as a Secret ID
func ValidateSecretID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecretID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Secret ID
func (id SecretId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/secrets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.SecretName)
}

// Segments returns a slice of Resource ID Segments which comprise this Secret ID
func (id SecretId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticSecrets", "secrets", "secrets"),
		resourceids.UserSpecifiedSegment("secretName", "secretValue"),
	}
}

// String returns a human-readable description of this Secret ID
func (id SecretId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Secret Name: %q", id.SecretName),
	}
	return fmt.Sprintf("Secret (%s)", strings.Join(components, "\n"))
}
