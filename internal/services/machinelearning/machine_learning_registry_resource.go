// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/registries"
	registry "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// pollForRegistryCreation polls the GET endpoint until the registry is available
// This handles the known Azure ML Registry API bug where CreateOrUpdate returns 202 with no body
func pollForRegistryCreation(ctx context.Context, client *registry.RegistryManagementClient, id registry.RegistryId, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.RegistriesGet(ctx, id)
		if err == nil && resp.Model != nil {
			return nil
		}

		if err != nil && !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("unexpected error while polling for registry: %+v", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
		}
	}

	return fmt.Errorf("timed out waiting for registry to become available")
}

type MachineLearningRegistry struct{}

type ReplicationRegion struct {
	Location                     string `tfschema:"location"`
	CustomStorageAccountId       string `tfschema:"custom_storage_account_id"`
	CustomAcrAccountId           string `tfschema:"custom_acr_account_id"`
	StorageAccountType           string `tfschema:"storage_account_type"`
	PublicAccessEnabled          bool   `tfschema:"public_access_enabled"`
	HsnEnabled                   bool   `tfschema:"hsn_enabled"`
	SystemCreateStorageAccountId string `tfschema:"system_create_storage_account_id"`
	AcrSku                       string `tfschema:"acr_sku"`
	SystemCreatedAcrId           string `tfschema:"system_created_acr_id"`
}

