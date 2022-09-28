package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
)

// ValidateDnsZoneIDInsensitively checks that 'input' can be parsed as a Dns Zone ID
func ValidateDnsZoneIDInsensitively(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := zones.ParseDnsZoneIDInsensitively(v); err != nil {
		errors = append(errors, err)
	}

	return
}
