// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"fmt"
	"strings"
)

type resourceId interface {
	ID() string
}

// AsGeneratedID returns a SchemaValidateFunc which validates the value with the supplied
// case-insensitive resource ID parser (a `Parse<Type>IDInsensitively` function), whilst only
// tolerating the two casing deviations the legacy resourceids.ParseAzureResourceID parser accepted:
// an all-lowercase `resourcegroups` static segment and any casing of the resource provider
// namespace. Every other segment must match the canonical casing, exactly as the legacy parser
// required.
//
// todo 6.0 - remove this and move the call sites to the strict SDK `Validate<Type>ID` functions.
// The call sites were migrated from locally generated validators built on the legacy
// resourceids.ParseAzureResourceID, which treated the casing of the `resourceGroups` segment and
// the resource provider namespace insensitively, unlike the SDK validation functions. IDs with that
// non-canonical casing can therefore exist in configurations and in state (e.g. from older
// imports), so switching to case-sensitive validation requires a state migration on the referenced
// resources first.
func AsGeneratedID[T resourceId](parser func(string) (T, error)) func(interface{}, string) ([]string, []error) {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
		}

		id, err := parser(v)
		if err != nil {
			return nil, []error{err}
		}

		if err := checkLegacyCasing(v, id.ID()); err != nil {
			return nil, []error{fmt.Errorf("parsing %q: %+v", v, err)}
		}

		return nil, nil
	}
}

// checkLegacyCasing compares an insensitively parsed input against the canonical form of the ID and
// permits only the casing deviations the legacy resourceids.ParseAzureResourceID parser tolerated:
// the all-lowercase `resourcegroups` fallback and the casing of the resource provider namespace.
func checkLegacyCasing(input string, canonical string) error {
	inputSegments := strings.Split(strings.Trim(input, "/"), "/")
	canonicalSegments := strings.Split(strings.Trim(canonical, "/"), "/")
	if len(inputSegments) != len(canonicalSegments) {
		return fmt.Errorf("expected the ID to be in the format %q", canonical)
	}

	for i, segment := range inputSegments {
		if segment == canonicalSegments[i] {
			continue
		}

		if canonicalSegments[i] == "resourceGroups" && segment == "resourcegroups" {
			continue
		}

		if i > 0 && canonicalSegments[i-1] == "providers" && strings.EqualFold(segment, canonicalSegments[i]) {
			continue
		}

		return fmt.Errorf("the segment %q should be %q", segment, canonicalSegments[i])
	}

	return nil
}
