package parse

import (
	"fmt"
	"regexp"
	"strings"
)

type ManagementGroupId struct {
	Name string
}

func ManagementGroupID(input string) (*ManagementGroupId, error) {
	regex := regexp.MustCompile(`^/providers/[Mm]icrosoft\.[Mm]anagement/[Mm]anagement[Gg]roups/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q", input)
	}

	// Split the input ID by the regex
	segments := regex.Split(input, -1)
	if len(segments) != 2 {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q: expected id to have two segments after splitting", input)
	}

	groupID := segments[1]
	if groupID == "" {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: management group name is empty", input)
	}
	if segments := strings.Split(groupID, "/"); len(segments) != 1 {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: ID has extra segments", input)
	}

	id := ManagementGroupId{
		Name: groupID,
	}

	return &id, nil
}

func NewManagementGroupId(managementGroupName string) ManagementGroupId {
	return ManagementGroupId{
		Name: managementGroupName,
	}
}

func (r ManagementGroupId) ID() string {
	managementGroupIdFmt := "/providers/Microsoft.Management/managementGroups/%s"
	return fmt.Sprintf(managementGroupIdFmt, r.Name)
}
