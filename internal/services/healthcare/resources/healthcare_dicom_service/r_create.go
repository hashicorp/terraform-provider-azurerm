// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_dicom_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisDicomServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare workspace error: %+v", err)
	}

	id := dicomservices.NewDicomServiceID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_healthcare_dicom_service", id.ID())
		}
	}

	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	t := d.Get("tags").(map[string]any)

	parameters := dicomservices.DicomService{
		Identity: i,
		Properties: &dicomservices.DicomServiceProperties{
			PublicNetworkAccess: pointer.To(dicomservices.PublicNetworkAccessEnabled),
		},
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(t),
	}

	if v, ok := d.GetOk("data_partitions_enabled"); ok {
		parameters.Properties.EnableDataPartitions = pointer.To(v.(bool))
	}

	cors := expandDicomServiceCorsConfiguration(d.Get("cors").([]any))
	if cors != nil {
		parameters.Properties.CorsConfiguration = cors
	}

	if v, ok := d.GetOk("encryption_key_url"); ok && v.(string) != "" {
		parameters.Properties.Encryption = &dicomservices.Encryption{
			CustomerManagedKeyEncryption: &dicomservices.EncryptionCustomerManagedKeyEncryption{
				KeyEncryptionKeyURL: pointer.To(v.(string)),
			},
		}
	}

	storage := expandStorageConfiguration(d.Get("storage").([]any))
	if storage != nil {
		parameters.Properties.StorageConfiguration = storage
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.Properties.PublicNetworkAccess = pointer.To(dicomservices.PublicNetworkAccessDisabled)
	}

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func expandStorageConfiguration(input []any) *dicomservices.StorageConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	storageSettings := input[0].(map[string]any)

	var storageResourceId *string
	if v, ok := storageSettings["storage_account_id"]; ok {
		storageResourceId = pointer.To(v.(string))
	}

	var fileSystemName *string
	if v, ok := storageSettings["file_system_name"]; ok {
		fileSystemName = pointer.To(v.(string))
	}

	return &dicomservices.StorageConfiguration{
		FileSystemName:    fileSystemName,
		StorageResourceId: storageResourceId,
	}
}

func expandDicomServiceCorsConfiguration(inputList []any) *dicomservices.CorsConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList[0].(map[string]any)
	output := dicomservices.CorsConfiguration{}

	if v, ok := input["allowed_origins"]; ok {
		output.Origins = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if v, ok := input["allowed_headers"]; ok {
		output.Headers = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if v, ok := input["allowed_methods"]; ok {
		output.Methods = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if v, ok := input["max_age_in_seconds"]; ok {
		output.MaxAge = pointer.To(int64(v.(int)))
	}

	if v, ok := input["allow_credentials"]; ok {
		output.AllowCredentials = pointer.To(v.(bool))
	}

	return &output
}
