// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_shared_access_policy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubSharedAccessPolicyCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSharedAccessPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	expandedAccessPolicy := devices.SharedAccessSignatureAuthorizationRule{
		KeyName: &id.IotHubKeyName,
		Rights:  devices.AccessRights(expandAccessRights(d)),
	}

	accessPolicies := make([]devices.SharedAccessSignatureAuthorizationRule, 0)

	alreadyExists := false
	for accessPolicyIterator, err := client.ListKeysComplete(ctx, id.ResourceGroup, id.IotHubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s: %+v", id, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if strings.EqualFold(*existingAccessPolicy.KeyName, id.IotHubKeyName) {
			if d.IsNewResource() {
				if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
					return tf.ImportAsExistsError("azurerm_iothub_shared_access_policy", id.ID())
				}
			}

			if existingAccessPolicy.PrimaryKey != nil {
				expandedAccessPolicy.PrimaryKey = existingAccessPolicy.PrimaryKey
			}

			if existingAccessPolicy.SecondaryKey != nil {
				expandedAccessPolicy.SecondaryKey = existingAccessPolicy.SecondaryKey
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
		return fmt.Errorf("dnable to find %s", id)
	}

	iothub.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return resourceIotHubSharedAccessPolicyRead(d, meta)
}

func expandAccessRights(d *pluginsdk.ResourceData) string {
	possibleAccessRights := []struct {
		schema string
		right  string
	}{
		{"registry_read", "RegistryRead"},
		{"registry_write", "RegistryWrite"},
		{"service_connect", "ServiceConnect"},
		{"device_connect", "DeviceConnect"},
	}
	actualRights := make([]string, 0)
	// iteration order is important here, so we cannot use a map
	for _, possibleRight := range possibleAccessRights {
		if d.Get(possibleRight.schema).(bool) {
			actualRights = append(actualRights, possibleRight.right)
		}
	}
	return strings.Join(actualRights, ", ")
}
