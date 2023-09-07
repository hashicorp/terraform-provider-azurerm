package validate

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func OrganizationID(i interface{}, k string) (warnings []string, errors []error) {
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
		errors = append(errors, fmt.Errorf("cannot parse azure iot central organization ID: %s", err))
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 3 {
		errors = append(errors, fmt.Errorf("iot central organization should have 3 segments, found %d segment(s) in %q", len(components), id))
	}

	apiString := components[0]
	if apiString != "api" {
		errors = append(errors, fmt.Errorf("iot central organization should have api as first segment, found %q", apiString))
	}

	organizationsString := components[1]
	if organizationsString != "organizations" {
		errors = append(errors, fmt.Errorf("iot central organization should have organizations as second segment, found %q", organizationsString))
	}

	return warnings, errors
}
