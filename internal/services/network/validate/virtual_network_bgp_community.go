// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func VirtualNetworkBgpCommunity(i any, k string) (warnings []string, errs []error) {
	v, ok := i.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	segments := strings.Split(v, ":")
	if len(segments) != 2 {
		errs = append(errs, errors.New(`invalid notation of bgp community: expected "x:y"`))
		return
	}

	asn, err := strconv.Atoi(segments[0])
	if err != nil {
		errs = append(errs, fmt.Errorf(`converting asn %q: %v`, segments[0], err))
		return
	}
	if asn <= 0 || asn >= 65535 {
		errs = append(errs, fmt.Errorf(`asn %d exceeds range: [0, 65535]`, asn))
		return
	}

	comm, err := strconv.Atoi(segments[1])
	if err != nil {
		errs = append(errs, fmt.Errorf(`converting community value %q: %v`, segments[1], err))
		return
	}
	if comm <= 0 || comm >= 65535 {
		errs = append(errs, fmt.Errorf(`community value %d exceeds range: [0, 65535]`, comm))
		return
	}
	return warnings, errs
}
