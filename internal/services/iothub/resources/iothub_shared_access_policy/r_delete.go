// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_shared_access_policy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubSharedAccessPolicyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if response.WasNotFound(iothub.Response.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	accessPolicies := make([]devices.SharedAccessSignatureAuthorizationRule, 0)

	for accessPolicyIterator, err := client.ListKeysComplete(ctx, id.ResourceGroup, id.IotHubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s: %+v", id, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if !strings.EqualFold(*existingAccessPolicy.KeyName, id.IotHubKeyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothub.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}
