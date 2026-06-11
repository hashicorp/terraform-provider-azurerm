// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineLearningRegistry struct{}

type ReplicationRegion struct {
	Location                           string `tfschema:"location"`
	SystemCreatedStorageAccountType    string `tfschema:"system_created_storage_account_type"`
	HnsEnabled                         bool   `tfschema:"system_created_storage_account_hns_enabled"`
	StorageAccountBlobPublicAccess     bool   `tfschema:"system_created_storage_account_blob_public_access_enabled"`
	SystemCreatedContainerRegistrySku  string `tfschema:"system_created_container_registry_sku"`
	SystemCreatedStorageAccountId      string `tfschema:"system_created_storage_account_id"`
	SystemCreatedStorageAccountName    string `tfschema:"system_created_storage_account_name"`
	SystemCreatedAcrId                 string `tfschema:"system_created_container_registry_id"`
	SystemCreatedContainerRegistryName string `tfschema:"system_created_container_registry_name"`
}

type RegistryPrivateLinkServiceConnectionState struct {
	ActionsRequired string `tfschema:"actions_required"`
	Description     string `tfschema:"description"`
	Status          string `tfschema:"status"`
}

type RegistryPrivateEndpointConnection struct {
	Id                string                                      `tfschema:"id"`
	Location          string                                      `tfschema:"location"`
	GroupIds          []string                                    `tfschema:"group_ids"`
	SubnetId          string                                      `tfschema:"subnet_id"`
	ProvisioningState string                                      `tfschema:"provisioning_state"`
	ConnectionState   []RegistryPrivateLinkServiceConnectionState `tfschema:"connection_state"`
}

type MachineLearningRegistryModel struct {
	Name                       string                                     `tfschema:"name"`
	ResourceGroupName          string                                     `tfschema:"resource_group_name"`
	Location                   string                                     `tfschema:"location"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccessEnabled bool                                       `tfschema:"public_network_access_enabled"`

	SystemCreatedStorageAccountType                    string `tfschema:"system_created_storage_account_type"`
	SystemCreatedStorageAccountHnsEnabled              bool   `tfschema:"system_created_storage_account_hns_enabled"`
	SystemCreatedStorageAccountBlobPublicAccessEnabled bool   `tfschema:"system_created_storage_account_blob_public_access_enabled"`
	SystemCreatedContainerRegistrySku                  string `tfschema:"system_created_container_registry_sku"`

	SystemCreatedStorageAccountId      string `tfschema:"system_created_storage_account_id"`
	SystemCreatedStorageAccountName    string `tfschema:"system_created_storage_account_name"`
	SystemCreatedContainerRegistryId   string `tfschema:"system_created_container_registry_id"`
	SystemCreatedContainerRegistryName string `tfschema:"system_created_container_registry_name"`

	ReplicationRegion []ReplicationRegion `tfschema:"replication_region"`

	PrivateEndpointConnection      []RegistryPrivateEndpointConnection `tfschema:"private_endpoint_connections"`
	DiscoveryUrl                   string                              `tfschema:"discovery_url"`
	MachineLearningFlowRegistryUri string                              `tfschema:"machine_learning_flow_registry_uri"`
	ManagedResourceGroup           string                              `tfschema:"managed_resource_group_id"`
	Tags                           map[string]string                   `tfschema:"tags"`
}

func (r MachineLearningRegistry) ModelObject() interface{} {
	return &MachineLearningRegistryModel{}
}

func (r MachineLearningRegistry) ResourceType() string {
	return "azurerm_machine_learning_registry"
}

func (r MachineLearningRegistry) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return registrymanagement.ValidateRegistryID
}

var _ sdk.ResourceWithUpdate = MachineLearningRegistry{}

var _ sdk.ResourceWithCustomizeDiff = MachineLearningRegistry{}

func (r MachineLearningRegistry) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{1,32}$`),
				"Machine Learning Registry name must be 2 - 33 characters long. Its first character has to be alphanumeric, and the rest may contain hyphens and underscores. No whitespace is allowed.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"replication_region": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: replicationRegionSchema(),
			},
		},

		"tags": commonschema.Tags(),
	}

	for k, v := range registryRegionConfigSchema() {
		arguments[k] = v
	}

	return arguments
}

func registryRegionConfigSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"system_created_container_registry_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(registrymanagement.SkuTierPremium),
			ValidateFunc: validation.StringInSlice([]string{
				string(registrymanagement.SkuTierPremium),
			}, false),
		},

		"system_created_storage_account_blob_public_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"system_created_storage_account_hns_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"system_created_storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  systemCreatedStorageAccountTypeDefault,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard_LRS",
				"Standard_GRS",
				"Standard_RAGRS",
				"Standard_ZRS",
				"Standard_GZRS",
				"Standard_RAGZRS",
				"Premium_LRS",
				"Premium_ZRS",
			}, false),
		},
	}
}

func registryRegionComputedSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"system_created_container_registry_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_created_container_registry_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_created_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_created_storage_account_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func replicationRegionSchema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"location": commonschema.LocationWithoutForceNew(),
	}

	for k, v := range registryRegionConfigSchema() {
		schema[k] = v
	}

	for k, v := range registryRegionComputedSchema() {
		schema[k] = v
	}

	return schema
}

