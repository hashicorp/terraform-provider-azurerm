package frontdoor

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

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

func ValidateFrontdoorSettings(d *schema.ResourceDiff) error {
	routingRules := d.Get("routing_rule").([]interface{})
	configFrontendEndpoints := d.Get("frontend_endpoint").([]interface{})
	backendPools := d.Get("backend_pool").([]interface{})
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]interface{})
	healthProbeSettings := d.Get("backend_pool_health_probe").([]interface{})

	if len(configFrontendEndpoints) == 0 {
		return fmt.Errorf(`"frontend_endpoint": must have at least one "frontend_endpoint" defined, found 0`)
	}

	// Loop over all of the Routing Rules and validate that only one type of configuration is defined per Routing Rule
	for _, rr := range routingRules {
		routingRule := rr.(map[string]interface{})
		routingRuleName := routingRule["name"]
		redirectConfig := routingRule["redirect_configuration"].([]interface{})
		forwardConfig := routingRule["forwarding_configuration"].([]interface{})

		// Check 0. validate that at least one routing configuration exists per routing rule
		if len(redirectConfig) == 0 && len(forwardConfig) == 0 {
			return fmt.Errorf(`"routing_rule":%q is invalid. you must have either a "redirect_configuration" or a "forwarding_configuration" defined for the "routing_rule":%q `, routingRuleName, routingRuleName)
		}

		// Check 1. validate that only one configuration type is defined per routing rule
		if len(redirectConfig) == 1 && len(forwardConfig) == 1 {
			return fmt.Errorf(`"routing_rule":%q is invalid. "redirect_configuration" conflicts with "forwarding_configuration". You can only have one configuration type per each routing rule`, routingRuleName)
		}

		// Check 2. routing rule is a forwarding_configuration type make sure the backend_pool_name exists in the configuration file
		if len(forwardConfig) > 0 {
			fc := forwardConfig[0].(map[string]interface{})
			if err := VerifyBackendPoolExists(fc["backend_pool_name"].(string), backendPools); err != nil {
				return fmt.Errorf(`"routing_rule":%q is invalid. %+v`, routingRuleName, err)
			}
		}

		// Check 3. validate that each routing rule frontend_endpoints are actually defined in the resource schema
		if routingRuleFrontends := routingRule["frontend_endpoints"].([]interface{}); len(routingRuleFrontends) > 0 {
			if err := VerifyRoutingRuleFrontendEndpoints(routingRuleFrontends, configFrontendEndpoints); err != nil {
				return fmt.Errorf(`"routing_rule":%q %+v`, routingRuleName, err)
			}
		} else {
			return fmt.Errorf(`"routing_rule": %q must have at least one "frontend_endpoints" defined`, routingRuleName)
		}
	}

	// Verify backend pool load balancing settings and health probe settings are defined in the resource schema
	if err := VerifyLoadBalancingAndHealthProbeSettings(backendPools, loadBalancingSettings, healthProbeSettings); err != nil {
		return fmt.Errorf(`%+v`, err)
	}

	// Verify frontend endpoints custom https configuration is valid if defined
	if err := VerifyCustomHttpsConfiguration(configFrontendEndpoints); err != nil {
		return fmt.Errorf(`%+v`, err)
	}

	return nil
}
