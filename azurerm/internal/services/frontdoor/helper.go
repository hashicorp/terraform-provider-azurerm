package frontdoor

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
)

func VerifyBackendPoolExists(backendPoolName string, backendPools []interface{}) error {
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

func AzureKeyVaultCertificateHasValues(customHttpsConfiguration map[string]interface{}, MatchAllKeys bool) bool {
	certificateSecretName := customHttpsConfiguration["azure_key_vault_certificate_secret_name"]
	certificateSecretVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"]
	certificateVaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"]

	if MatchAllKeys {
		if strings.TrimSpace(certificateSecretName.(string)) != "" && strings.TrimSpace(certificateSecretVersion.(string)) != "" && strings.TrimSpace(certificateVaultId.(string)) != "" {
			return true
		}
	} else {
		if strings.TrimSpace(certificateSecretName.(string)) != "" || strings.TrimSpace(certificateSecretVersion.(string)) != "" || strings.TrimSpace(certificateVaultId.(string)) != "" {
			return true
		}
	}

	return false
}

func IsFrontDoorFrontendEndpointConfigurable(currentState frontdoor.CustomHTTPSProvisioningState, customHttpsProvisioningEnabled bool, frontendEndpointName string, resourceGroup string) error {
	action := "disable"
	if customHttpsProvisioningEnabled {
		action = "enable"
	}

	switch currentState {
	case frontdoor.CustomHTTPSProvisioningStateDisabling, frontdoor.CustomHTTPSProvisioningStateEnabling, frontdoor.CustomHTTPSProvisioningStateFailed:
		return fmt.Errorf("Unable to %s the Front Door Frontend Endpoint %q (Resource Group %q) Custom Domain HTTPS state because the Frontend Endpoint is currently in the %q state", action, frontendEndpointName, resourceGroup, currentState)
	default:
		return nil
	}
}

func NormalizeCustomHTTPSProvisioningStateToBool(provisioningState frontdoor.CustomHTTPSProvisioningState) bool {
	isEnabled := false
	if provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled || provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabling {
		isEnabled = true
	}

	return isEnabled
}

func GetFrontDoorBasicRouteConfigurationType(i interface{}) string {
	_, ok := i.(frontdoor.ForwardingConfiguration)
	if !ok {
		_, ok := i.(frontdoor.RedirectConfiguration)
		if !ok {
			return ""
		}
		return "RedirectConfiguration"
	} else {
		return "ForwardingConfiguration"
	}
}

func VerifyRoutingRuleFrontendEndpoints(routingRuleFrontends []interface{}, configFrontendEndpoints []interface{}) error {
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

func VerifyLoadBalancingAndHealthProbeSettings(backendPools []interface{}, loadBalancingSettings []interface{}, healthProbeSettings []interface{}) error {
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

func VerifyCustomHttpsConfiguration(configFrontendEndpoints []interface{}) error {
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
					if !AzureKeyVaultCertificateHasValues(customHttpsConfiguration, true) {
						return fmt.Errorf(`"frontend_endpoint":%q "custom_https_configuration" is invalid, all of the following keys must have values in the "custom_https_configuration" block: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`, FrontendName)
					}
				} else {
					if AzureKeyVaultCertificateHasValues(customHttpsConfiguration, false) {
						return fmt.Errorf(`"frontend_endpoint":%q "custom_https_configuration" is invalid, all of the following keys must be removed from the "custom_https_configuration" block: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`, FrontendName)
					}
				}
			} else if customHttpsEnabled {
				return fmt.Errorf(`"frontend_endpoint":%q configuration is invalid because "custom_https_provisioning_enabled" is set to "true" and the "custom_https_configuration" block is undefined. please add the "custom_https_configuration" block to the configuration file`, FrontendName)
			}
		}
	}

	return nil
}

func FlattenTransformSlice(input *[]frontdoor.TransformType) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func FlattenFrontendEndpointLinkSlice(input *[]frontdoor.FrontendEndpointLink) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, *item.ID)
		}
	}
	return result
}
