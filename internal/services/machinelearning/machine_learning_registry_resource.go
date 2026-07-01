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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name machine_learning_registry -service-package-name machinelearning -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type MachineLearningRegistryResource struct{}

type ReplicationRegion struct {
	Location                           string `tfschema:"location"`
	SystemCreatedStorageAccountType    string `tfschema:"system_created_storage_account_type"`
	HierarchicalNamespaceEnabled       bool   `tfschema:"system_created_storage_account_hierarchical_namespace_enabled"`
	SystemCreatedContainerRegistrySku  string `tfschema:"system_created_container_registry_sku"`
	SystemCreatedStorageAccountId      string `tfschema:"system_created_storage_account_id"`
	SystemCreatedStorageAccountName    string `tfschema:"system_created_storage_account_name"`
	SystemCreatedAcrId                 string `tfschema:"system_created_container_registry_id"`
	SystemCreatedContainerRegistryName string `tfschema:"system_created_container_registry_name"`
}

type MachineLearningRegistryModel struct {
	Name                                                    string                                     `tfschema:"name"`
	ResourceGroupName                                       string                                     `tfschema:"resource_group_name"`
	Location                                                string                                     `tfschema:"location"`
	Identity                                                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccessEnabled                              bool                                       `tfschema:"public_network_access_enabled"`
	SystemCreatedStorageAccountType                         string                                     `tfschema:"system_created_storage_account_type"`
	SystemCreatedStorageAccountHierarchicalNamespaceEnabled bool                                       `tfschema:"system_created_storage_account_hierarchical_namespace_enabled"`
	SystemCreatedContainerRegistrySku                       string                                     `tfschema:"system_created_container_registry_sku"`
	SystemCreatedStorageAccountId                           string                                     `tfschema:"system_created_storage_account_id"`
	SystemCreatedStorageAccountName                         string                                     `tfschema:"system_created_storage_account_name"`
	SystemCreatedContainerRegistryId                        string                                     `tfschema:"system_created_container_registry_id"`
	SystemCreatedContainerRegistryName                      string                                     `tfschema:"system_created_container_registry_name"`
	ReplicationRegion                                       []ReplicationRegion                        `tfschema:"replication_region"`
	DiscoveryUrl                                            string                                     `tfschema:"discovery_url"`
	IntellectualPropertyPublisher                           string                                     `tfschema:"intellectual_property_publisher"`
	MachineLearningFlowRegistryUri                          string                                     `tfschema:"machine_learning_flow_registry_uri"`
	ManagedResourceGroup                                    string                                     `tfschema:"managed_resource_group_id"`
	Tags                                                    map[string]string                          `tfschema:"tags"`
}

func (r MachineLearningRegistryResource) ModelObject() interface{} {
	return &MachineLearningRegistryModel{}
}

func (r MachineLearningRegistryResource) ResourceType() string {
	return "azurerm_machine_learning_registry"
}

func (r MachineLearningRegistryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return registrymanagement.ValidateRegistryID
}

func (r MachineLearningRegistryResource) Identity() resourceids.ResourceId {
	return &registrymanagement.RegistryId{}
}

var (
	_ sdk.ResourceWithUpdate        = MachineLearningRegistryResource{}
	_ sdk.ResourceWithCustomizeDiff = MachineLearningRegistryResource{}
	_ sdk.ResourceWithIdentity      = MachineLearningRegistryResource{}
)

