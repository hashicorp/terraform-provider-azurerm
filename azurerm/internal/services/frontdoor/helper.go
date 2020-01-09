package frontdoor

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2019-04-01/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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

// ParseAzureResourceIDLowerPath converts a long-form Azure Resource Manager ID
// into a ResourceID. We make assumptions about the structure of URLs,
// which is obviously not good, but the best thing available given the
// SDK. I had to normalize the key casing of Path because the Front Door API
// via Portal does not have consistent casing within the resource, for example:
//
// In the backendPools block the casing of the HealthProbeSettings is (notice the lowercase 'h'):
// portal-front-door/ -> healthProbeSettings/healthProbeSettings-1571100669337
//
// but in the HealthProbeSettings block the casing of the HealthProbeSettings is (notice the uppercase 'H')::
// portal-front-door/ -> HealthProbeSettings/healthProbeSettings-1571100669337
//
// so if I need to parse the name of the resource from its ID string I would be
// unable to do so with the current implementation so I normalize the key into
// a known format so I can reliable parse the ID string.
//
// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
func ParseAzureResourceIDLowerPath(id string) (*azure.ResourceID, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	// We should have an even number of key-value pairs.
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	var subscriptionID string

	// Put the constituent key-value pairs into a map
	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := strings.ToLower(components[current])
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}

		// Catch the subscriptionID before it can be overwritten by another "subscriptions"
		// value in the ID which is the case for the Service Bus subscription resource
		if key == "subscriptions" && subscriptionID == "" {
			subscriptionID = value
		} else {
			componentMap[key] = value
		}
	}

	// Build up a ResourceID from the map
	idObj := &azure.ResourceID{}
	idObj.Path = componentMap

	if subscriptionID != "" {
		idObj.SubscriptionID = subscriptionID
	} else {
		return nil, fmt.Errorf("No subscription ID found in: %q", path)
	}

	if resourceGroup, ok := componentMap["resourceGroups"]; ok {
		idObj.ResourceGroup = resourceGroup
		delete(componentMap, "resourceGroups")
	} else {
		// Some Azure APIs are weird and provide things in lower case...
		// However it's not clear whether the casing of other elements in the URI
		// matter, so we explicitly look for that case here.
		if resourceGroup, ok := componentMap["resourcegroups"]; ok {
			idObj.ResourceGroup = resourceGroup
			delete(componentMap, "resourcegroups")
		}
	}

	// It is OK not to have a provider in the case of a resource group
	if provider, ok := componentMap["providers"]; ok {
		idObj.Provider = provider
		delete(componentMap, "providers")
	}

	return idObj, nil
}
