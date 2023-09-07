package validate

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ID(i interface{}, k string) (warnings []string, errors []error) {
	if warnings, errors = validation.StringIsNotEmpty(i, k); len(errors) > 0 {
		return warnings, errors
	}

	id, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return warnings, errors
	}

	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		errors = append(errors, fmt.Errorf("cannot parse azure iot central organization ID as URI: %s", err))
		return warnings, errors
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 3 {
		errors = append(errors, fmt.Errorf("iot central organization should have 3 segments, found %d segment(s) in %q", len(components), id))
		return warnings, errors
	}

	apiString := components[0]
	if apiString != "api" {
		errors = append(errors, fmt.Errorf("iot central organization should have api as first segment, found %q", apiString))
		return warnings, errors
	}

	organizationsString := components[1]
	if organizationsString != "organizations" {
		errors = append(errors, fmt.Errorf("iot central organization should have organizations as second segment, found %q", organizationsString))
		return warnings, errors
	}

	organizationIdString := components[2]
	err = validateOrganizationId(organizationIdString)
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func OrganizationID(i interface{}, k string) (warnings []string, errors []error) {
	id, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return warnings, errors
	}

	err := validateOrganizationId(id)
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func validateOrganizationId(id string) error {
	// Ensure the string follows the desired format.
	// Regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$
	// The negative lookahead (?!-) is not supported in Go's standard regexp package
	formatPattern := `^[a-z0-9-]{1,48}[a-z0-9]$`
	formatRegex, err := regexp.Compile(formatPattern)
	if err != nil {
		return fmt.Errorf("error compiling format regex: %s error: %+v", formatPattern, err)
	}

	if !formatRegex.MatchString(id) {
		return fmt.Errorf("iot central organizationId %q is invalid, regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$", id)
	}

	// Ensure the string does not start with a hyphen.
	// Solves for (?!-)
	startHyphenPattern := `^-`
	startHyphenRegex, err := regexp.Compile(startHyphenPattern)
	if err != nil {
		return fmt.Errorf("error compiling start hyphen regex: %s error: %+v", startHyphenPattern, err)
	}

	if startHyphenRegex.MatchString(id) {
		return fmt.Errorf("iot central organizationId %q is invalid, regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$", id)
	}

	return nil
}
