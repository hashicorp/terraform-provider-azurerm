// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId struct {
	Scope string
}

func NewSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(scope string) SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId {
	return SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{
		Scope: scope,
	}
}

func SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(input string) (*SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentIDInsensitively(input string) (*SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	return nil
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) ID() string {
	fmtString := "/%s/providers/Microsoft.ScVmm/virtualMachineInstances/default/guestAgents/default"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"))
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.HybridCompute/machines/some-machine"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticVirtualMachineInstances", "virtualMachineInstances", "virtualMachineInstances"),
		resourceids.StaticSegment("staticVirtualMachineInstanceDefault", "default", "default"),
		resourceids.StaticSegment("staticGuestAgents", "guestAgents", "guestAgents"),
		resourceids.StaticSegment("staticGuestAgentDefault", "default", "default"),
	}
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
	}
	return fmt.Sprintf("System Center Virtual Machine Manager Virtual Machine Instance Guest Agent (%s)", strings.Join(components, "\n"))
}
