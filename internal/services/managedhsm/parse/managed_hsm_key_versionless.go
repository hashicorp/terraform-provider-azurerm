// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ManagedHSMKeyVersionlessId{}

type ManagedHSMKeyVersionlessId struct {
	SubscriptionId string
	ResourceGroup  string
	ManagedHSMName string
	KeyName        string
}

func NewManagedHSMKeyVersionlessId(subscriptionId, resourceGroup, managedHSMName, keyName string) ManagedHSMKeyVersionlessId {
	return ManagedHSMKeyVersionlessId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ManagedHSMName: managedHSMName,
		KeyName:        keyName,
	}
}

// FromParseResult implements resourceids.ResourceId.
func (id *ManagedHSMKeyVersionlessId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroup, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedHSMName, ok = input.Parsed["managedHSMName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedHSMName", input)
	}

	if id.KeyName, ok = input.Parsed["keyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "keyName", input)
	}

	return nil
}

// Segments implements resourceids.ResourceId.
func (id *ManagedHSMKeyVersionlessId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticManagedHSMs", "managedHSMs", "managedHSMs"),
		resourceids.UserSpecifiedSegment("managedHSMName", "managedHSMValue"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyValue"),
	}
}

func (id ManagedHSMKeyVersionlessId) String() string {
	segments := []string{
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M  Key With Version", segmentsStr)
}

func (id ManagedHSMKeyVersionlessId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.KeyName)
}

// ManagedHSMKeyVersionlessId parses a ManagedHSMKeyWithVersion ID into an ManagedHSMKeyVersionlessId struct
func ParseManagedHSMKeyVersionlessId(input string) (*ManagedHSMKeyVersionlessId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMKeyVersionlessId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMKeyVersionlessId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ParseManagedHSMKeyVersionlessIdInsensitively(input string) (*ManagedHSMKeyVersionlessId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMKeyVersionlessId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMKeyVersionlessId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}
func ValidateManagedHSMKeyVersionlessID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedHSMKeyVersionlessId(v); err != nil {
		errors = append(errors, err)
	}

	return
}
