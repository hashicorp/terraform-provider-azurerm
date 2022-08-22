package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
)

func ServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := servers.ParseServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
