// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IsCommaSeparatedCIDRs(i interface{}, k string) (warnings []string, errors []error) {
	s, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			errors = append(errors, fmt.Errorf("%q: contains an empty CIDR entry", k))
			continue
		}
		if ws, errs := validation.IsCIDR(p, k); len(warnings) > 0 || len(errs) > 0 {
			warnings = append(warnings, ws...)
			errors = append(errors, errs...)
		}
	}
	return
}

func DomainName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	re := regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	if !re.MatchString(v) {
		errors = append(errors, fmt.Errorf("`%s` must be a valid domain name, for example, \"oracle.local\"", k))
	}
	return
}

func DomainNames(i interface{}, k string) (warnings []string, errors []error) {
	s, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			errors = append(errors, fmt.Errorf("%q: contains an empty Domain Name", k))
			continue
		}
		if ws, errs := DomainName(p, k); len(warnings) > 0 || len(errs) > 0 {
			warnings = append(warnings, ws...)
			errors = append(errors, errs...)
		}
	}
	return
}
