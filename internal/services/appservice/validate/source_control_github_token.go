package validate

import "fmt"

const expectedID = "/providers/Microsoft.Web/sourcecontrols/GitHub"

// TODO - Should this be genericised for the other 3 Possible token types (DropBox, BitBucket, and OneDrive), or one per?

func SourceControlGitHubTokenID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}
	if v != expectedID {
		errors = append(errors, fmt.Errorf("ID must be exactly %q", expectedID))
	}
	return
}
