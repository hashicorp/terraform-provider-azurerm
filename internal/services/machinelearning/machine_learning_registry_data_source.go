// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineLearningRegistryDataSource struct{}

type MachineLearningRegistryDataSourceModel struct {
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d MachineLearningRegistryDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"main_region": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hns_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
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
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"hns_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
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

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

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

		"tags": commonschema.TagsDataSource(),
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
				return fmt.Errorf("decoding Machine Learning Registry data source model: %+v", err)
			}

			id := registrymanagement.NewRegistryID(subscriptionId, model.ResourceGroupName, model.Name)

			resp, err := client.RegistriesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("reading nil model %s", id)
			}

			identityIds, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(resp.Model.Identity)
			if err != nil {
				return fmt.Errorf("flattening identity %s: %+v", id, err)
			}

			prop := resp.Model.Properties
			model = MachineLearningRegistryDataSourceModel{
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

			regions := flattenRegistryRegionDetails(prop.RegionDetails)
			for i, region := range regions {
				if i == 0 {
					model.MainRegion = []ReplicationRegion{region}
				} else {
					model.ReplicationRegion = append(model.ReplicationRegion, region)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
