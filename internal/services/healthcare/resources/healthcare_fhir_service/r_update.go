// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_fhir_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisFhirServiceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	id := fhirservices.NewFhirServiceID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := fhirservices.FhirService{
		Identity: i,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     pointer.ToEnum[fhirservices.FhirServiceKind](d.Get("kind").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: &fhirservices.FhirServiceProperties{
			AuthenticationConfiguration: expandFhirAuthentication(d.Get("authentication").([]any)),
			CorsConfiguration:           expandFhirCorsConfiguration(d.Get("cors").([]any)),
			AccessPolicies:              expandAccessPolicy(d.Get("access_policy_object_ids").(*pluginsdk.Set).List()),
		},
	}

	storageAcc, hasValues := d.GetOk("configuration_export_storage_account_name")
	if hasValues {
		parameters.Properties.ExportConfiguration = &fhirservices.FhirServiceExportConfiguration{
			StorageAccountName: pointer.To(storageAcc.(string)),
		}
	}

	acrConfig := fhirservices.FhirServiceAcrConfiguration{}
	ociArtifactsRaw, hasValues := d.GetOk("oci_artifact")
	if hasValues {
		acrConfig.OciArtifacts = expandOciArtifacts(ociArtifactsRaw.([]any))
	}
	loginServersRaw, hasValues := d.GetOk("container_registry_login_server_url")
	if hasValues {
		acrConfig.LoginServers = expandFhirAcrLoginServer(loginServersRaw.(*pluginsdk.Set).List())
	}
	parameters.Properties.AcrConfiguration = &acrConfig

	if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}
