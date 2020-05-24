package validate

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func FrontDoorName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{3,61})([\da-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 5 and 63 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}
	return nil, nil
}

func FrontDoorWAFName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{0,127})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 128 characters in length, must begin with a letter and may only contain letters and numbers.`, k))
	}

	return nil, nil
}

func FrontDoorBackendPoolRoutingRuleName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z])([-\da-zA-Z]{1,88})([\da-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`, k))
	}

	return nil, nil
}

func FrontdoorCustomBlockResponseBody(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q contains invalid characters, %q must contain only alphanumeric and equals sign characters.`, k, k))
	}

	return nil, nil
}

func FrontdoorSettings(d *schema.ResourceDiff) error {
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
			return fmt.Errorf(`routing_rule %s block is invalid. you must have either a "redirect_configuration" or a "forwarding_configuration" defined for the routing_rule %s`, routingRuleName, routingRuleName)
		}

		// Check 1. validate that only one configuration type is defined per routing rule
		if len(redirectConfig) == 1 && len(forwardConfig) == 1 {
			return fmt.Errorf(`routing_rule %s block is invalid. "redirect_configuration" conflicts with "forwarding_configuration". You can only have one configuration type per each routing rule`, routingRuleName)
		}

		// Check 2. routing rule is a forwarding_configuration type make sure the backend_pool_name exists in the configuration file
		if len(forwardConfig) > 0 {
			fc := forwardConfig[0].(map[string]interface{})

			if err := verifyBackendPoolExists(fc["backend_pool_name"].(string), backendPools); err != nil {
				return fmt.Errorf(`routing_rule %s is invalid. %+v`, routingRuleName, err)
			}
		}

		// Check 3. validate that each routing rule frontend_endpoints are actually defined in the resource schema
		if routingRuleFrontends := routingRule["frontend_endpoints"].([]interface{}); len(routingRuleFrontends) > 0 {
			if err := verifyRoutingRuleFrontendEndpoints(routingRuleFrontends, configFrontendEndpoints); err != nil {
				return fmt.Errorf(`"routing_rule":%q %+v`, routingRuleName, err)
			}
		} else {
			return fmt.Errorf(`"routing_rule": %q must have at least one "frontend_endpoints" defined`, routingRuleName)
		}
	}

	// Verify backend pool load balancing settings and health probe settings are defined in the resource schema
	if err := verifyLoadBalancingAndHealthProbeSettings(backendPools, loadBalancingSettings, healthProbeSettings); err != nil {
		return fmt.Errorf(`%+v`, err)
	}

	// Verify frontend endpoints custom https configuration is valid if defined
	if err := verifyCustomHttpsConfiguration(configFrontendEndpoints); err != nil {
		return fmt.Errorf(`%+v`, err)
	}

	return nil
}

func verifyBackendPoolExists(backendPoolName string, backendPools []interface{}) error {
	if backendPoolName == "" {
		return fmt.Errorf(`"backend_pool_name" cannot be empty`)
	}

	for _, bps := range backendPools {
		backendPool := bps.(map[string]interface{})
		if backendPool["name"].(string) == backendPoolName {
			return nil
		}
	}

	return fmt.Errorf(`unable to locate "backend_pool_name":%q in configuration file`, backendPoolName)
}

func verifyRoutingRuleFrontendEndpoints(routingRuleFrontends []interface{}, configFrontendEndpoints []interface{}) error {
	for _, routingRuleFrontend := range routingRuleFrontends {
		// Get the name of the frontend defined in the routing rule
		routingRulefrontendName := routingRuleFrontend.(string)
		found := false

		// Loop over all of the defined frontend endpoints in the config
		// seeing if we find the routing rule frontend in the list
		for _, configFrontendEndpoint := range configFrontendEndpoints {
			configFrontend := configFrontendEndpoint.(map[string]interface{})
			configFrontendName := configFrontend["name"]
			if routingRulefrontendName == configFrontendName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf(`"frontend_endpoints":%q was not found in the configuration file. verify you have the "frontend_endpoint":%q defined in the configuration file`, routingRulefrontendName, routingRulefrontendName)
		}
	}

	return nil
}

func verifyLoadBalancingAndHealthProbeSettings(backendPools []interface{}, loadBalancingSettings []interface{}, healthProbeSettings []interface{}) error {
	for _, bps := range backendPools {
		backendPool := bps.(map[string]interface{})
		backendPoolName := backendPool["name"]
		backendPoolLoadBalancingName := backendPool["load_balancing_name"]
		backendPoolHealthProbeName := backendPool["health_probe_name"]
		found := false

		// Verify backend pool load balancing settings name exists
		if len(loadBalancingSettings) > 0 {
			for _, lbs := range loadBalancingSettings {
				loadBalancing := lbs.(map[string]interface{})
				loadBalancingName := loadBalancing["name"]

				if loadBalancingName == backendPoolLoadBalancingName {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf(`"backend_pool":%q "load_balancing_name":%q was not found in the configuration file. verify you have the "backend_pool_load_balancing":%q defined in the configuration file`, backendPoolName, backendPoolLoadBalancingName, backendPoolLoadBalancingName)
			}
		}

		found = false

		// Verify health probe settings name exists
		if len(healthProbeSettings) > 0 {
			for _, hps := range healthProbeSettings {
				healthProbe := hps.(map[string]interface{})
				healthProbeName := healthProbe["name"]

				if healthProbeName == backendPoolHealthProbeName {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf(`"backend_pool":%q "health_probe_name":%q was not found in the configuration file. verify you have the "backend_pool_health_probe":%q defined in the configuration file`, backendPoolName, backendPoolHealthProbeName, backendPoolHealthProbeName)
			}
		}
	}

	return nil
}

func verifyCustomHttpsConfiguration(configFrontendEndpoints []interface{}) error {
	for _, configFrontendEndpoint := range configFrontendEndpoints {
		if configFrontend := configFrontendEndpoint.(map[string]interface{}); len(configFrontend) > 0 {
			FrontendName := configFrontend["name"]
			customHttpsEnabled := configFrontend["custom_https_provisioning_enabled"].(bool)

			if chc := configFrontend["custom_https_configuration"].([]interface{}); len(chc) > 0 {
				if !customHttpsEnabled {
					return fmt.Errorf(`"frontend_endpoint":%q "custom_https_configuration" is invalid because "custom_https_provisioning_enabled" is set to "false". please remove the "custom_https_configuration" block from the configuration file`, FrontendName)
				}

				customHttpsConfiguration := chc[0].(map[string]interface{})
				certificateSource := customHttpsConfiguration["certificate_source"]
				if certificateSource == string(frontdoor.CertificateSourceAzureKeyVault) {
					if !azureKeyVaultCertificateHasValues(customHttpsConfiguration, true) {
						return fmt.Errorf(`"frontend_endpoint":%q "custom_https_configuration" is invalid, all of the following keys must have values in the "custom_https_configuration" block: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`, FrontendName)
					}
				} else if azureKeyVaultCertificateHasValues(customHttpsConfiguration, false) {
					return fmt.Errorf(`"frontend_endpoint":%q "custom_https_configuration" is invalid, all of the following keys must be removed from the "custom_https_configuration" block: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`, FrontendName)
				}
			} else if customHttpsEnabled {
				return fmt.Errorf(`"frontend_endpoint":%q configuration is invalid because "custom_https_provisioning_enabled" is set to "true" and the "custom_https_configuration" block is undefined. please add the "custom_https_configuration" block to the configuration file`, FrontendName)
			}
		}
	}

	return nil
}

func azureKeyVaultCertificateHasValues(customHttpsConfiguration map[string]interface{}, matchAllKeys bool) bool {
	certificateSecretName := customHttpsConfiguration["azure_key_vault_certificate_secret_name"]
	certificateSecretVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"]
	certificateVaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"]

	if matchAllKeys {
		if strings.TrimSpace(certificateSecretName.(string)) != "" && strings.TrimSpace(certificateSecretVersion.(string)) != "" && strings.TrimSpace(certificateVaultId.(string)) != "" {
			return true
		}
	} else if strings.TrimSpace(certificateSecretName.(string)) != "" || strings.TrimSpace(certificateSecretVersion.(string)) != "" || strings.TrimSpace(certificateVaultId.(string)) != "" {
		return true
	}

	return false
}
