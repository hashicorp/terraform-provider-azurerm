package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId struct {
	ScopeId commonids.ScopeId
}

func NewSystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(scopeId string) SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId {
	return SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId{
		ScopeId: commonids.ScopeId{
			Scope: scopeId,
		},
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

	if id.ScopeId.Scope, ok = input.Parsed["scopeId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scopeId", input)
	}

	return nil
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) ID() string {
	fmtString := "/%s/providers/Microsoft.ScVmm/virtualMachineInstances/default/guestAgents/default"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ScopeId.Scope, "/"))
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scopeId", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.HybridCompute/machines/machine1"),
	}
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentId) String() string {
	components := []string{
		fmt.Sprintf("ScopeId: %q", id.ScopeId),
	}
	return fmt.Sprintf("System Center Virtual Machine Manager Virtual Machine Instance Guest Agent (%s)", strings.Join(components, "\n"))
}
