package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SystemCenterVirtualMachineManagerVirtualMachineInstanceId{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceId struct {
	Scope string
}

func NewSystemCenterVirtualMachineManagerVirtualMachineInstanceID(scope string) SystemCenterVirtualMachineManagerVirtualMachineInstanceId {
	return SystemCenterVirtualMachineManagerVirtualMachineInstanceId{
		Scope: scope,
	}
}

func SystemCenterVirtualMachineManagerVirtualMachineInstanceID(input string) (*SystemCenterVirtualMachineManagerVirtualMachineInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemCenterVirtualMachineManagerVirtualMachineInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemCenterVirtualMachineManagerVirtualMachineInstanceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func SystemCenterVirtualMachineManagerVirtualMachineInstanceIDInsensitively(input string) (*SystemCenterVirtualMachineManagerVirtualMachineInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemCenterVirtualMachineManagerVirtualMachineInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemCenterVirtualMachineManagerVirtualMachineInstanceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SystemCenterVirtualMachineManagerVirtualMachineInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	return nil
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceId) ID() string {
	fmtString := "/%s/providers/Microsoft.ScVmm/virtualMachineInstances/default"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"))
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.HybridCompute/machines/some-machine"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticVirtualMachineInstances", "virtualMachineInstances", "virtualMachineInstances"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
	}
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
	}
	return fmt.Sprintf("System Center Virtual Machine Manager Virtual Machine Instance (%s)", strings.Join(components, "\n"))
}
