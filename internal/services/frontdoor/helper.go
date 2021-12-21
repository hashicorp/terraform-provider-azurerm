package frontdoor

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/sdk/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/sdk/2020-05-01/frontdoors"
)

func isFrontDoorFrontendEndpointConfigurable(currentState frontdoors.CustomHttpsProvisioningState, customHttpsProvisioningEnabled bool, frontendEndpointId frontdoors.FrontendEndpointId) error {
	action := "disable"
	if customHttpsProvisioningEnabled {
		action = "enable"
	}

	switch currentState {
	case frontdoors.CustomHttpsProvisioningStateDisabling, frontdoors.CustomHttpsProvisioningStateEnabling, frontdoors.CustomHttpsProvisioningStateFailed:
		return fmt.Errorf("unable to %s %s Custom Domain HTTPS state because the Frontend Endpoint is currently in the %q state", action, frontendEndpointId, currentState)
	default:
		return nil
	}
}

func NormalizeCustomHTTPSProvisioningStateToBool(provisioningState frontdoors.CustomHttpsProvisioningState) bool {
	return provisioningState == frontdoors.CustomHttpsProvisioningStateEnabled || provisioningState == frontdoors.CustomHttpsProvisioningStateEnabling
}

func FlattenTransformSlice(input *[]webapplicationfirewallpolicies.TransformType) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func FlattenFrontendEndpointLinkSlice(input *[]webapplicationfirewallpolicies.FrontendEndpointLink) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			if item.Id == nil {
				continue
			}

			result = append(result, *item.Id)
		}
	}
	return result
}
