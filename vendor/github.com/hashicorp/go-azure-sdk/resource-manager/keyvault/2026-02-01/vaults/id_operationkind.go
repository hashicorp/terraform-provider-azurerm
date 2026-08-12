package vaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OperationKindId{})
}

var _ resourceids.ResourceId = &OperationKindId{}

// OperationKindId is a struct representing the Resource ID for a Operation Kind
type OperationKindId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
	OperationKind     AccessPolicyUpdateKind
}

// NewOperationKindID returns a new OperationKindId struct
func NewOperationKindID(subscriptionId string, resourceGroupName string, vaultName string, operationKind AccessPolicyUpdateKind) OperationKindId {
	return OperationKindId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
		OperationKind:     operationKind,
	}
}

// ParseOperationKindID parses 'input' into a OperationKindId
func ParseOperationKindID(input string) (*OperationKindId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationKindId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationKindId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOperationKindIDInsensitively parses 'input' case-insensitively into a OperationKindId
// note: this method should only be used for API response data and not user input
func ParseOperationKindIDInsensitively(input string) (*OperationKindId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationKindId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationKindId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OperationKindId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VaultName, ok = input.Parsed["vaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vaultName", input)
	}

	if v, ok := input.Parsed["operationKind"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "operationKind", input)
		}

		operationKind, err := parseAccessPolicyUpdateKind(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.OperationKind = *operationKind
	}

	return nil
}

// ValidateOperationKindID checks that 'input' can be parsed as a Operation Kind ID
func ValidateOperationKindID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationKindID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Kind ID
func (id OperationKindId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, string(id.OperationKind))
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Kind ID
func (id OperationKindId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticAccessPolicies", "accessPolicies", "accessPolicies"),
		resourceids.ConstantSegment("operationKind", PossibleValuesForAccessPolicyUpdateKind(), "add"),
	}
}

// String returns a human-readable description of this Operation Kind ID
func (id OperationKindId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Operation Kind: %q", string(id.OperationKind)),
	}
	return fmt.Sprintf("Operation Kind (%s)", strings.Join(components, "\n"))
}
