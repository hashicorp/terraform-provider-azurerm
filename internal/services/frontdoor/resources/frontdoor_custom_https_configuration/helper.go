// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
)

func isFrontDoorFrontendEndpointConfigurable(currentState frontdoors.CustomHTTPSProvisioningState, customHttpsProvisioningEnabled bool, frontendEndpointId frontdoors.FrontendEndpointId) error {
	action := "disable"
	if customHttpsProvisioningEnabled {
		action = "enable"
	}

	switch currentState {
	case frontdoors.CustomHTTPSProvisioningStateDisabling, frontdoors.CustomHTTPSProvisioningStateEnabling, frontdoors.CustomHTTPSProvisioningStateFailed:
		return fmt.Errorf("unable to %s %s Custom Domain HTTPS state because the Frontend Endpoint is currently in the %q state", action, frontendEndpointId, currentState)
	default:
		return nil
	}
}

func NormalizeCustomHTTPSProvisioningStateToBool(provisioningState frontdoors.CustomHTTPSProvisioningState) bool {
	return provisioningState == frontdoors.CustomHTTPSProvisioningStateEnabled || provisioningState == frontdoors.CustomHTTPSProvisioningStateEnabling
}
