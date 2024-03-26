// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ManagedHSMKeyWithVersionId{}

type ManagedHSMKeyWithVersionId struct {
	SubscriptionId string
	ResourceGroup  string
	ManagedHSMName string
	KeyName        string
	VersionName    string
}

func NewManagedHSMKeyWithVersionID(subscriptionId, resourceGroup, managedHSMName, keyName, versionName string) ManagedHSMKeyWithVersionId {
	return ManagedHSMKeyWithVersionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ManagedHSMName: managedHSMName,
		KeyName:        keyName,
		VersionName:    versionName,
	}
}

// FromParseResult implements resourceids.ResourceId.
func (id *ManagedHSMKeyWithVersionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
}

// Segments implements resourceids.ResourceId.
func (id *ManagedHSMKeyWithVersionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticVersion", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionValue"),
	}
}

func (id ManagedHSMKeyWithVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M  Key With Version", segmentsStr)
}

func (id ManagedHSMKeyWithVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/keys/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.KeyName, id.VersionName)
}

// ManagedHSMKeyWithVersionID parses a ManagedHSMKeyWithVersion ID into an ManagedHSMKeyWithVersionId struct
func ParseManagedHSMKeyWithVersionID(input string) (*ManagedHSMKeyWithVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMKeyWithVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMKeyWithVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ParseManagedHSMKeyWithVersionIDInsensitively(input string) (*ManagedHSMKeyWithVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedHSMKeyWithVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedHSMKeyWithVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func ValidateManagedHSMKeyWithVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedHSMKeyWithVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
