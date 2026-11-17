// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventhubs.ParseEventhubID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("partition_count") {
		o, n := d.GetChange("partition_count")
		if o.(int) > n.(int) {
			return fmt.Errorf("`partition_count` cannot be decreased")
		}

		client := meta.(*clients.Client).Eventhub.NamespacesClient
		namespaceId := namespaces.NewNamespaceID(subscriptionId, id.ResourceGroupName, id.NamespaceName)
		resp, err := client.Get(ctx, namespaceId)
		if err != nil {
			return err
		}
		if model := resp.Model; model != nil {
			if model.Sku.Name != namespaces.SkuNamePremium && model.Properties.ClusterArmId == nil {
				return fmt.Errorf("`partition_count` cannot be changed on shared namespaces unless the namespace sku is `Premium`")
			}
		}
	}

	eventhubStatus := eventhubs.EntityStatus(d.Get("status").(string))
	parameters := eventhubs.Eventhub{
		Properties: &eventhubs.EventhubProperties{
			PartitionCount:     pointer.To(int64(d.Get("partition_count").(int))),
			Status:             &eventhubStatus,
			CaptureDescription: expandEventHubCaptureDescription(d),
		},
	}

	rawConfig := d.GetRawConfig()
	messageRetention := rawConfig.GetAttr("message_retention")
	if messageRetention.IsKnown() && !messageRetention.IsNull() {
		parameters.Properties.MessageRetentionInDays = pointer.To(int64(d.Get("message_retention").(int)))
	}

	retentionDescription := rawConfig.GetAttr("retention_description")
	if retentionDescription.IsKnown() && !retentionDescription.IsNull() {
		parameters.Properties.RetentionDescription = expandEventHubRetentionDescription(d)
	}

	if d.HasChange("capture_description") {
		parameters.Properties.CaptureDescription = expandEventHubCaptureDescription(d)
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceEventHubRead(d, meta)
}