func (r MachineLearningRegistry) Attributes() map[string]*pluginsdk.Schema {
	attributes := map[string]*pluginsdk.Schema{
		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"machine_learning_flow_registry_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"managed_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_endpoint_connections": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"location": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"group_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"connection_state": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"status": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"description": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"actions_required": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"provisioning_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	for k, v := range registryRegionComputedSchema() {
		attributes[k] = v
	}

	return attributes
}

func (r MachineLearningRegistry) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MachineLearningRegistryModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			primaryLocation := location.Normalize(model.Location)
			for _, region := range model.ReplicationRegion {
				if location.Normalize(region.Location) == primaryLocation {
					return fmt.Errorf("`replication_region` cannot contain the primary region `%s` specified in `location`", model.Location)
				}
			}

			return nil
		},
	}
}

func (r MachineLearningRegistry) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.RegistryManagement
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningRegistryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding Machine Learning Registry model %+v", err)
			}

			id := registrymanagement.NewRegistryID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.RegistriesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_registry", id.ID())
			}

			param := registrymanagement.RegistryTrackedResource{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: registrymanagement.Registry{
					PublicNetworkAccess: pointer.To(string(PublicNetworkAccessStateDisabled)),
				},
			}

			if model.PublicNetworkAccessEnabled {
				param.Properties.PublicNetworkAccess = pointer.To(string(PublicNetworkAccessStateEnabled))
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}

			param.Identity = expandedIdentity

			regions := expandRegistryRegions(model)
			param.Properties.RegionDetails = &regions

			if err := client.RegistriesCreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningRegistry) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.RegistryManagement

			id, err := registrymanagement.ParseRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.RegistriesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("model is nil %s", *id)
			}

			identityIds, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(resp.Model.Identity)
			if err != nil {
				return fmt.Errorf("flatten identity %s: %+v", *id, err)
			}

			prop := resp.Model.Properties
			model := MachineLearningRegistryModel{
				Name:                           id.RegistryName,
				ResourceGroupName:              id.ResourceGroupName,
				Identity:                       identityIds,
				Location:                       resp.Model.Location,
				PublicNetworkAccessEnabled:     pointer.From(prop.PublicNetworkAccess) == string(PublicNetworkAccessStateEnabled),
				Tags:                           pointer.From(resp.Model.Tags),
				MachineLearningFlowRegistryUri: pointer.From(prop.MlFlowRegistryUri),
				DiscoveryUrl:                   pointer.From(prop.DiscoveryURL),
				PrivateEndpointConnection:      flattenRegistryPrivateEndpointConnections(prop.RegistryPrivateEndpointConnections),
			}

			if prop.ManagedResourceGroup != nil {
				model.ManagedResourceGroup = pointer.From(prop.ManagedResourceGroup.ResourceId)
			}

			regions := flattenRegistryRegionDetails(prop.RegionDetails)
			for i, region := range regions {
				if i == 0 {
					model.SystemCreatedStorageAccountType = region.SystemCreatedStorageAccountType
					model.SystemCreatedStorageAccountHnsEnabled = region.HnsEnabled
					model.SystemCreatedStorageAccountBlobPublicAccessEnabled = metadata.ResourceData.Get("system_created_storage_account_blob_public_access_enabled").(bool)
					model.SystemCreatedContainerRegistrySku = region.SystemCreatedContainerRegistrySku
					model.SystemCreatedStorageAccountId = region.SystemCreatedStorageAccountId
					model.SystemCreatedStorageAccountName = region.SystemCreatedStorageAccountName
					model.SystemCreatedContainerRegistryId = region.SystemCreatedAcrId
					model.SystemCreatedContainerRegistryName = region.SystemCreatedContainerRegistryName
					continue
				}
				region.StorageAccountBlobPublicAccess = metadata.ResourceData.Get(fmt.Sprintf("replication_region.%d.system_created_storage_account_blob_public_access_enabled", i-1)).(bool)
				model.ReplicationRegion = append(model.ReplicationRegion, region)
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MachineLearningRegistry) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.RegistryManagement

			var model MachineLearningRegistryModel
			id, err := registrymanagement.ParseRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding Machine Learning Registry model %+v", err)
			}

			existing, err := client.RegistriesGet(ctx, *id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving existing model for %s", id)
			}

			param := *existing.Model

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding identity: %+v", err)
				}

				param.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				param.Properties.PublicNetworkAccess = pointer.To(string(PublicNetworkAccessStateDisabled))
				if model.PublicNetworkAccessEnabled {
					param.Properties.PublicNetworkAccess = pointer.To(string(PublicNetworkAccessStateEnabled))
				}
			}

			if metadata.ResourceData.HasChanges(
				"replication_region",
				"system_created_storage_account_type",
				"system_created_storage_account_hns_enabled",
				"system_created_storage_account_blob_public_access_enabled",
				"system_created_container_registry_sku",
			) {
				regions := expandRegistryRegions(model)
				param.Properties.RegionDetails = &regions
			}

			if err := client.RegistriesCreateOrUpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningRegistry) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.RegistryManagement

			id, err := registrymanagement.ParseRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.RegistriesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandRegistryRegionDetail(input ReplicationRegion) registrymanagement.RegistryRegionArmDetails {
	acrSku := input.SystemCreatedContainerRegistrySku
	if acrSku == "" {
		acrSku = string(registrymanagement.SkuTierPremium)
	}

	return registrymanagement.RegistryRegionArmDetails{
		Location: pointer.To(input.Location),
		AcrDetails: &[]registrymanagement.AcrDetails{
			{
				SystemCreatedAcrAccount: &registrymanagement.SystemCreatedAcrAccount{
					AcrAccountSku: pointer.To(acrSku),
				},
			},
		},
		StorageAccountDetails: &[]registrymanagement.StorageAccountDetails{
			{
				SystemCreatedStorageAccount: &registrymanagement.SystemCreatedStorageAccount{
					StorageAccountHnsEnabled: pointer.To(input.HnsEnabled),
					StorageAccountType:       pointer.To(input.SystemCreatedStorageAccountType),
					AllowBlobPublicAccess:    pointer.To(input.StorageAccountBlobPublicAccess),
				},
			},
		},
	}
}

