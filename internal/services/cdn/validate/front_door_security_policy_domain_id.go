// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
)

func FrontDoorSecurityPolicyDomainID(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
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

	// Fix for issue #18551: If the ID failed to parse try it again but case insensitively.
	// I believe this is ok because it is just a pass through value and the service doesn't care
	// about the passed values casing...
	if err != nil {
		if _, err := parse.FrontDoorCustomDomainIDInsensitively(v); err == nil {
			return nil, nil
		}

		if _, err := parse.FrontDoorEndpointIDInsensitively(v); err == nil {
			return nil, nil
		}
	}

	return nil, []error{fmt.Errorf("%q is invalid: the %q needs to be a valid CDN Front Door Custom Domain ID or a valid CDN Front Door Endpoint ID: Unable to parse passed domain resource ID value %q", "domain", k, v)}
}
