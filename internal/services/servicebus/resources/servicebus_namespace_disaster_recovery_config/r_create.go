// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_disaster_recovery_config

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	servicebusClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceServiceBusNamespaceDisasterRecoveryConfigCreate(d *pluginsdk.ResourceData, meta any) error {
	client := servicebusClient.From(meta).ArmDisasterRecoveriesClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	namespaceId, err := armdisasterrecoveries.ParseNamespaceID(d.Get("primary_namespace_id").(string))
	if err != nil {
		return err
	}

	partnerNamespaceId := d.Get("partner_namespace_id").(string)

	id := armdisasterrecoveries.NewDisasterRecoveryConfigID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.DisasterRecoveryConfigsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_disaster_recovery_config", id.ID())
		}
	}

	parameters := armdisasterrecoveries.ArmDisasterRecovery{
		Properties: &armdisasterrecoveries.ArmDisasterRecoveryProperties{
			PartnerNamespace: pointer.To(partnerNamespaceId),
		},
	}

	if _, err := client.DisasterRecoveryConfigsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	return resourceServiceBusNamespaceDisasterRecoveryConfigRead(d, meta)
}
