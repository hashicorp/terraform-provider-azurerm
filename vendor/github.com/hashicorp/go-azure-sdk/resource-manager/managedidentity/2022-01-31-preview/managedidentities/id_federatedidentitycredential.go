package managedidentities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = FederatedIdentityCredentialId{}

// FederatedIdentityCredentialId is a struct representing the Resource ID for a Federated Identity Credential
type FederatedIdentityCredentialId struct {
	SubscriptionId                          string
	ResourceGroupName                       string
	ResourceName                            string
	FederatedIdentityCredentialResourceName string
}

// NewFederatedIdentityCredentialID returns a new FederatedIdentityCredentialId struct
func NewFederatedIdentityCredentialID(subscriptionId string, resourceGroupName string, resourceName string, federatedIdentityCredentialResourceName string) FederatedIdentityCredentialId {
	return FederatedIdentityCredentialId{
		SubscriptionId:                          subscriptionId,
		ResourceGroupName:                       resourceGroupName,
		ResourceName:                            resourceName,
		FederatedIdentityCredentialResourceName: federatedIdentityCredentialResourceName,
	}
}

// ParseFederatedIdentityCredentialID parses 'input' into a FederatedIdentityCredentialId
func ParseFederatedIdentityCredentialID(input string) (*FederatedIdentityCredentialId, error) {
	parser := resourceids.NewParserFromResourceIdType(FederatedIdentityCredentialId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FederatedIdentityCredentialId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.FederatedIdentityCredentialResourceName, ok = parsed.Parsed["federatedIdentityCredentialResourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'federatedIdentityCredentialResourceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseFederatedIdentityCredentialIDInsensitively parses 'input' case-insensitively into a FederatedIdentityCredentialId
// note: this method should only be used for API response data and not user input
func ParseFederatedIdentityCredentialIDInsensitively(input string) (*FederatedIdentityCredentialId, error) {
	parser := resourceids.NewParserFromResourceIdType(FederatedIdentityCredentialId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FederatedIdentityCredentialId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.FederatedIdentityCredentialResourceName, ok = parsed.Parsed["federatedIdentityCredentialResourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'federatedIdentityCredentialResourceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateFederatedIdentityCredentialID checks that 'input' can be parsed as a Federated Identity Credential ID
func ValidateFederatedIdentityCredentialID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFederatedIdentityCredentialID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Federated Identity Credential ID
func (id FederatedIdentityCredentialId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s/federatedIdentityCredentials/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.FederatedIdentityCredentialResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Federated Identity Credential ID
func (id FederatedIdentityCredentialId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagedIdentity", "Microsoft.ManagedIdentity", "Microsoft.ManagedIdentity"),
		resourceids.StaticSegment("staticUserAssignedIdentities", "userAssignedIdentities", "userAssignedIdentities"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticFederatedIdentityCredentials", "federatedIdentityCredentials", "federatedIdentityCredentials"),
		resourceids.UserSpecifiedSegment("federatedIdentityCredentialResourceName", "federatedIdentityCredentialResourceValue"),
	}
}

// String returns a human-readable description of this Federated Identity Credential ID
func (id FederatedIdentityCredentialId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Federated Identity Credential Resource Name: %q", id.FederatedIdentityCredentialResourceName),
	}
	return fmt.Sprintf("Federated Identity Credential (%s)", strings.Join(components, "\n"))
}
