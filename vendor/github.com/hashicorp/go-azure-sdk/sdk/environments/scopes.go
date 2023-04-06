// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func Resource(endpoint Api) (*string, error) {
	resource, ok := endpoint.ResourceIdentifier()
	if !ok {
		// if this API has been defined in-line it may not have a Resource Identifier - however
		// the resource will be the endpoint instead, so we can best-effort to obtain the auth token
		resource2, ok2 := endpoint.Endpoint()
		if !ok2 {
			return nil, fmt.Errorf("the endpoint %q is not supported in this Azure Environment", endpoint.Name())
		}
		resource = resource2
	}

	return resource, nil
}

func Scope(endpoint Api) (*string, error) {
	e, ok := endpoint.ResourceIdentifier()
	if !ok {
		// if this API has been defined in-line it may not have a Resource Identifier - however
		// the scope will be the endpoint instead, so we can best-effort to obtain the auth token
		e2, ok2 := endpoint.Endpoint()
		if !ok2 {
			return nil, fmt.Errorf("the endpoint %q is not supported in this Azure Environment", endpoint.Name())
		}
		e = e2
	}
	out := fmt.Sprintf("%s/.default", *e)
	return pointer.To(out), nil
}
