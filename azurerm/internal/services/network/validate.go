package network

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
		errors = append(errors, fmt.Errorf(`%q is an invalid fqdn.`, k))
	}

	if len(v) > 253 {
		errors = append(errors, fmt.Errorf(`FQDNs can not be longer than 253 characters in length, got %d characters.`, len(v)))
	}

	// TODO: Remove empty entries, this is a bug with the trailing . format
	segments := strings.Split(v, ".")
	index := 0

	for _,label := range segments {
		index++
		fmt.Println(label)
		
		if index == len(segments) {
			if label != "" && len(label) < 2 {
				errors = append(errors,fmt.Errorf(`the last label of an FQDN must be at least 2 characters long, %q is only 1 character in length.`, label))
			}
		} else {
			if len(label) > 63 {
				errors = append(errors,fmt.Errorf(`labels of an FQDN must not be longer than 63 characters, got %d character.`, len(label)))
			}
		}
	}


	return nil, errors
}