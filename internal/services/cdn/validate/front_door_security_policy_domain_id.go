package validate

import (
	"fmt"
)

func FrontDoorSecurityPolicyDomainID(i interface{}, k string) (_ []string, errors []error) {
	_, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q is invalid: expected type of %q to be string", "domain", k)}
	}

	var err []error

	if _, err = FrontDoorCustomDomainID(i, k); err == nil {
		return nil, nil
	}

	if _, err = FrontDoorEndpointID(i, k); err == nil {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("%q is invalid: the %q needs to be a valid Frontdoor Custom Domain ID or a Frontdoor Endpoint ID: %+v", "domain", k, err)}
}
