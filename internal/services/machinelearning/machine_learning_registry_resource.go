// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/registries"
	registry "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineLearningRegistry struct{}

type ReplicationRegion struct {
	Location                      string `tfschema:"location"`
	CustomStorageAccountId        string `tfschema:"custom_storage_account_id"`
	CustomAcrAccountId            string `tfschema:"custom_container_registry_account_id"`
	StorageAccountType            string `tfschema:"storage_account_type"`
	HsnEnabled                    bool   `tfschema:"hsn_enabled"`
	SystemCreatedStorageAccountId string `tfschema:"system_created_storage_account_id"`
	SystemCreatedAcrId            string `tfschema:"system_created_container_registry_id"`
	// AcrSku                       string `tfschema:"container_registry_sku"` // Only allowed value is "Premium"
	// PublicAccessEnabled          bool   `tfschema:"public_access_enabled"` Not returned by API
}

type MachineLearningRegistryModel struct {
	Name                          string                                     `tfschema:"name"`
	ResourceGroupName             string                                     `tfschema:"resource_group_name"`
	PublicNetworkAccessEnabled    bool                                       `tfschema:"public_network_access_enabled"`
	MainRegion                    []ReplicationRegion                        `tfschema:"main_region"`
	ReplicationRegion             []ReplicationRegion                        `tfschema:"replication_region"`
	Location                      string                                     `tfschema:"location"`
	Identity                      []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	DiscoveryUrl                  string                                     `tfschema:"discovery_url"`
	IntellectualPropertyPublisher string                                     `tfschema:"intellectual_property_publisher"`
	MlFlowRegistryUri             string                                     `tfschema:"ml_flow_registry_uri"`
	ManagedResourceGroup          string                                     `tfschema:"managed_resource_group"`
	Tags                          map[string]string                          `tfschema:"tags"`
}

func (r MachineLearningRegistry) ModelObject() interface{} {
	return &MachineLearningRegistryModel{}
}

func (r MachineLearningRegistry) ResourceType() string {
	return "azurerm_machine_learning_registry"
}

func (r MachineLearningRegistry) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return registry.ValidateRegistryID
}

var _ sdk.ResourceWithUpdate = MachineLearningRegistry{}

