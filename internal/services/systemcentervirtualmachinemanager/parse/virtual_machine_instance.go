package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SystemCenterVirtualMachineManagerVirtualMachineInstanceId{}

type SystemCenterVirtualMachineManagerVirtualMachineInstanceId struct {
	ScopeId commonids.ScopeId
}

func NewSystemCenterVirtualMachineManagerVirtualMachineInstanceID(scopeId string) SystemCenterVirtualMachineManagerVirtualMachineInstanceId {
	return SystemCenterVirtualMachineManagerVirtualMachineInstanceId{
		ScopeId: commonids.ScopeId{
			Scope: scopeId,
		},
	}
}

func SystemCenterVirtualMachineManagerVirtualMachineInstanceID(input string) (*SystemCenterVirtualMachineManagerVirtualMachineInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SystemCenterVirtualMachineManagerVirtualMachineInstanceId{})
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

func (id *SystemCenterVirtualMachineManagerVirtualMachineInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ScopeId.Scope, ok = input.Parsed["scopeId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scopeId", input)
	}

	return nil
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceId) ID() string {
	fmtString := "/%s/providers/Microsoft.ScVmm/virtualMachineInstances/default"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.ScopeId.Scope, "/"))
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scopeId", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.HybridCompute/machines/machine1"),
	}
}

func (id SystemCenterVirtualMachineManagerVirtualMachineInstanceId) String() string {
	components := []string{
		fmt.Sprintf("ScopeId: %q", id.ScopeId),
	}
	return fmt.Sprintf("System Center Virtual Machine Manager Virtual Machine Instance (%s)", strings.Join(components, "\n"))
}
