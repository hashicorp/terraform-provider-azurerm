package encryptionscopes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = EncryptionScopeId{}

// EncryptionScopeId is a struct representing the Resource ID for a Encryption Scope
type EncryptionScopeId struct {
	SubscriptionId      string
	ResourceGroupName   string
	AccountName         string
	EncryptionScopeName string
}

// NewEncryptionScopeID returns a new EncryptionScopeId struct
func NewEncryptionScopeID(subscriptionId string, resourceGroupName string, accountName string, encryptionScopeName string) EncryptionScopeId {
	return EncryptionScopeId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		AccountName:         accountName,
		EncryptionScopeName: encryptionScopeName,
	}
}

// ParseEncryptionScopeID parses 'input' into a EncryptionScopeId
func ParseEncryptionScopeID(input string) (*EncryptionScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(EncryptionScopeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EncryptionScopeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.EncryptionScopeName, ok = parsed.Parsed["encryptionScopeName"]; !ok {
		return nil, fmt.Errorf("the segment 'encryptionScopeName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseEncryptionScopeIDInsensitively parses 'input' case-insensitively into a EncryptionScopeId
// note: this method should only be used for API response data and not user input
func ParseEncryptionScopeIDInsensitively(input string) (*EncryptionScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(EncryptionScopeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EncryptionScopeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.EncryptionScopeName, ok = parsed.Parsed["encryptionScopeName"]; !ok {
		return nil, fmt.Errorf("the segment 'encryptionScopeName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateEncryptionScopeID checks that 'input' can be parsed as a Encryption Scope ID
func ValidateEncryptionScopeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEncryptionScopeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Encryption Scope ID
func (id EncryptionScopeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/encryptionScopes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.EncryptionScopeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Encryption Scope ID
func (id EncryptionScopeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticEncryptionScopes", "encryptionScopes", "encryptionScopes"),
		resourceids.UserSpecifiedSegment("encryptionScopeName", "encryptionScopeValue"),
	}
}

// String returns a human-readable description of this Encryption Scope ID
func (id EncryptionScopeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Encryption Scope Name: %q", id.EncryptionScopeName),
	}
	return fmt.Sprintf("Encryption Scope (%s)", strings.Join(components, "\n"))
}