func (r MachineLearningRegistryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{2,32}$`),
				"Machine Learning Registry name must be between 3 and 33 characters long. Its first character has to be alphanumeric, and the rest may contain hyphens and underscores. No whitespace is allowed.",
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

		"system_created_container_registry_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(registrymanagement.SkuTierPremium),
			ValidateFunc: validation.StringInSlice([]string{
				string(registrymanagement.SkuTierPremium),
			}, false),
		},

		"system_created_storage_account_hierarchical_namespace_enabled": {
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

		"tags": commonschema.Tags(),
	}
}

func replicationRegionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationWithoutForceNew(),

		"system_created_container_registry_sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(registrymanagement.SkuTierPremium),
			ValidateFunc: validation.StringInSlice([]string{
				string(registrymanagement.SkuTierPremium),
			}, false),
		},

		"system_created_storage_account_hierarchical_namespace_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"system_created_storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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

func (r MachineLearningRegistryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"intellectual_property_publisher": {
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

func (r MachineLearningRegistryResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MachineLearningRegistryModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			for _, region := range model.ReplicationRegion {
				if location.Normalize(region.Location) == location.Normalize(model.Location) {
					return fmt.Errorf("`replication_region` cannot contain the primary `Location`: `%s`", model.Location)
				}
			}

			if metadata.ResourceDiff.Id() != "" {
				oldVal, newVal := metadata.ResourceDiff.GetChange("replication_region")
				oldRegions := oldVal.([]interface{})
				newRegions := newVal.([]interface{})

				newRegionsByLocation := make(map[string]map[string]interface{})
				for _, r := range newRegions {
					region := r.(map[string]interface{})
					newRegionsByLocation[location.Normalize(region["location"].(string))] = region
				}

				for _, r := range oldRegions {
					oldRegion := r.(map[string]interface{})
					oldLocation := location.Normalize(oldRegion["location"].(string))
					newRegion, exists := newRegionsByLocation[oldLocation]
					if !exists {
						return fmt.Errorf("removing a `replication_region` is not supported, region `%s` cannot be removed", oldRegion["location"].(string))
					}

					if oldRegion["system_created_storage_account_type"].(string) != newRegion["system_created_storage_account_type"].(string) ||
						oldRegion["system_created_storage_account_hierarchical_namespace_enabled"].(bool) != newRegion["system_created_storage_account_hierarchical_namespace_enabled"].(bool) ||
						oldRegion["system_created_container_registry_sku"].(string) != newRegion["system_created_container_registry_sku"].(string) {
						return fmt.Errorf("updating properties of an existing `replication_region` is not supported, region `%s` cannot be modified", oldRegion["location"].(string))
					}
				}
			}

			return nil
		},
	}
}

func (r MachineLearningRegistryResource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.RegistriesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := registrymanagement.RegistryTrackedResource{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: registrymanagement.Registry{
					PublicNetworkAccess: pointer.To(string(PublicNetworkAccessStateEnabled)),
					RegionDetails:       pointer.To(expandRegistryRegions(model)),
				},
			}

			if !model.PublicNetworkAccessEnabled {
				param.Properties.PublicNetworkAccess = pointer.To(string(PublicNetworkAccessStateDisabled))
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			param.Identity = expandedIdentity

			if err := client.RegistriesCreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r MachineLearningRegistryResource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r MachineLearningRegistryResource) flatten(metadata sdk.ResourceMetaData, id *registrymanagement.RegistryId, registry *registrymanagement.RegistryTrackedResource) error {
	identityIds, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(registry.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity` %s: %+v", *id, err)
	}

	prop := registry.Properties
	model := MachineLearningRegistryModel{
		Name:                           id.RegistryName,
		ResourceGroupName:              id.ResourceGroupName,
		Identity:                       identityIds,
		Location:                       registry.Location,
		PublicNetworkAccessEnabled:     pointer.From(prop.PublicNetworkAccess) == string(PublicNetworkAccessStateEnabled),
		Tags:                           pointer.From(registry.Tags),
		MachineLearningFlowRegistryUri: pointer.From(prop.MlFlowRegistryUri),
		DiscoveryUrl:                   pointer.From(prop.DiscoveryURL),
		IntellectualPropertyPublisher:  pointer.From(prop.IntellectualPropertyPublisher),
	}

	if prop.ManagedResourceGroup != nil {
		resourceGroupId, err := commonids.ParseResourceGroupID(pointer.From(prop.ManagedResourceGroup.ResourceId))
		if err != nil {
			return err
		}
		model.ManagedResourceGroup = resourceGroupId.ID()
	}

	regions, err := flattenRegistryRegionDetails(prop.RegionDetails)
	if err != nil {
		return fmt.Errorf("flattening `region_details` %s: %+v", *id, err)
	}

	for _, region := range regions {
		if location.Normalize(region.Location) == location.Normalize(registry.Location) {
			model.SystemCreatedStorageAccountType = region.SystemCreatedStorageAccountType
			model.SystemCreatedStorageAccountHierarchicalNamespaceEnabled = region.HierarchicalNamespaceEnabled
			model.SystemCreatedContainerRegistrySku = region.SystemCreatedContainerRegistrySku
			model.SystemCreatedStorageAccountId = region.SystemCreatedStorageAccountId
			model.SystemCreatedStorageAccountName = region.SystemCreatedStorageAccountName
			model.SystemCreatedContainerRegistryId = region.SystemCreatedAcrId
			model.SystemCreatedContainerRegistryName = region.SystemCreatedContainerRegistryName
		} else {
			model.ReplicationRegion = append(model.ReplicationRegion, region)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&model)
}

