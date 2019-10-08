package network

import (
	"fmt"
	"net"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ValidatePrivateLinkServiceName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z]){1,}([\d\._\-a-zA-Z]{0,77})([\da-zA-Z_]$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must be between 1 and 80 characters, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.`, k))
	}

	return nil, errors
}

func ValidatePrivateLinkServiceSubsciptionGuid(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^([0-9A-Fa-f]{8}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{4}[-][0-9A-Fa-f]{12})$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q is an invalid subscription GUID.`, k))
	}

	return nil, errors
}

func ValidatePrivateLinkServiceSubsciptionFqdn(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^(([a-zA-Z\d]|[a-zA-Z\d][a-zA-Z\d\-]*[a-zA-Z\d])\.){1,}([a-zA-Z\d]|[a-zA-Z\d][a-zA-Z\d\-]*[a-zA-Z\d\.]){1,}$`); !m {
		errors = append(errors, fmt.Errorf(`%q is an invalid FQDN`, v))
	}

	// I use 255 here because the string contains the upto three . characters in them
	if len(v) > 255 {
		errors = append(errors, fmt.Errorf(`FQDNs can not be longer than 255 characters in length, got %d characters`, len(v)))
	}

	segments := utils.SplitRemoveEmptyEntries(v, ".", false)
	index := 0

	for _, label := range segments {
		index++
		if index == len(segments) {
			if len(label) < 2 {
				errors = append(errors, fmt.Errorf(`the last label of an FQDN must be at least 2 characters, got 1 character`))
			}
		} else {
			if len(label) > 63 {
				errors = append(errors, fmt.Errorf(`FQDN labels must not be longer than 63 characters, got %d characters`, len(label)))
			}
		}
	}

	return nil, errors
}

func ValidatePrivateLinkServiceIPv4Address(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	result := net.ParseIP(v)
	if result.To4() == nil {
		errors = append(errors, fmt.Errorf("%q is not a valid IPv4 IP address", v))
	}

	return nil, errors
}

func ValidatePrivateLinkServiceIsResourceId(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if _, err := azure.ParseAzureResourceID(v); err != nil {
		return nil, []error{fmt.Errorf("%q is not a valid %q", v, k)}
	}

	return nil, errors
}
