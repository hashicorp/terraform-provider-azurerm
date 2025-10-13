// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
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
	StorageAccountType            string `tfschema:"storage_account_type"`
	HnsEnabled                    bool   `tfschema:"hns_enabled"`
	SystemCreatedStorageAccountId string `tfschema:"system_created_storage_account_id"`
	SystemCreatedAcrId            string `tfschema:"system_created_container_registry_id"`
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
	return registrymanagement.ValidateRegistryID
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
			Default:  true,
		},

		"main_region": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: regionSchema(),
			},
		},

		"replication_region": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: regionSchema(),
			},
		},

		"tags": commonschema.Tags(),
	}
}

func regionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationWithoutForceNew(),

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Standard_LRS",
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

		"hns_enabled": {
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
				Name:     pointer.To(model.Name),
				Location: model.Location,
				Tags:     pointer.To(model.Tags),
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}
			param.Identity = expandedIdentity

			var prop = registrymanagement.Registry{
				PublicNetworkAccess: pointer.To("Disabled"),
			}
			if model.PublicNetworkAccessEnabled {
				prop.PublicNetworkAccess = pointer.To("Enabled")
			}

			var regions = []registrymanagement.RegistryRegionArmDetails{expandRegistryRegionDetail(model.MainRegion[0])}

			for _, region := range model.ReplicationRegion {
				regions = append(regions, expandRegistryRegionDetail(region))
			}
			prop.RegionDetails = &regions
			param.Properties = prop

			response, err := client.RegistriesCreateOrUpdate(ctx, id, param)
			if response.HttpResponse != nil {
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

			id, err := registrymanagement.ParseRegistryID(metadata.ResourceData.Id())
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

			if existing.Model == nil {
				return fmt.Errorf("retrieving existing model for %s", id)
			}

			param := existing.Model

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
				param.Properties.PublicNetworkAccess = pointer.To("Disabled")
				if model.PublicNetworkAccessEnabled {
					param.Properties.PublicNetworkAccess = pointer.To("Enabled")
				}
			}

			if metadata.ResourceData.HasChanges("main_region", "replication_region") {
				var regions []registrymanagement.RegistryRegionArmDetails

				if len(model.MainRegion) > 0 {
					regions = append(regions, expandRegistryRegionDetail(model.MainRegion[0]))
				}

				for _, region := range model.ReplicationRegion {
					regions = append(regions, expandRegistryRegionDetail(region))
				}

				param.Properties.RegionDetails = &regions
			}

			response, err := client.RegistriesCreateOrUpdate(ctx, id, *param)
			if response.HttpResponse != nil {
				pollerType, err := custompollers.NewMachineLearningRegistryPoller(client, id, response.HttpResponse)
				if err != nil {
					return fmt.Errorf("creating poller: %+v", err)
				}
				if pollerType == nil {
					return fmt.Errorf("no poller created for updating %s", id)
				}
				poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("polling update of %s: %+v", id, err)
				}
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
	return registrymanagement.RegistryRegionArmDetails{
		Location: pointer.To(input.Location),
		AcrDetails: &[]registrymanagement.AcrDetails{
			{
				SystemCreatedAcrAccount: &registrymanagement.SystemCreatedAcrAccount{
					AcrAccountSku: pointer.To(string(registrymanagement.SkuTierPremium)),
				},
			},
		},
		StorageAccountDetails: &[]registrymanagement.StorageAccountDetails{
			{
				SystemCreatedStorageAccount: &registrymanagement.SystemCreatedStorageAccount{
					StorageAccountHnsEnabled: pointer.To(input.HnsEnabled),
					StorageAccountType:       pointer.To(input.StorageAccountType),
				},
			},
		},
	}
}

func flattenRegistryRegionDetails(input *[]registrymanagement.RegistryRegionArmDetails, location string) []ReplicationRegion {
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
			if systemAccount := sa[0].SystemCreatedStorageAccount; systemAccount != nil {
				region.StorageAccountType = pointer.From(systemAccount.StorageAccountType)
				region.HnsEnabled = pointer.From(systemAccount.StorageAccountHnsEnabled)
				region.SystemCreatedStorageAccountId = pointer.From(pointer.From(systemAccount.ArmResourceId).ResourceId)
			}
		}

		if acr := pointer.From(item.AcrDetails); len(acr) > 0 {
			if systemAcr := acr[0].SystemCreatedAcrAccount; systemAcr != nil {
				region.SystemCreatedAcrId = pointer.From(pointer.From(systemAcr.ArmResourceId).ResourceId)
			}
		}

		result = append(result, region)
	}
	return result
}
