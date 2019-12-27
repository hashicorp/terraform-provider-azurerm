package network

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func ValidatePrivateEndpointSettings(d *schema.ResourceData) error {
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

func ValidatePrivateLinkNatIpConfiguration(d *schema.ResourceDiff) error {
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	ipConfigurations := d.Get("nat_ip_configuration").([]interface{})

	for i, item := range ipConfigurations {
		v := item.(map[string]interface{})
		p := fmt.Sprintf("nat_ip_configuration.%d.private_ip_address", i)
		s := fmt.Sprintf("nat_ip_configuration.%d.subnet_id", i)
		isPrimary := v["primary"].(bool)
		in := v["name"].(string)

		if d.HasChange(p) {
			o, n := d.GetChange(p)
			if o != "" && n == "" {
				return fmt.Errorf("Private Link Service %q (Resource Group %q) nat_ip_configuration %q private_ip_address once assigned can not be removed", name, resourceGroup, in)
			}
		}

		if isPrimary && d.HasChange(s) {
			o, _ := d.GetChange(s)
			if o != "" {
				return fmt.Errorf("Private Link Service %q (Resource Group %q) nat_ip_configuration %q primary subnet_id once assigned can not be changed", name, resourceGroup, in)
			}
		}
	}

	return nil
}

func ValidatePrivateLinkName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules per the Nat Gateway service team are (Friday, October 18, 2019 4:20 PM):
	// 1. Must not be empty.
	// 2. Must be between 1 and 80 characters.
	// 3. The attribute must:
	//    a) begin with a letter or number
	//    b) end with a letter, number or underscore
	//    c) may contain only letters, numbers, underscores, periods, or hyphens.

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

func ValidateNatGatewayName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	// The name attribute rules per the Nat Gateway service team are (Friday, October 18, 2019 4:20 PM):
	// 1. Must not be empty.
	// 2. Must be between 1 and 80 characters.
	// 3. The attribute must:
	//    a) begin with a letter or number
	//    b) end with a letter, number or underscore
	//    c) may contain only letters, numbers, underscores, periods, or hyphens.

	if len(v) == 1 {
		if matched := regexp.MustCompile(`^([a-zA-Z\d])`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%s must begin with a letter or number", k))
		}
	} else {
		if matched := regexp.MustCompile(`^([a-zA-Z\d])([a-zA-Z\d-\_\.]{0,78})([a-zA-Z\d\_])$`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens", k))
		}
	}

	return warnings, errors
}

func ValidatePrivateLinkSubResourceName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if len(strings.TrimSpace(v)) >= 3 {
		if m, _ := validate.RegExHelper(i, k, `^([a-zA-Z0-9])([\w\.-]{1,61})([a-zA-Z0-9])$`); !m {
			errors = append(errors, fmt.Errorf("%s must begin and end with a alphanumeric character, be between 3 and 63 characters in length, only contain letters, numbers, underscores, periods, and dashes", k))
		}
	} else {
		errors = append(errors, fmt.Errorf("%s must be at least 3 character in length", k))
	}

	return nil, errors
}