func expandRegistryRegions(model MachineLearningRegistryModel) []registrymanagement.RegistryRegionArmDetails {
	regions := make([]registrymanagement.RegistryRegionArmDetails, 0)

	regions = append(regions, expandRegistryRegionDetail(ReplicationRegion{
		Location:                          model.Location,
		SystemCreatedStorageAccountType:   model.SystemCreatedStorageAccountType,
		HnsEnabled:                        model.SystemCreatedStorageAccountHnsEnabled,
		StorageAccountBlobPublicAccess:    model.SystemCreatedStorageAccountBlobPublicAccessEnabled,
		SystemCreatedContainerRegistrySku: model.SystemCreatedContainerRegistrySku,
	}))

	for _, region := range model.ReplicationRegion {
		regions = append(regions, expandRegistryRegionDetail(region))
	}

	return regions
}

func flattenRegistryRegionDetails(input *[]registrymanagement.RegistryRegionArmDetails) []ReplicationRegion {
	result := make([]ReplicationRegion, 0)
	if input == nil || len(*input) == 0 {
		return result
	}

	for _, item := range *input {
		var region ReplicationRegion
		region.Location = pointer.From(item.Location)

		if sa := pointer.From(item.StorageAccountDetails); len(sa) > 0 {
			if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
				region.SystemCreatedStorageAccountType = pointer.From(systemAccount.StorageAccountType)
				region.HnsEnabled = pointer.From(systemAccount.StorageAccountHnsEnabled)
				region.StorageAccountBlobPublicAccess = pointer.From(systemAccount.AllowBlobPublicAccess)
				region.SystemCreatedStorageAccountName = pointer.From(systemAccount.StorageAccountName)

				if systemAccount.ArmResourceId != nil {
					region.SystemCreatedStorageAccountId = pointer.From(systemAccount.ArmResourceId.ResourceId)
				}
			}
		}

		if acr := pointer.From(item.AcrDetails); len(acr) > 0 {
			if systemAcr := acr[0].SystemCreatedAcrAccount; systemAcr != nil {
				region.SystemCreatedContainerRegistrySku = pointer.From(systemAcr.AcrAccountSku)
				region.SystemCreatedContainerRegistryName = pointer.From(systemAcr.AcrAccountName)

				if systemAcr.ArmResourceId != nil {
					region.SystemCreatedAcrId = pointer.From(systemAcr.ArmResourceId.ResourceId)
				}
			}
		}

		result = append(result, region)
	}
	return result
}

func flattenRegistryPrivateEndpointConnections(input *[]registrymanagement.RegistryPrivateEndpointConnection) []RegistryPrivateEndpointConnection {
	result := make([]RegistryPrivateEndpointConnection, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		connection := RegistryPrivateEndpointConnection{
			Id:       pointer.From(item.Id),
			Location: pointer.From(item.Location),
		}

		if props := item.Properties; props != nil {
			connection.GroupIds = pointer.From(props.GroupIds)
			connection.ProvisioningState = pointer.From(props.ProvisioningState)

			if props.PrivateEndpoint != nil {
				connection.SubnetId = pointer.From(props.PrivateEndpoint.SubnetArmId)
			}

			if state := props.RegistryPrivateLinkServiceConnectionState; state != nil {
				connection.ConnectionState = []RegistryPrivateLinkServiceConnectionState{
					{
						Status:          string(pointer.From(state.Status)),
						Description:     pointer.From(state.Description),
						ActionsRequired: pointer.From(state.ActionsRequired),
					},
				}
			}
		}

		result = append(result, connection)
	}

	return result
}

type PublicNetworkAccessState string

const (
	PublicNetworkAccessStateEnabled  PublicNetworkAccessState = "Enabled"
	PublicNetworkAccessStateDisabled PublicNetworkAccessState = "Disabled"
)

const systemCreatedStorageAccountTypeDefault = "Standard_LRS"
