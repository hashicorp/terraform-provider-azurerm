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
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineLearningRegistryDataSource struct{}

type MachineLearningRegistryDataSourceModel struct {
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
	ReplicationRegion                                       []ReplicationRegion                        `tfschema:"replication_regions"`
	DiscoveryUrl                                            string                                     `tfschema:"discovery_url"`
	IntellectualPropertyPublisher                           string                                     `tfschema:"intellectual_property_publisher"`
	MlFlowRegistryUri                                       string                                     `tfschema:"machine_learning_flow_registry_uri"`
	ManagedResourceGroup                                    string                                     `tfschema:"managed_resource_group_id"`
	Tags                                                    map[string]string                          `tfschema:"tags"`
}

func (d MachineLearningRegistryDataSource) ModelObject() interface{} {
	return &MachineLearningRegistryDataSourceModel{}
}

func (d MachineLearningRegistryDataSource) ResourceType() string {
	return "azurerm_machine_learning_registry"
}

func (d MachineLearningRegistryDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return registrymanagement.ValidateRegistryID
}

func (d MachineLearningRegistryDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{2,32}$`),
				"Machine Learning Registry name must be between 3 and 33 characters long. Its first character has to be alphanumeric, and the rest may contain hyphens and underscores. No whitespace is allowed.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d MachineLearningRegistryDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := map[string]*pluginsdk.Schema{
		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

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

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"replication_regions": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: registryRegionDataSourceSchema(),
			},
		},

		"tags": commonschema.TagsDataSource(),
	}

	for k, v := range registryRegionDataSourceSchema() {
		attributes[k] = v
	}

	return attributes
}

func registryRegionDataSourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"system_created_container_registry_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_created_container_registry_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_created_container_registry_sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"system_created_storage_account_hierarchical_namespace_enabled": {
			Type:     pluginsdk.TypeBool,
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

		"system_created_storage_account_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d MachineLearningRegistryDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.RegistryManagement
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningRegistryDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := registrymanagement.NewRegistryID(subscriptionId, model.ResourceGroupName, model.Name)

			resp, err := client.RegistriesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			identityIds, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(resp.Model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity` %s: %+v", id, err)
			}

			prop := resp.Model.Properties
			model = MachineLearningRegistryDataSourceModel{
				Name:                          id.RegistryName,
				ResourceGroupName:             id.ResourceGroupName,
				Identity:                      identityIds,
				Location:                      resp.Model.Location,
				PublicNetworkAccessEnabled:    pointer.From(prop.PublicNetworkAccess) == string(PublicNetworkAccessStateEnabled),
				Tags:                          pointer.From(resp.Model.Tags),
				MlFlowRegistryUri:             pointer.From(prop.MlFlowRegistryUri),
				DiscoveryUrl:                  pointer.From(prop.DiscoveryURL),
				IntellectualPropertyPublisher: pointer.From(prop.IntellectualPropertyPublisher),
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
				return fmt.Errorf("flattening `region_details` %s: %+v", id, err)
			}

			for _, region := range regions {
				if location.Normalize(region.Location) == location.Normalize(resp.Model.Location) {
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

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
