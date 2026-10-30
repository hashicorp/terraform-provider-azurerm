// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package notification_hub_namespace

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNotificationHubNamespaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NotificationHubs.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := namespaces.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_notification_hub_namespace", id.ID())
		}
	}

	zoneRedundancy := namespaces.ZoneRedundancyPreferenceDisabled
	if v, ok := d.GetOk("zone_redundancy_enabled"); ok && v.(bool) {
		zoneRedundancy = namespaces.ZoneRedundancyPreferenceEnabled
	}

	namespaceType := namespaces.NamespaceType(d.Get("namespace_type").(string))
	parameters := namespaces.NamespaceResource{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: namespaces.Sku{
			Name: namespaces.SkuName(d.Get("sku_name").(string)),
		},
		Properties: &namespaces.NamespaceProperties{
			NamespaceType:  &namespaceType,
			Enabled:        pointer.To(d.Get("enabled").(bool)),
			ZoneRedundancy: pointer.To(zoneRedundancy),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("replication_region"); ok {
		parameters.Properties.ReplicationRegion = pointer.ToEnum[namespaces.ReplicationRegion](location.Normalize(v.(string)))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be created..", id)
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              15 * time.Second,
		TargetStatusCode:      pointer.To(http.StatusOK),
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %ss to finish replicating: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNotificationHubNamespaceRead(d, meta)
}
