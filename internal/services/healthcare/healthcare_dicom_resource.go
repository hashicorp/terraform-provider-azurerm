// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisDicomService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisDicomServiceCreate,
		Read:   resourceHealthcareApisDicomServiceRead,
		Update: resourceHealthcareApisDicomServiceUpdate,
		Delete: resourceHealthcareApisDicomServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareDicomServiceV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dicomservices.ParseDicomServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DicomServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"authentication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"audience": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					},
				},
			},

			"private_endpoint": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// NOTE: O+C API sets defaults for these if omitted
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"allowed_headers": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"max_age_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 99998),
						},

						"allow_credentials": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"data_partitions_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"encryption_key_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// NOTE: O+C API sets defaults for these if omitted
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"file_system_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"storage_account_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: commonids.ValidateStorageAccountID,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisDicomServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Dicom Service creation.")

	workspace, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare workspace error: %+v", err)
	}

	id := dicomservices.NewDicomServiceID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
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

	i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	t := d.Get("tags").(map[string]interface{})

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

	cors := expandDicomServiceCorsConfiguration(d.Get("cors").([]interface{}))
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

	storage := expandStorageConfiguration(d.Get("storage").([]interface{}))
	if storage != nil {
		parameters.Properties.StorageConfiguration = storage
	}

	if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
		parameters.Properties.PublicNetworkAccess = pointer.To(dicomservices.PublicNetworkAccessDisabled)
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func resourceHealthcareApisDicomServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dicomservices.ParseDicomServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Dicom service error: %+v", err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DicomServiceName)
	d.Set("workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))

		if props := m.Properties; props != nil {
			d.Set("authentication", flattenDicomAuthentication(props.AuthenticationConfiguration))
			d.Set("private_endpoint", flattenDicomServicePrivateEndpoint(props.PrivateEndpointConnections))
			d.Set("service_url", props.ServiceURL)

			if pna := pointer.From(props.PublicNetworkAccess); pna != "" {
				d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == dicomservices.PublicNetworkAccessEnabled)
			}

			d.Set("data_partitions_enabled", pointer.From(props.EnableDataPartitions))

			d.Set("cors", flattenDicomServiceCorsConfiguration(props.CorsConfiguration))

			if props.Encryption != nil && props.Encryption.CustomerManagedKeyEncryption != nil {
				d.Set("encryption_key_url", pointer.From(props.Encryption.CustomerManagedKeyEncryption.KeyEncryptionKeyURL))
			}

			storage, err := flattenStorageConfiguration(props.StorageConfiguration)
			if err != nil {
				return fmt.Errorf("flattening `storage`: %+v", err)
			}
			d.Set("storage", storage)
		}

		i, err := identity.FlattenLegacySystemAndUserAssignedMap(m.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", i); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		return tags.FlattenAndSet(d, m.Tags)
	}
	return nil
}

func resourceHealthcareApisDicomServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dicomservices.ParseDicomServiceID(d.Id())
	if err != nil {
		return err
	}

	// Retrieve the existing resource
	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("cors") {
		payload.Properties.CorsConfiguration = expandDicomServiceCorsConfiguration(d.Get("cors").([]interface{}))
	}

	if d.HasChange("encryption_key_url") {
		payload.Properties.Encryption = &dicomservices.Encryption{
			CustomerManagedKeyEncryption: &dicomservices.EncryptionCustomerManagedKeyEncryption{
				KeyEncryptionKeyURL: pointer.To(d.Get("encryption_key_url").(string)),
			},
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if enabled := d.Get("public_network_access_enabled").(bool); !enabled {
			payload.Properties.PublicNetworkAccess = pointer.To(dicomservices.PublicNetworkAccessDisabled)
		} else {
			payload.Properties.PublicNetworkAccess = pointer.To(dicomservices.PublicNetworkAccessEnabled)
		}
	}

	if d.HasChange("identity") {
		i, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = i
	}

	if d.HasChange("location") {
		payload.Location = pointer.To(location.Normalize(d.Get("location").(string)))
	}

	if d.HasChange("tags") {
		if err := updateTags(d, meta); err != nil {
			return fmt.Errorf("updating tags error: %+v", err)
		}
	}

	err = client.CreateOrUpdateThenPoll(ctx, *id, *payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareApisDicomServiceRead(d, meta)
}

func resourceHealthcareApisDicomServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dicomservices.ParseDicomServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Dicom service error: %+v", err)
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Healthcare Dicom Service %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Pending"},
		Target:                    []string{"Deleted"},
		Refresh:                   dicomServiceStateStatusCodeRefreshFunc(ctx, client, *id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 3,
		PollInterval:              10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func dicomServiceStateStatusCodeRefreshFunc(ctx context.Context, client *dicomservices.DicomServicesClient, id dicomservices.DicomServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return resp, "Pending", nil
	}
}

func updateTags(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceDicomServiceClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dicomservices.ParseDicomServiceID(d.Id())
	if err != nil {
		return err
	}

	update := dicomservices.DicomServicePatchResource{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	err = client.UpdateThenPoll(ctx, *id, update)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return nil
}

func flattenDicomAuthentication(input *dicomservices.DicomServiceAuthenticationConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	authBlock := make(map[string]interface{})
	if input.Authority != nil {
		authBlock["authority"] = *input.Authority
	}

	audience := make([]interface{}, 0)
	if input.Audiences != nil {
		for _, data := range *input.Audiences {
			audience = append(audience, data)
		}
	}
	authBlock["audience"] = audience

	return []interface{}{authBlock}
}

func flattenDicomServicePrivateEndpoint(input *[]dicomservices.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]interface{}{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}

		if endpoint.Id != nil {
			result["id"] = *endpoint.Id
		}
	}
	return results
}

func expandStorageConfiguration(input []interface{}) *dicomservices.StorageConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	storageSettings := input[0].(map[string]interface{})

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

func flattenStorageConfiguration(configuration *dicomservices.StorageConfiguration) (interface{}, error) {
	if configuration == nil {
		return []interface{}{}, nil
	}

	result := make(map[string]interface{})
	if configuration.FileSystemName != nil {
		result["file_system_name"] = pointer.From(configuration.FileSystemName)
	}

	if v := pointer.From(configuration.StorageResourceId); v != "" {
		id, err := commonids.ParseStorageAccountID(v)
		if err != nil {
			return nil, err
		}
		result["storage_account_id"] = id.ID()
	}

	return []interface{}{result}, nil
}

func expandDicomServiceCorsConfiguration(inputList []interface{}) *dicomservices.CorsConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList[0].(map[string]interface{})
	output := dicomservices.CorsConfiguration{}

	if v, ok := input["allowed_origins"]; ok {
		output.Origins = utils.ExpandStringSlice(v.([]interface{}))
	}

	if v, ok := input["allowed_headers"]; ok {
		output.Headers = utils.ExpandStringSlice(v.([]interface{}))
	}

	if v, ok := input["allowed_methods"]; ok {
		output.Methods = utils.ExpandStringSlice(v.([]interface{}))
	}

	if v, ok := input["max_age_in_seconds"]; ok {
		output.MaxAge = pointer.To(int64(v.(int)))
	}

	if v, ok := input["allow_credentials"]; ok {
		output.AllowCredentials = pointer.To(v.(bool))
	}

	return &output
}

func flattenDicomServiceCorsConfiguration(input *dicomservices.CorsConfiguration) []interface{} {
	outputList := make([]interface{}, 0)
	if input == nil {
		return outputList
	}

	output := make(map[string]interface{})
	output["allow_credentials"] = pointer.From(input.AllowCredentials)

	if input.Headers != nil {
		output["allowed_headers"] = utils.FlattenStringSlice(input.Headers)
	}

	output["max_age_in_seconds"] = pointer.From(input.MaxAge)

	if input.Methods != nil {
		output["allowed_methods"] = utils.FlattenStringSlice(input.Methods)
	}

	if input.Origins != nil {
		output["allowed_origins"] = utils.FlattenStringSlice(input.Origins)
	}

	return append(outputList, output)
}
