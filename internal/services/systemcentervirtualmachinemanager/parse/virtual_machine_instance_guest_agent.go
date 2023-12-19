package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId struct {
	Scope                      string
	VirtualMachineInstanceName string
	GuestAgentName             string
}

func NewSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(scope string, virtualMachineInstanceName string, guestAgentName string) SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId {
	return SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{
		Scope:                      scope,
		VirtualMachineInstanceName: virtualMachineInstanceName,
		GuestAgentName:             guestAgentName,
	}
}

func SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(input string) (*SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId, error) {
	parser := resourceids.NewParserFromResourceIdType(SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{})
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

func (id *SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.VirtualMachineInstanceName, ok = input.Parsed["virtualMachineInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineInstanceName", input)
	}

	if id.GuestAgentName, ok = input.Parsed["guestAgentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "guestAgentName", input)
	}

	return nil
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) ID() string {
	fmtString := "/%s/providers/Microsoft.ScVmm/virtualMachineInstances/%s/guestAgents/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.VirtualMachineInstanceName, id.GuestAgentName)
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.HybridCompute/machines/machine1"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticVirtualMachineInstances", "virtualMachineInstances", "virtualMachineInstances"),
		resourceids.UserSpecifiedSegment("virtualMachineInstanceName", "virtualMachineInstanceValue"),
		resourceids.StaticSegment("staticGuestAgents", "guestAgents", "guestAgents"),
		resourceids.UserSpecifiedSegment("guestAgentName", "guestAgentValue"),
	}
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Virtual Machine Instance Name: %q", id.VirtualMachineInstanceName),
		fmt.Sprintf("Guest Agent Name: %q", id.GuestAgentName),
	}
	return fmt.Sprintf("System Center Virtual Machine Manager Virtual Machine Instance Guest Agent (%s)", strings.Join(components, "\n"))
}
