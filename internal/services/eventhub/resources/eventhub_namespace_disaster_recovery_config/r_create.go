// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_disaster_recovery_config

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceDisasterRecoveryConfigCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_disaster_recovery_config", id.ID())
		}
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	parameters := disasterrecoveryconfigs.ArmDisasterRecovery{
		Properties: &disasterrecoveryconfigs.ArmDisasterRecoveryProperties{
			PartnerNamespace: pointer.To(d.Get("partner_namespace_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id); err != nil {
		return fmt.Errorf("waiting for replication of %s: %+v", id, err)
	}

	return resourceEventHubNamespaceDisasterRecoveryConfigRead(d, meta)
}
