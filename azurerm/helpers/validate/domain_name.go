package validate

import (
	"fmt"
	"regexp"
)

// Fqdn validates that a domain name, including the host portion, is valid to RFC requirements
// e.g. portal.azure.com
func DomainName(i interface{}, k string) (warnings []string, errors []error) {
	//
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be string", k))
		return
	}
	// The following regexp is a good example of why pre-processing support would be nice in GoLang
	if matched := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,62}\.)|(?:[0-9]+/[0-9]{2})\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`).Match([]byte(v)); !matched || len(v) > 253 {
		errors = append(errors, fmt.Errorf("%q must be a valid CNAME", k))
	}

	return
}
