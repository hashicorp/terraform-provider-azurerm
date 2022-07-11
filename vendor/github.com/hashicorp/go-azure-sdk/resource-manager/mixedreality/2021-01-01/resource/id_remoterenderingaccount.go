package resource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = RemoteRenderingAccountId{}

// RemoteRenderingAccountId is a struct representing the Resource ID for a Remote Rendering Account
type RemoteRenderingAccountId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
}

// NewRemoteRenderingAccountID returns a new RemoteRenderingAccountId struct
func NewRemoteRenderingAccountID(subscriptionId string, resourceGroupName string, accountName string) RemoteRenderingAccountId {
	return RemoteRenderingAccountId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
	}
}

// ParseRemoteRenderingAccountID parses 'input' into a RemoteRenderingAccountId
func ParseRemoteRenderingAccountID(input string) (*RemoteRenderingAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(RemoteRenderingAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RemoteRenderingAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseRemoteRenderingAccountIDInsensitively parses 'input' case-insensitively into a RemoteRenderingAccountId
// note: this method should only be used for API response data and not user input
func ParseRemoteRenderingAccountIDInsensitively(input string) (*RemoteRenderingAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(RemoteRenderingAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RemoteRenderingAccountId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateRemoteRenderingAccountID checks that 'input' can be parsed as a Remote Rendering Account ID
func ValidateRemoteRenderingAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRemoteRenderingAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Remote Rendering Account ID
func (id RemoteRenderingAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MixedReality/remoteRenderingAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Remote Rendering Account ID
func (id RemoteRenderingAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMixedReality", "Microsoft.MixedReality", "Microsoft.MixedReality"),
		resourceids.StaticSegment("staticRemoteRenderingAccounts", "remoteRenderingAccounts", "remoteRenderingAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
	}
}

// String returns a human-readable description of this Remote Rendering Account ID
func (id RemoteRenderingAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
	}
	return fmt.Sprintf("Remote Rendering Account (%s)", strings.Join(components, "\n"))
}
