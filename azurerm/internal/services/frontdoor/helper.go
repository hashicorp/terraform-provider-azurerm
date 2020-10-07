package frontdoor

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
)

func isFrontDoorFrontendEndpointConfigurable(currentState frontdoor.CustomHTTPSProvisioningState, customHttpsProvisioningEnabled bool, frontendEndpointId parse.FrontendEndpointId) error {
	action := "disable"
	if customHttpsProvisioningEnabled {
		action = "enable"
	}

	switch currentState {
	case frontdoor.CustomHTTPSProvisioningStateDisabling, frontdoor.CustomHTTPSProvisioningStateEnabling, frontdoor.CustomHTTPSProvisioningStateFailed:
		return fmt.Errorf("Unable to %s the Front Door Frontend Endpoint %q (Resource Group %q) Custom Domain HTTPS state because the Frontend Endpoint is currently in the %q state", action, frontendEndpointId.Name, frontendEndpointId.ResourceGroup, currentState)
	default:
		return nil
	}
}

func NormalizeCustomHTTPSProvisioningStateToBool(provisioningState frontdoor.CustomHTTPSProvisioningState) bool {
	return provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled || provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabling
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
			if item.ID == nil {
				continue
			}

			result = append(result, *item.ID)
		}
	}
	return result
}
