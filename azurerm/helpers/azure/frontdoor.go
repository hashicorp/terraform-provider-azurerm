package azure

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

//Frontdoor name must begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.
func ValidateFrontDoorName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{3,61})([\da-zA-Z]$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must be between 5 and 63 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}

	return nil, errors
}

func ValidateBackendPoolRoutingRuleName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{1,88})([\da-zA-Z]$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}

	return nil, errors
}

func ValidateFrontdoor(d *schema.ResourceData) error {
	routingRules := d.Get("routing_rule").([]interface{})
	configFrontendEndpoints := d.Get("frontend_endpoint").([]interface{})

	// if len(routingRules) == 0 {
	// 	return nil
	// }

	if len(configFrontendEndpoints) == 0 {
		return fmt.Errorf(`"frontend_endpoint": must have at least one "frontend_endpoint" defined, found 0`)
	}

	// Loop over all of the Routing Rules and validate that only one type of configuration is defined per Routing Rule
	// start routing rule validation
	if routingRules != nil {
		for _, rr := range routingRules {
			routingRule := rr.(map[string]interface{})
			routingRuleName := routingRule["name"]
			found := false

			redirectConfig := routingRule["redirect_configuration"].([]interface{})
			forwardConfig := routingRule["forwarding_configuration"].([]interface{})

			// 1. validate that only one type of redirect configuration is set per routing rule
			if len(redirectConfig) == 1 && len(forwardConfig) == 1 {
				return fmt.Errorf(`"routing_rule":%q is invalid. "redirect_configuration" conflicts with "forwarding_configuration". You can only have one configuration type per routing rule`, routingRuleName)
			}

			// 2. validate that each routing rule frontend_endpoints are actually defined in the resource schema
			if rrFrontends := routingRule["frontend_endpoints"].([]interface{}); len(rrFrontends) > 0 {

				for _, rrFrontend := range rrFrontends {
					// Get the name of the frontend defined in the routing rule
					frontendsName := rrFrontend.(string)
					found = false

					// Loop over all of the defined frontend endpoints in the config 
					// seeing if we find the routing rule frontend in the list
					for _, configFrontendEndpoint := range configFrontendEndpoints {
						cFrontend := configFrontendEndpoint.(map[string]interface{})
						configFrontendName := cFrontend["name"]
						if( frontendsName == configFrontendName){
							found = true
							break
						}
					}

					if !found {
						return fmt.Errorf(`"routing_rule":%q "frontend_endpoints":%q was not found in the configuration file. verify you have the "frontend_endpoint":%q defined in the configuration file`, routingRuleName, frontendsName, frontendsName)
					}
				}
			} else {
				return fmt.Errorf(`"routing_rule": %q must have at least one "frontend_endpoints" defined`, routingRuleName)
			}

		} // end routing rule validation

	} // end routing rule nil check

	return nil
}