package frontdoor

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
)

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
