// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func VirtualNetworkBgpCommunity(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	segments := strings.Split(v, ":")
	if len(segments) != 2 {
		errors = append(errors, fmt.Errorf(`invalid notation of bgp community: expected "x:y"`))
		return warnings, errors
	}

	asn, err := strconv.Atoi(segments[0])
	if err != nil {
		errors = append(errors, fmt.Errorf(`converting asn %q: %v`, segments[0], err))
		return warnings, errors
	}
	if asn <= 0 || asn >= 65535 {
		errors = append(errors, fmt.Errorf(`asn %d exceeds range: [0, 65535]`, asn))
		return warnings, errors
	}

	comm, err := strconv.Atoi(segments[1])
	if err != nil {
		errors = append(errors, fmt.Errorf(`converting community value %q: %v`, segments[1], err))
		return warnings, errors
	}
	if comm <= 0 || comm >= 65535 {
		errors = append(errors, fmt.Errorf(`community value %d exceeds range: [0, 65535]`, comm))
		return warnings, errors
	}
	return warnings, errors
}