func (r MachineLearningRegistry) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{1,31}$`),
				"Machine Learning Registry name must be 2 - 32 characters long, begin and end with alphanumeric and may contain only alphanumeric, hyphen or underscore.",
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"main_region": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// Azure requires main region location to be the same as resource location, no need to specify main region location
					"location": commonschema.LocationOptional(),

					"custom_storage_account_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ConflictsWith: []string{"main_region.0.storage_account_type"},
						ValidateFunc:  commonids.ValidateStorageAccountID,
					},

					"custom_container_registry_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: registries.ValidateRegistryID,
					},

					"storage_account_type": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Default:       "Standard_LRS",
						ConflictsWith: []string{"main_region.0.custom_storage_account_id"},
						AtLeastOneOf:  []string{"main_region.0.custom_storage_account_id", "main_region.0.storage_account_type"},
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

					"hsn_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"system_created_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"system_created_container_registry_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"replication_region": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.Location(),

					"custom_storage_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"custom_container_registry_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: registries.ValidateRegistryID,
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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

					"hsn_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"system_created_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"system_created_container_registry_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r MachineLearningRegistry) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"intellectual_property_publisher": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ml_flow_registry_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// API will return a generated mrg even provided in request
		"managed_resource_group": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

			id := registry.NewRegistryID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.RegistriesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_registry", id.ID())
			}

			param := registry.RegistryTrackedResource{
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}
			param.Identity = expandedIdentity

			var prop = registry.Registry{
				PublicNetworkAccess: pointer.To("Disabled"),
			}
			if model.PublicNetworkAccessEnabled {
				prop.PublicNetworkAccess = pointer.To("Enabled")
			}
			var regions []registry.RegistryRegionArmDetails

			regions = []registry.RegistryRegionArmDetails{
				{
					Location: pointer.To(model.Location),
					AcrDetails: &[]registry.AcrDetails{
						{
							SystemCreatedAcrAccount: &registry.SystemCreatedAcrAccount{},
						},
					},
					StorageAccountDetails: &[]registry.StorageAccountDetails{
						{
							SystemCreatedStorageAccount: &registry.SystemCreatedStorageAccount{},
						},
					},
				},
			}
			if len(model.MainRegion) > 0 {
				regions[0] = expandRegistryRegionDetail(model.MainRegion[0])
			}

			for _, region := range model.ReplicationRegion {
				regions = append(regions, expandRegistryRegionDetail(region))
			}
			prop.RegionDetails = &regions
			param.Properties = prop

			response, err := client.RegistriesCreateOrUpdate(ctx, id, param)

			if err != nil {
				pollerType, err := custompollers.NewMachineLearningRegistryPoller(client, id, response.HttpResponse)
				if err != nil {
					return fmt.Errorf("creating poller: %+v", err)
				}
				if pollerType == nil {
					return fmt.Errorf("no poller created for creating %s", id)
				}
				poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("polling creation of %s: %+v", id, err)
				}
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

			id, err := registry.ParseRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.RegistriesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("reading nil model %s", *id)
			}

			identityIds, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(resp.Model.Identity)
			if err != nil {
				return fmt.Errorf("flatten identity %s: %+v", *id, err)
			}
			prop := resp.Model.Properties
			model := MachineLearningRegistryModel{
				Name:                          id.RegistryName,
				ResourceGroupName:             id.ResourceGroupName,
				Identity:                      identityIds,
				Location:                      resp.Model.Location,
				PublicNetworkAccessEnabled:    pointer.From(prop.PublicNetworkAccess) == "Enabled",
				Tags:                          pointer.From(resp.Model.Tags),
				MlFlowRegistryUri:             pointer.From(prop.MlFlowRegistryUri),
				DiscoveryUrl:                  pointer.From(prop.DiscoveryURL),
				IntellectualPropertyPublisher: pointer.From(prop.IntellectualPropertyPublisher),
				ManagedResourceGroup:          pointer.From(pointer.From(prop.ManagedResourceGroup).ResourceId),
			}

			regions := flattenRegistryRegionDetails(prop.RegionDetails, resp.Model.Location)
			for i, region := range regions {
				if i == 0 {
					model.MainRegion = []ReplicationRegion{region}
				} else {
					model.ReplicationRegion = append(model.ReplicationRegion, region)
				}
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

			id, err := registry.ParseRegistryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MachineLearningRegistryModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decode state: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				var update registry.PartialRegistryPartialTrackedResource
				update.Identity, err = identity.ExpandLegacySystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding identity: %+v", err)
				}

				_, err = client.RegistriesUpdate(ctx, *id, update)
				if err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}

			if o, n := metadata.ResourceData.GetChange("replication_region"); !cmp.Equal(o, n) {
				// remove first: if exists in o but not in n
				var toRemove []string
				oldLocations := map[string]bool{}
				for _, item := range o.([]interface{}) {
					loc := item.(map[string]interface{})["location"].(string)
					oldLocations[loc] = true
				}
				for _, item := range n.([]interface{}) {
					loc := item.(map[string]interface{})["location"].(string)
					if _, ok := oldLocations[loc]; !ok {
						toRemove = append(toRemove, loc)
					}
				}
				if len(toRemove) > 0 {
					var regions []registry.RegistryRegionArmDetails
					for _, item := range toRemove {
						regions = append(regions, registry.RegistryRegionArmDetails{
							Location: pointer.To(item),
						})
					}
					var req = registry.RegistryTrackedResource{
						Properties: registry.Registry{
							RegionDetails: &regions,
						},
					}
					if err = client.RegistriesRemoveRegionsThenPoll(ctx, *id, req); err != nil {
						return fmt.Errorf("remove regions %s: %+v", *id, err)
					}
				}
				// add regions and remove regions separately
				req := registry.RegistryTrackedResource{
					Location: state.Location,
				}
				var regions []registry.RegistryRegionArmDetails
				mainCopy := state.MainRegion[0]
				mainCopy.Location = state.Location
				regions = append(regions, expandRegistryRegionDetail(mainCopy))
				for _, region := range state.ReplicationRegion {
					regions = append(regions, expandRegistryRegionDetail(region))
				}
				req.Properties = registry.Registry{
					RegionDetails: &regions,
				}
				response, err := client.RegistriesCreateOrUpdate(ctx, *id, req)

				if err != nil {
					pollerType, err := custompollers.NewMachineLearningRegistryPoller(client, *id, response.HttpResponse)
					if err != nil {
						return fmt.Errorf("creating poller: %+v", err)
					}
					if pollerType == nil {
						return fmt.Errorf("no poller created for updating %s", *id)
					}
					poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
					if err := poller.PollUntilDone(ctx); err != nil {
						return fmt.Errorf("polling update of %s: %+v", *id, err)
					}
				}
			}

			return nil
		},
	}
}

func (r MachineLearningRegistry) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.RegistryManagement

			id, err := registry.ParseRegistryID(metadata.ResourceData.Id())
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

func expandRegistryRegionDetail(input ReplicationRegion) registry.RegistryRegionArmDetails {
	var result registry.RegistryRegionArmDetails
	if result.Location == nil {
		result.Location = pointer.To(input.Location)
	}
	var sa registry.StorageAccountDetails
	if input.CustomStorageAccountId != "" {
		// sa.UserCreatedStorageAccount = &registry.UserCreatedStorageAccount{
		// 	ArmResourceId: &registry.ArmResourceId{
		// 		ResourceId: pointer.To(input.CustomAcrAccountId),
		// 	},
		// }
	} else {
		sa.SystemCreatedStorageAccount = &registry.SystemCreatedStorageAccount{
			StorageAccountHnsEnabled: pointer.To(input.HsnEnabled),
			StorageAccountType:       pointer.To(input.StorageAccountType),
		}
	}
	result.StorageAccountDetails = &[]registry.StorageAccountDetails{sa}

	var acr registry.AcrDetails
	if input.CustomAcrAccountId != "" {
		// acr.UserCreatedAcrAccount = &registry.UserCreatedAcrAccount{
		// 	ArmResourceId: &registry.ArmResourceId{
		// 		ResourceId: pointer.To(input.CustomAcrAccountId),
		// 	},
		// }
	} else {
		acr.SystemCreatedAcrAccount = &registry.SystemCreatedAcrAccount{
			AcrAccountSku: pointer.To(string(registry.SkuTierPremium)),
		}
	}

	result.AcrDetails = &[]registry.AcrDetails{acr}
	return result
}

func flattenRegistryRegionDetails(input *[]registry.RegistryRegionArmDetails, location string) []ReplicationRegion {
	var result []ReplicationRegion
	if input == nil || len(*input) == 0 {
		return result
	}

	for i, item := range *input {
		var region ReplicationRegion
		region.Location = pointer.From(item.Location)
		if i == 0 {
			region.Location = location
		}
		if sa := pointer.From(item.StorageAccountDetails); len(sa) > 0 {
			// if customAccount := sa[0].UserCreatedStorageAccount; customAccount != nil {
			// 	region.CustomStorageAccountId = pointer.From(pointer.From(customAccount.ArmResourceId).ResourceId)
			// } else if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
			if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
				region.StorageAccountType = pointer.From(systemAccount.StorageAccountType)
				region.HsnEnabled = pointer.From(systemAccount.StorageAccountHnsEnabled)
				region.SystemCreatedStorageAccountId = pointer.From(pointer.From(systemAccount.ArmResourceId).ResourceId)
			}
		}

		if acr := pointer.From(item.AcrDetails); len(acr) > 0 {
			// if customAcr := acr[0].UserCreatedAcrAccount; customAcr != nil {
			// 	region.CustomAcrAccountId = pointer.From(pointer.From(customAcr.ArmResourceId).ResourceId)
			if systemAcr := acr[0].SystemCreatedAcrAccount; systemAcr != nil {
				region.SystemCreatedAcrId = pointer.From(pointer.From(systemAcr.ArmResourceId).ResourceId)
			}
		}

		result = append(result, region)
	}
	return result
}
