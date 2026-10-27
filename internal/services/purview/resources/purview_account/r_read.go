// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package purview_account

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-12-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePurviewAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("managed_resources", flattenPurviewAccountManagedResources(props.ManagedResources)); err != nil {
				return fmt.Errorf("flattening `managed_resources`: %+v", err)
			}

			d.Set("public_network_enabled", pointer.From(props.PublicNetworkAccess) == account.PublicNetworkAccessEnabled)

			d.Set("managed_event_hub_enabled", pointer.From(props.ManagedEventHubState) == account.ManagedEventHubStateEnabled)

			d.Set("managed_resource_group_name", pointer.From(props.ManagedResourceGroupName))

			if endpoints := props.Endpoints; endpoints != nil {
				d.Set("catalog_endpoint", endpoints.Catalog)
				d.Set("guardian_endpoint", endpoints.Guardian)
				d.Set("scan_endpoint", endpoints.Scan)
			}

			if props.CloudConnectors != nil {
				d.Set("aws_external_id", pointer.From(props.CloudConnectors.AwsExternalId))
			}
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if keys, err := client.ListKeys(ctx, *id); err == nil {
		if model := keys.Model; model != nil {
			d.Set("atlas_kafka_endpoint_primary_connection_string", model.AtlasKafkaPrimaryEndpoint)
			d.Set("atlas_kafka_endpoint_secondary_connection_string", model.AtlasKafkaSecondaryEndpoint)
		}
	} else {
		// if eventhubs have been disabled we will get a response was not found, so we can ignore that error
		if !response.WasNotFound(keys.HttpResponse) {
			return fmt.Errorf("retrieving Keys for %s: %+v", *id, err)
		}

		d.Set("atlas_kafka_endpoint_primary_connection_string", "")
		d.Set("atlas_kafka_endpoint_secondary_connection_string", "")
	}

	return nil
}

func flattenPurviewAccountManagedResources(managedResources *account.ManagedResources) any {
	if managedResources == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"resource_group_id":      pointer.From(managedResources.ResourceGroup),
			"storage_account_id":     pointer.From(managedResources.StorageAccount),
			"event_hub_namespace_id": pointer.From(managedResources.EventHubNamespace),
		},
	}
}

func resourcePurviewAccountIdentityIdsSetHash(v any) int {
	var buf bytes.Buffer

	buf.WriteString(strings.ToLower(v.(string)))

	return pluginsdk.HashString(buf.String())
}
