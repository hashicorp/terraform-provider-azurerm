// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_shared_access_policy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSSharedAccessPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotdpsresource.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ProvisioningServiceName, IothubResourceName)
	defer locks.UnlockByName(id.ProvisioningServiceName, IothubResourceName)

	iothubDpsId := commonids.NewProvisioningServiceID(id.SubscriptionId, id.ResourceGroupName, id.ProvisioningServiceName)
	iothubDps, err := client.Get(ctx, iothubDpsId)
	if err != nil {
		if response.WasNotFound(iothubDps.HttpResponse) {
			return fmt.Errorf("IotHub DPS %q was not found", id)
		}

		return fmt.Errorf("loading IotHub DPS %q: %+v", id, err)
	}

	accessPolicies := make([]iotdpsresource.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	existingAccessPolicies, err := client.ListKeysComplete(ctx, iothubDpsId)
	if err != nil {
		return fmt.Errorf("loading %s: %+v", id, err)
	}
	for _, existingAccessPolicy := range existingAccessPolicies.Items {
		if !strings.EqualFold(existingAccessPolicy.KeyName, id.KeyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothubDps.Model.Properties.AuthorizationPolicies = &accessPolicies

	if err := client.CreateOrUpdateThenPoll(ctx, iothubDpsId, *iothubDps.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return nil
}