type MachineLearningRegistryModel struct {
	Name                       string                                     `tfschema:"name"`
	ResourceGroupName          string                                     `tfschema:"resource_group_name"`
	PublicNetworkAccessEnabled bool                                       `tfschema:"public_network_access_enabled"`
	MainRegion                 []ReplicationRegion                        `tfschema:"main_region"`
	ReplicationRegions         []ReplicationRegion                        `tfschema:"replication_regions"`
	Location                   string                                     `tfschema:"location"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	// DiscoveryUrl                  string                                     `tfschema:"discovery_url"`
	IntellectualPropertyPublisher string            `tfschema:"intellectual_property_publisher"`
	MlFlowRegistryUri             string            `tfschema:"ml_flow_registry_uri"`
	ManagedResourceGroup          string            `tfschema:"managed_resource_group"`
	Tags                          map[string]string `tfschema:"tags"`
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

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"main_region": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.LocationComputed(),

					"custom_storage_account_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ConflictsWith: []string{"main_region.0.storage_account_type"},
						ValidateFunc:  commonids.ValidateStorageAccountID,
					},

					"custom_acr_account_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  registries.ValidateRegistryID,
						ConflictsWith: []string{"main_region.0.acr_sku"},
					},

					"storage_account_type": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
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

					"public_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"hsn_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"system_create_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"acr_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							// string(registries.SkuNameBasic),
							// string(registries.SkuNameStandard),
							string(registries.SkuNamePremium),
						}, false),
						ConflictsWith: []string{"main_region.0.custom_acr_account_id"},
					},

					"system_created_acr_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"replication_regions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": commonschema.Location(),

					"custom_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// ConflictsWith: []string{"replication_regions.0.storage_account_type"},
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"custom_acr_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: registries.ValidateRegistryID,
						// ConflictsWith: []string{"replication_regions.0.acr_sku"},
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// ConflictsWith: []string{"replication_regions.0.custom_storage_account_id"},
						// AtLeastOneOf:  []string{"replication_regions.0.custom_storage_account_id", "replication_regions.0.storage_account_type"},
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

					"public_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"hsn_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"system_create_storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"acr_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							// string(registries.SkuNameBasic),
							// string(registries.SkuNameStandard),
							string(registries.SkuNamePremium),
						}, false),
						// ConflictsWith: []string{"replication_regions.0.custom_acr_account_id"},
					},

					"system_created_acr_id": {
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
		// "discovery_url": {
		// 	Type:     pluginsdk.TypeString,
		// 	Computed: true,
		// },

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

			identityIns, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %v", err)
			}

			trackedResource := registry.RegistryTrackedResource{
				Name:     pointer.To(model.Name),
				Identity: identityIns,
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
			}

			var prop = registry.Registry{
				PublicNetworkAccess: pointer.To("Disabled"),
			}
			if model.PublicNetworkAccessEnabled {
				prop.PublicNetworkAccess = pointer.To("Enabled")
			}

			var regions []registry.RegistryRegionArmDetails
			mainCopy := model.MainRegion[0]
			mainCopy.Location = model.Location
			regions = append(regions, expandRegistryRegionDetail(mainCopy))
			for _, region := range model.ReplicationRegions {
				regions = append(regions, expandRegistryRegionDetail(region))
			}
			prop.RegionDetails = &regions
			trackedResource.Properties = prop

			// Try the normal operation first
			err = client.RegistriesCreateOrUpdateThenPoll(ctx, id, trackedResource)
			if err != nil {
				// Check if this is the known issue with 202 + no body
				if strings.Contains(err.Error(), "unexpected status 202") && strings.Contains(err.Error(), "received with no body") {
					// This is the known API bug - poll manually until the resource is available
					if pollErr := pollForRegistryCreation(ctx, client, id, 30*time.Minute); pollErr != nil {
						return fmt.Errorf("registry creation failed: initial 202 error followed by polling failure: %+v", pollErr)
					}
				} else {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

// only update tags and remove regions
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

			if o, n := metadata.ResourceData.GetChange("replication_regions"); !cmp.Equal(o, n) {
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
				for _, region := range state.ReplicationRegions {
					regions = append(regions, expandRegistryRegionDetail(region))
				}
				req.Properties = registry.Registry{
					RegionDetails: &regions,
				}
				// Try the normal operation first
				err = client.RegistriesCreateOrUpdateThenPoll(ctx, *id, req)
				if err != nil {
					// Check if this is the known issue with 202 + no body
					if strings.Contains(err.Error(), "unexpected status 202") && strings.Contains(err.Error(), "received with no body") {
						// This is the known API bug - poll manually until the resource is available
						if pollErr := pollForRegistryCreation(ctx, client, *id, 30*time.Minute); pollErr != nil {
							return fmt.Errorf("registry update failed: initial 202 error followed by polling failure: %+v", pollErr)
						}
					} else {
						return fmt.Errorf("updating %s: %+v", *id, err)
					}
				}

			}

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

			identityIns, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(resp.Model.Identity)
			if err != nil {
				return fmt.Errorf("flatten identity %s: %+v", *id, err)
			}
			prop := resp.Model.Properties
			model := MachineLearningRegistryModel{
				Name:                       id.RegistryName,
				ResourceGroupName:          id.ResourceGroupName,
				Identity:                   identityIns,
				Location:                   resp.Model.Location,
				PublicNetworkAccessEnabled: pointer.From(prop.PublicNetworkAccess) == ,
				Tags:                       pointer.From(resp.Model.Tags),
				MlFlowRegistryUri:          pointer.From(prop.MlFlowRegistryUri),
				// DiscoveryUrl:                  pointer.From(prop.DiscoveryUrl),
				IntellectualPropertyPublisher: pointer.From(prop.IntellectualPropertyPublisher),
				ManagedResourceGroup:          pointer.From(pointer.From(prop.ManagedResourceGroup).ResourceId),
			}

			regions := flattenRegistryRegionDetails(prop.RegionDetails)
			for _, region := range regions {
				if region.Location == resp.Model.Location {
					model.MainRegion = []ReplicationRegion{region}
				} else {
					model.ReplicationRegions = append(model.ReplicationRegions, region)
				}
			}

			return metadata.Encode(&model)
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
	result.Location = pointer.To(input.Location)
	var sa registry.StorageAccountDetails
	if input.CustomStorageAccountId != "" {
		// sa.UserCreatedStorageAccount = &registry.UserCreatedStorageAccount{
		// 	ArmResourceId: &registry.ArmResourceId{
		// 		ResourceId: pointer.To(input.CustomAcrAccountId),
		// 	},
		// }
	} else {
		sa.SystemCreatedStorageAccount = &registry.SystemCreatedStorageAccount{
			AllowBlobPublicAccess:    pointer.To(input.PublicAccessEnabled),
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
			AcrAccountSku: pointer.To(input.AcrSku),
		}
	}

	result.AcrDetails = &[]registry.AcrDetails{acr}
	return result
}

func flattenRegistryRegionDetails(input *[]registry.RegistryRegionArmDetails) []ReplicationRegion {
	var result []ReplicationRegion
	if input == nil || len(*input) == 0 {
		return result
	}

	for _, item := range *input {
		var region ReplicationRegion
		region.Location = pointer.From(item.Location)
		if sa := pointer.From(item.StorageAccountDetails); len(sa) > 0 {
			// if customAccount := sa[0].UserCreatedStorageAccount; customAccount != nil {
			// 	region.CustomStorageAccountId = pointer.From(pointer.From(customAccount.ArmResourceId).ResourceId)
			// } else if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
			if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
				region.StorageAccountType = pointer.From(systemAccount.StorageAccountType)
				// Don't overwrite PublicAccessEnabled if the API doesn't return AllowBlobPublicAccess
				// This allows the configured value to be preserved
				region.HsnEnabled = pointer.From(systemAccount.StorageAccountHnsEnabled)
				region.SystemCreateStorageAccountId = pointer.From(pointer.From(systemAccount.ArmResourceId).ResourceId)
			}
		}

		if acr := pointer.From(item.AcrDetails); len(acr) > 0 {
			// if customAcr := acr[0].UserCreatedAcrAccount; customAcr != nil {
			// 	region.CustomAcrAccountId = pointer.From(pointer.From(customAcr.ArmResourceId).ResourceId)
			if systemAcr := acr[0].SystemCreatedAcrAccount; systemAcr != nil {
				region.AcrSku = pointer.From(systemAcr.AcrAccountSku)
				region.SystemCreatedAcrId = pointer.From(pointer.From(systemAcr.ArmResourceId).ResourceId)
			}
		}

		result = append(result, region)
	}
	return result
}
