package helper

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
)

func DoesBackendPoolExists(backendPoolName string, backendPools []interface{}) error {
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

func ConvertToPolicyEnabledStateFromBool(isEnabled bool) frontdoor.PolicyEnabledState {
	if isEnabled {
		return frontdoor.PolicyEnabledStateEnabled 
	}

	return frontdoor.PolicyEnabledStateDisabled 
}

func ConvertBoolToCustomRuleEnabledState(isEnabled bool) frontdoor.CustomRuleEnabledState {
	if isEnabled {
		return frontdoor.CustomRuleEnabledStateEnabled 
	}

	return frontdoor.CustomRuleEnabledStateDisabled 
}

func ConvertToPolicyModeFromString(policyMode string) frontdoor.PolicyMode {
	if policyMode == "Detection" {
		return frontdoor.Detection  
	}

	return frontdoor.Prevention  
}

func ConvertStringToRuleType(ruleType string) frontdoor.RuleType  {
	if ruleType == string(frontdoor.MatchRule) {
		return frontdoor.MatchRule  
	}

	return frontdoor.RateLimitRule  
}

func ConvertStringToActionType(actionType string) frontdoor.ActionType {
	switch actionType {
		case string(frontdoor.Allow):
			return frontdoor.Allow
		case string(frontdoor.Block):
			return frontdoor.Block
		case string(frontdoor.Log):
			return frontdoor.Log
		case string(frontdoor.Redirect):
			return frontdoor.Redirect
	}
}

func GetFrontDoorSubResourceId(subscriptionId string, resourceGroup string, frontDoorName string, resourceType string, resourceName string) string {
	if strings.TrimSpace(subscriptionId) == "" || strings.TrimSpace(resourceGroup) == "" || strings.TrimSpace(frontDoorName) == "" || strings.TrimSpace(resourceType) == "" || strings.TrimSpace(resourceName) == "" {
		return ""
	}

	return fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.Network/Frontdoors/%s/%s/%s", subscriptionId, resourceGroup, frontDoorName, resourceType, resourceName)
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
