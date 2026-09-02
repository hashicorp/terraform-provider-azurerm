// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PipelineGroupReceiverEndpoint validates a Syslog or OTLP receiver endpoint, which must be
// in `<host>:<port>` form, for example `0.0.0.0:514`. Semantic restrictions on the host (for
// example disallowed address ranges) are left to the service to enforce.
// https://learn.microsoft.com/en-us/azure/templates/microsoft.monitor/pipelinegroups
func PipelineGroupReceiverEndpoint(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %q to be string", k))
	}

	if v == "" {
		return nil, append(errors, fmt.Errorf("%q cannot be empty, expected a `<host>:<port>` endpoint such as `0.0.0.0:514`", k))
	}

	port, ok := pipelineGroupReceiverEndpointPort(v)
	if !ok {
		return nil, append(errors, fmt.Errorf("%q must be in the form `<host>:<port>`, got %q", k, v))
	}

	if portNumber, err := strconv.Atoi(port); err != nil || portNumber < 1 || portNumber > 65535 {
		errors = append(errors, fmt.Errorf("%q must have a numeric port between `1` and `65535`, got %q", k, v))
	}

	return warnings, errors
}

// pipelineGroupReceiverEndpointPort extracts the port component of an endpoint. It accepts
// `<host>:<port>`, `:<port>`, `[<ipv6>]:<port>`, and `<scheme>://<host>:<port>` forms.
func pipelineGroupReceiverEndpointPort(endpoint string) (port string, ok bool) {
	if idx := strings.Index(endpoint, "://"); idx != -1 {
		endpoint = endpoint[idx+3:]
		if slash := strings.IndexAny(endpoint, "/?#"); slash != -1 {
			endpoint = endpoint[:slash]
		}
	}

	if endpoint == "" {
		return "", false
	}

	if _, p, err := net.SplitHostPort(endpoint); err == nil {
		return p, true
	}

	if strings.HasPrefix(endpoint, ":") && !strings.Contains(endpoint[1:], ":") {
		return endpoint[1:], true
	}

	return "", false
}
