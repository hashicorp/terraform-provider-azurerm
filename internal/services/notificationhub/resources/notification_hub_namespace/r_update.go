// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_namespace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubNamespaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	parameters := namespaces.NamespacePatchParameters{
		Properties: &namespaces.NamespaceProperties{
			NamespaceType: pointer.ToEnum[namespaces.NamespaceType](d.Get("namespace_type").(string)),
			Enabled:       pointer.To(d.Get("enabled").(bool)),
		},
	}

	if d.HasChange("sku_name") {
		parameters.Sku = &namespaces.Sku{
			Name: namespaces.SkuName(d.Get("sku_name").(string)),
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceNotificationHubNamespaceRead(d, meta)
}
