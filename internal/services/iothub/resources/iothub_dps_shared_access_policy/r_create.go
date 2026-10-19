// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_shared_access_policy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSSharedAccessPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubDpsId := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_dps_name").(string))

	locks.ByName(iothubDpsId.ProvisioningServiceName, IothubResourceName)
	defer locks.UnlockByName(iothubDpsId.ProvisioningServiceName, IothubResourceName)

	iothubDps, err := client.Get(ctx, iothubDpsId)
	if err != nil {
		if response.WasNotFound(iothubDps.HttpResponse) {
			return fmt.Errorf("IotHub DPS %s was not found", iothubDpsId.String())
		}

		return fmt.Errorf("retrieving IotHub DPS %s: %+v", iothubDpsId.String(), err)
	}

	if iothubDps.Model == nil {
		return fmt.Errorf("retrieving IotHub DPS %s: ID was nil", iothubDpsId.String())
	}

	id := iotdpsresource.NewKeyID(iothubDpsId.SubscriptionId, iothubDpsId.ResourceGroupName, iothubDpsId.ProvisioningServiceName, d.Get("name").(string))

	accessRights := dpsAccessRights{
		enrollmentRead:    d.Get("enrollment_read").(bool),
		enrollmentWrite:   d.Get("enrollment_write").(bool),
		registrationRead:  d.Get("registration_read").(bool),
		registrationWrite: d.Get("registration_write").(bool),
		serviceConfig:     d.Get("service_config").(bool),
	}

	if err := accessRights.validate(); err != nil {
		return fmt.Errorf("building Access Rights: %s", err)
	}

	expandedAccessPolicy := iotdpsresource.SharedAccessSignatureAuthorizationRuleAccessRightsDescription{
		KeyName: id.KeyName,
		Rights:  iotdpsresource.AccessRightsDescription(expandDpsAccessRights(accessRights)),
	}

	accessPolicies := make([]iotdpsresource.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	alreadyExists := false
	existingAccessPolicies, err := client.ListKeysComplete(ctx, iothubDpsId)
	if err != nil {
		return fmt.Errorf("loading %s: %+v", id, err)
	}
	for _, existingAccessPolicy := range existingAccessPolicies.Items {
		if strings.EqualFold(existingAccessPolicy.KeyName, id.KeyName) {
			if d.IsNewResource() {
				if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return tf.ImportAsExistsError("azurerm_iothub_dps_shared_access_policy", id.ID())
				}
			}
			accessPolicies = append(accessPolicies, expandedAccessPolicy)
			alreadyExists = true
		} else {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	if d.IsNewResource() {
		accessPolicies = append(accessPolicies, expandedAccessPolicy)
	} else if !alreadyExists {
		return fmt.Errorf("unable to find %s", id)
	}

	iothubDps.Model.Properties.AuthorizationPolicies = &accessPolicies

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, iothubDpsId, *iothubDps.Model, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("updating IotHub DPS %s with Shared Access Policy %s: %+v", iothubDpsId, id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, iothubDpsId, *iothubDps.Model); err != nil {
			return fmt.Errorf("updating IotHub DPS %s with Shared Access Policy %s: %+v", iothubDpsId, id, err)
		}
	}

	return resourceIotHubDPSSharedAccessPolicyRead(d, meta)
}

func expandDpsAccessRights(input dpsAccessRights) string {
	actualRights := make([]string, 0)

	// NOTE: the iteration order is important here
	if input.enrollmentRead {
		actualRights = append(actualRights, "EnrollmentRead")
	}

	if input.enrollmentWrite {
		actualRights = append(actualRights, "EnrollmentWrite")
	}

	if input.registrationRead {
		actualRights = append(actualRights, "RegistrationStatusRead")
	}

	if input.registrationWrite {
		actualRights = append(actualRights, "RegistrationStatusWrite")
	}

	if input.serviceConfig {
		actualRights = append(actualRights, "ServiceConfig")
	}

	return strings.Join(actualRights, ", ")
}
