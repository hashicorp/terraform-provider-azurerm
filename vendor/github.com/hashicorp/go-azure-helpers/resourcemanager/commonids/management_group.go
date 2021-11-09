package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ManagementGroupId{}

type ManagementGroupId struct {
	GroupId string
}

func NewManagementGroupID(groupId string) ManagementGroupId {
	return ManagementGroupId{
		GroupId: groupId,
	}
}

func ParseManagementGroupID(input string) (*ManagementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagementGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagementGroupId{}

	if id.GroupId, ok = parsed.Parsed["groupId"]; !ok {
		return nil, fmt.Errorf("the segment 'groupId' was not found in the resource id %q", input)
	}

	return &id, nil
}

func ParseManagementGroupIDInsensitively(input string) (*ManagementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagementGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagementGroupId{}

	if id.GroupId, ok = parsed.Parsed["groupId"]; !ok {
		return nil, fmt.Errorf("the segment 'groupId' was not found in the resource id %q", input)
	}

	return &id, nil
}

func (id ManagementGroupId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s"
	return fmt.Sprintf(fmtString, id.GroupId)
}

func (id ManagementGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("managementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("groupId", "groupIdValue"),
	}
}

func (id ManagementGroupId) String() string {
	components := []string{
		fmt.Sprintf("Group: %q", id.GroupId),
	}
	return fmt.Sprintf("Management Group (%s)", strings.Join(components, "\n"))
}
