// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resource.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ServicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if m := resp.Model; m != nil {
		d.Set("location", location.Normalize(m.Location))

		if kind := m.Kind; string(kind) != "" {
			d.Set("kind", kind)
		}

		if err := d.Set("identity", identity.FlattenSystemAssigned(m.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := m.Properties; props != nil {
			if err := d.Set("access_policy_object_ids", flattenAccessPolicies(props.AccessPolicies)); err != nil {
				return fmt.Errorf("setting `access_policy_object_ids`: %+v", err)
			}

			cosmodDbKeyVaultKeyVersionlessId := ""
			cosmosDbThroughput := 0
			if cosmos := props.CosmosDbConfiguration; cosmos != nil {
				if v := cosmos.OfferThroughput; v != nil {
					cosmosDbThroughput = int(*v)
				}
				if v := cosmos.KeyVaultKeyUri; v != nil {
					cosmodDbKeyVaultKeyVersionlessId = *v
				}
			}
			d.Set("cosmosdb_key_vault_key_versionless_id", cosmodDbKeyVaultKeyVersionlessId)
			d.Set("cosmosdb_throughput", cosmosDbThroughput)

			if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
				d.Set("configuration_export_storage_account_name", props.ExportConfiguration.StorageAccountName)
			}

			if pointer.From(props.PublicNetworkAccess) == resource.PublicNetworkAccessEnabled {
				d.Set("public_network_access_enabled", true)
			} else {
				d.Set("public_network_access_enabled", false)
			}

			if err := d.Set("authentication_configuration", flattenAuthentication(props.AuthenticationConfiguration)); err != nil {
				return fmt.Errorf("setting `authentication_configuration`: %+v", err)
			}

			if err := d.Set("cors_configuration", flattenCorsConfig(props.CorsConfiguration)); err != nil {
				return fmt.Errorf("setting `cors_configuration`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, m.Tags); err != nil {
			return err
		}
	}

	return nil
}
