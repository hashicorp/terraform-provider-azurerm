package parse

import (
	"fmt"
	"regexp"
)

type ManagementGroupId struct {
	GroupID string
}

func ManagementGroupID(input string) (*ManagementGroupId, error) {
	regex := regexp.MustCompile(`^/providers/[Mm]icrosoft\.[Mm]anagement/management[Gg]roups/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q", input)
	}

	// Split the input ID by the regex
	segments := regex.Split(input, -1)
	if len(segments) != 2 {
		return nil, fmt.Errorf("Unable to parse Management Group ID %q: expected id to have two segments after splitting", input)
	}
	groupID := segments[1]
	if _, errs := ValidateManagementGroupName(groupID, ""); len(errs) != 0 {
		return nil, fmt.Errorf("Unable to validate Management Group ID %q: %+v", input, errs)
	}

	id := ManagementGroupId{
		GroupID: groupID,
	}

	return &id, nil
}

func ValidateManagementGroupName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// portal says: The name can only be an ASCII letter, digit, -, _, (, ), . and have a maximum length constraint of 90
	if matched := regexp.MustCompile(`^[a-zA-Z0-9_().-]{1,90}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s can only consist of ASCII letters, digits, -, _, (, ), . , and cannot exceed the maximum length of 90", k))
	}
	return
}
