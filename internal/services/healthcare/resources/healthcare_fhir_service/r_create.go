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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisFhirServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceFhirServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	id := fhirservices.NewFhirServiceID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_healthcare_fhir_service", id.ID())
		}
	}

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
		},
	}

	accessPolicyObjectIds, hasValues := d.GetOk("access_policy_object_ids")
	if hasValues {
		parameters.Properties.AccessPolicies = expandAccessPolicy(accessPolicyObjectIds.(*pluginsdk.Set).List())
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

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisFhirServiceRead(d, meta)
}
