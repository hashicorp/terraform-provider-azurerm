package network

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// ValidatePrivateLinkNatIpConfiguration -  This rule makes sure that you can only go from a
// dynamic private ip address to a static private ip address. Once you have assigned a private
// ip address to a primary or secondary nat ip configuration block it is set in stone and can
// not become a dynamic private ip address again unless the resource is destroyed and recreated.
func ValidatePrivateLinkNatIpConfiguration(d *schema.ResourceDiff) error {
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	primaryIpConfiguration := d.Get("primary_nat_ip_configuration").([]interface{})
	secondaryIpConfigurations := d.Get("auxillery_nat_ip_configuration").([]interface{})

	for i, item := range primaryIpConfiguration {
		v := item.(map[string]interface{})
		p := fmt.Sprintf("primary_nat_ip_configuration.%d.private_ip_address", i)
		in := v["name"].(string)

		if d.HasChange(p) {
			o, n := d.GetChange(p)
			if o != "" && n == "" {
				return fmt.Errorf("Private Link Service %q (Resource Group %q) primary_nat_ip_configuration %q private_ip_address once assigned can not be removed", name, resourceGroup, in)
			}
		}
	}

	for i, item := range secondaryIpConfigurations {
		v := item.(map[string]interface{})
		p := fmt.Sprintf("auxillery_nat_ip_configuration.%d.private_ip_address", i)
		in := v["name"].(string)

		if d.HasChange(p) {
			o, n := d.GetChange(p)
			if o != "" && n == "" {
				return fmt.Errorf("Private Link Service %q (Resource Group %q) auxillery_nat_ip_configuration %q private_ip_address once assigned can not be removed", name, resourceGroup, in)
			}
		}
	}

	return nil
}

func ValidatePrivateLinkServiceName(i interface{}, k string) (_ []string, errors []error) {
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
			errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens", k))
		}
	}

	return nil, errors
}

func ValidatePrivateLinkServiceSubsciptionFqdn(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %q to be string", k))
	}

	if m, _ := validate.RegExHelper(i, k, `^(([a-zA-Z\d]|[a-zA-Z\d][a-zA-Z\d\-]*[a-zA-Z\d])\.){1,}([a-zA-Z\d]|[a-zA-Z\d][a-zA-Z\d\-]*[a-zA-Z\d\.]){1,}$`); !m {
		errors = append(errors, fmt.Errorf(`%q is an invalid FQDN`, v))
	}

	// I use 255 here because the string contains the upto three . characters in it
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
