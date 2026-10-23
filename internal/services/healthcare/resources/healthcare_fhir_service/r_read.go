// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_fhir_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisFhirServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fhirservices.ParseFhirServiceID(d.Id())
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
	d.Set("name", id.FhirServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	workSpaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("workspace_id", workSpaceId.ID())

	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))

		i, err := identity.FlattenLegacySystemAndUserAssignedMap(m.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		d.Set("kind", pointer.FromEnum(m.Kind))

		if props := m.Properties; props != nil {
			d.Set("access_policy_object_ids", flattenFhirAccessPolicy(props.AccessPolicies))
			d.Set("authentication", flattenFhirAuthentication(props.AuthenticationConfiguration))
			d.Set("cors", flattenFhirCorsConfiguration(props.CorsConfiguration))
			d.Set("container_registry_login_server_url", flattenFhirAcrLoginServer(props.AcrConfiguration))
			if acrConfig := props.AcrConfiguration; acrConfig != nil {
				if artifacts := acrConfig.OciArtifacts; artifacts != nil {
					d.Set("oci_artifact", flattenOciArtifacts(artifacts))
				}
			}
			if props.ExportConfiguration != nil && props.ExportConfiguration.StorageAccountName != nil {
				d.Set("configuration_export_storage_account_name", props.ExportConfiguration.StorageAccountName)
			}
			if props.PublicNetworkAccess != nil {
				d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == fhirservices.PublicNetworkAccessEnabled)
			}

			if err := tags.FlattenAndSet(d, m.Tags); err != nil {
				return err
			}
		}
	}
	return nil
}
