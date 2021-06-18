package frontdoor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func customizeHttpsConfigurationCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
	if v, ok := d.GetOk("frontend_endpoint_id"); ok && v.(string) != "" {
		id, err := parse.FrontendEndpointID(v.(string))
		if err != nil {
			return err
		}

		if err := customHttpsSettings(d); err != nil {
			return fmt.Errorf("validating Front Door Custom Https Configuration for Endpoint %q (Front Door %q / Resource Group %q): %+v", id.Name, id.FrontDoorName, id.ResourceGroup, err)
		}
	}

	return nil
}

func customHttpsSettings(d *pluginsdk.ResourceDiff) error {
	frontendEndpointCustomHttpsConfig := d.Get("custom_https_configuration").([]interface{})
	customHttpsEnabled := d.Get("custom_https_provisioning_enabled").(bool)

	if len(frontendEndpointCustomHttpsConfig) > 0 {
		if !customHttpsEnabled {
			return fmt.Errorf(`"custom_https_provisioning_enabled" is set to "false". please remove the "custom_https_configuration" block from the configuration file`)
		}

		// Verify frontend endpoints custom https configuration is valid if defined
		if err := verifyCustomHttpsConfiguration(frontendEndpointCustomHttpsConfig); err != nil {
			return err
		}
	} else if customHttpsEnabled {
		return fmt.Errorf(`"custom_https_provisioning_enabled" is set to "true". please add a "custom_https_configuration" block to the configuration file`)
	}

	return nil
}

func verifyCustomHttpsConfiguration(frontendEndpointCustomHttpsConfig []interface{}) error {
	if len(frontendEndpointCustomHttpsConfig) > 0 {
		customHttpsConfiguration := frontendEndpointCustomHttpsConfig[0].(map[string]interface{})
		certificateSource := customHttpsConfiguration["certificate_source"].(string)
		certificateVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)

		if certificateSource == string(frontdoor.CertificateSourceFrontDoor) {
			if azureKeyVaultCertificateHasValues(customHttpsConfiguration, true) {
				return fmt.Errorf(`a Front Door managed "custom_https_configuration" block does not support the following keys. Please remove the following keys from your configuration file: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`)
			}
		} else {
			// The latest secret version is no longer valid for key vaults
			if strings.EqualFold(certificateVersion, "latest") {
				return fmt.Errorf(`"azure_key_vault_certificate_secret_version" can not be set to "latest" please remove this attribute from the configuration file. Removing the value has the same functionality as setting it to "latest"`)
			}

			if !azureKeyVaultCertificateHasValues(customHttpsConfiguration, false) {
				if certificateVersion == "" {
					// If using latest, empty string is now equivalent to using the keyword latest
					return fmt.Errorf(`a "AzureKeyVault" managed "custom_https_configuration" block must have values in the following fileds: "azure_key_vault_certificate_secret_name" and "azure_key_vault_certificate_vault_id"`)
				} else {
					// If using a specific version of the secret
					return fmt.Errorf(`a "AzureKeyVault" managed "custom_https_configuration" block must have values in the following fileds: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`)
				}
			}
		}
	}

	return nil
}

func azureKeyVaultCertificateHasValues(customHttpsConfiguration map[string]interface{}, isFrontDoorManaged bool) bool {
	certificateSecretName := customHttpsConfiguration["azure_key_vault_certificate_secret_name"].(string)
	certificateSecretVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)
	certificateVaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"].(string)

	if isFrontDoorManaged {
		// if any of these keys have values it is invalid
		if strings.TrimSpace(certificateSecretName) != "" || strings.TrimSpace(certificateSecretVersion) != "" || strings.TrimSpace(certificateVaultId) != "" {
			return true
		}
	} else {
		if certificateSecretVersion == "" {
			// using latest ignore certificate secret version
			if strings.TrimSpace(certificateSecretName) != "" && strings.TrimSpace(certificateVaultId) != "" {
				return true
			}
		} else {
			// not using latest make sure all keys have values
			if strings.TrimSpace(certificateSecretName) != "" && strings.TrimSpace(certificateSecretVersion) != "" && strings.TrimSpace(certificateVaultId) != "" {
				return true
			}
		}
	}

	return false
}

func frontDoorCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
	if err := frontDoorSettings(d); err != nil {
		return fmt.Errorf("validating Front Door %q (Resource Group %q): %+v", d.Get("name").(string), d.Get("resource_group_name").(string), err)
	}

	return nil
}

func frontDoorSettings(d *pluginsdk.ResourceDiff) error {
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
