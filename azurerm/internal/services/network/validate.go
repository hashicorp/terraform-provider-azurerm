package network

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func ValidatePrivateLinkEndpointSettings(d *schema.ResourceData) error {
	privateServiceConnections := d.Get("private_service_connection").([]interface{})

	for _, psc := range privateServiceConnections {
		privateServiceConnection := psc.(map[string]interface{})
		name := privateServiceConnection["name"].(string)

		// If this is not a manual connection and the message is set return an error since this does not make sense.
		if !privateServiceConnection["is_manual_connection"].(bool) && privateServiceConnection["request_message"].(string) != "" {
			return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute cannot be set if the "is_manual_connection" attribute is "false"`, name)
		}
		// If this is a manual connection and the message isn't set return an error.
		if privateServiceConnection["is_manual_connection"].(bool) && strings.TrimSpace(privateServiceConnection["request_message"].(string)) == "" {
			return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute must not be empty`, name)
		}
	}

	return nil
}

func ValidatePrivateLinkEndpointName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// Message="Resource name foo^* is invalid. The name can be up to 80 characters long. It must begin with a
	// word character, and it must end with a word character or with '_'. The name may contain word characters
	// or '.', '-', '_'."
	if len(v) == 1 {
		if m, _ := validate.RegExHelper(i, k, `^([a-zA-Z\d])`); !m {
			errors = append(errors, fmt.Errorf("%s must begin with a letter or number", k))
		}
	} else {
		if m, _ := validate.RegExHelper(i, k, `^([a-zA-Z\d])([a-zA-Z\d-\_\.]{0,78})([a-zA-Z\d\_])$`); !m {
			errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, periods, hyphens or underscores", k))
		}
	}

	return nil, errors
}

func ValidateVirtualHubName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^.{1,256}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 256 characters in length.", k))
	}

	return warnings, errors
}