func (r MachineLearningRegistryResource) Update() sdk.ResourceFunc {
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
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.RegistriesGet(ctx, *id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			param := *existing.Model

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				param.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				param.Properties.PublicNetworkAccess = pointer.To(string(PublicNetworkAccessStateDisabled))
				if model.PublicNetworkAccessEnabled {
					param.Properties.PublicNetworkAccess = pointer.To(string(PublicNetworkAccessStateEnabled))
				}
			}

			if metadata.ResourceData.HasChange("replication_region") {
				param.Properties.RegionDetails = pointer.To(expandRegistryRegions(model))
			}

			if err := client.RegistriesCreateOrUpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MachineLearningRegistryResource) Delete() sdk.ResourceFunc {
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
	return registrymanagement.RegistryRegionArmDetails{
		Location: pointer.To(input.Location),
		AcrDetails: &[]registrymanagement.AcrDetails{
			{
				SystemCreatedAcrAccount: &registrymanagement.SystemCreatedAcrAccount{
					AcrAccountSku: pointer.To(input.SystemCreatedContainerRegistrySku),
				},
			},
		},
		StorageAccountDetails: &[]registrymanagement.StorageAccountDetails{
			{
				SystemCreatedStorageAccount: &registrymanagement.SystemCreatedStorageAccount{
					StorageAccountHnsEnabled: pointer.To(input.HierarchicalNamespaceEnabled),
					StorageAccountType:       pointer.To(input.SystemCreatedStorageAccountType),
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
		HierarchicalNamespaceEnabled:      model.SystemCreatedStorageAccountHierarchicalNamespaceEnabled,
		SystemCreatedContainerRegistrySku: model.SystemCreatedContainerRegistrySku,
	}))

	for _, region := range model.ReplicationRegion {
		regions = append(regions, expandRegistryRegionDetail(region))
	}

	return regions
}

func flattenRegistryRegionDetails(input *[]registrymanagement.RegistryRegionArmDetails) ([]ReplicationRegion, error) {
	result := make([]ReplicationRegion, 0)
	if input == nil || len(*input) == 0 {
		return result, nil
	}

	for _, item := range *input {
		var region ReplicationRegion
		region.Location = pointer.From(item.Location)

		if sa := pointer.From(item.StorageAccountDetails); len(sa) > 0 {
			if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
				region.SystemCreatedStorageAccountType = pointer.From(systemAccount.StorageAccountType)
				region.HierarchicalNamespaceEnabled = pointer.From(systemAccount.StorageAccountHnsEnabled)
				region.SystemCreatedStorageAccountName = pointer.From(systemAccount.StorageAccountName)

				if systemAccount.ArmResourceId != nil {
					storageAccountId, err := commonids.ParseStorageAccountID(pointer.From(systemAccount.ArmResourceId.ResourceId))
					if err != nil {
						return make([]ReplicationRegion, 0), err
					}
					region.SystemCreatedStorageAccountId = storageAccountId.ID()
				}
			}
		}

		if acr := pointer.From(item.AcrDetails); len(acr) > 0 {
			if systemAcr := acr[0].SystemCreatedAcrAccount; systemAcr != nil {
				region.SystemCreatedContainerRegistrySku = pointer.From(systemAcr.AcrAccountSku)
				region.SystemCreatedContainerRegistryName = pointer.From(systemAcr.AcrAccountName)

				if systemAcr.ArmResourceId != nil {
					containerRegistryId, err := registries.ParseRegistryID(pointer.From(systemAcr.ArmResourceId.ResourceId))
					if err != nil {
						return make([]ReplicationRegion, 0), err
					}
					region.SystemCreatedAcrId = containerRegistryId.ID()
				}
			}
		}

		result = append(result, region)
	}

	return result, nil
}

type PublicNetworkAccessState string

const (
	PublicNetworkAccessStateEnabled  PublicNetworkAccessState = "Enabled"
	PublicNetworkAccessStateDisabled PublicNetworkAccessState = "Disabled"
)

const systemCreatedStorageAccountTypeDefault = "Standard_LRS"
