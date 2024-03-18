// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.DataSource = StackHCIClusterDataSource{}

type StackHCIClusterDataSource struct{}

func (r StackHCIClusterDataSource) ResourceType() string {
	return "azurerm_stack_hci_cluster"
}

func (r StackHCIClusterDataSource) ModelObject() interface{} {
	return &StackHCIClusterDataSourceModel{}
}

type StackHCIClusterDataSourceModel struct {
	Name                      string                         `tfschema:"name"`
	ResourceGroupName         string                         `tfschema:"resource_group_name"`
	Location                  string                         `tfschema:"location"`
	ClientId                  string                         `tfschema:"client_id"`
	TenantId                  string                         `tfschema:"tenant_id"`
	AutomanageConfigurationId string                         `tfschema:"automanage_configuration_id"`
	CloudId                   string                         `tfschema:"cloud_id"`
	ServiceEndpoint           string                         `tfschema:"service_endpoint"`
	ResourceProviderObjectId  string                         `tfschema:"resource_provider_object_id"`
	Identity                  []identity.ModelSystemAssigned `tfschema:"identity"`
	Tags                      map[string]interface{}         `tfschema:"tags"`
}

func (r StackHCIClusterDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r StackHCIClusterDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location": commonschema.LocationComputed(),

		"client_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"automanage_configuration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cloud_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_provider_object_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedIdentityComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (r StackHCIClusterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Clusters
			subscriptionId := metadata.Client.Account.SubscriptionId
			hciAssignmentClient := metadata.Client.Automanage.HCIAssignmentClient

			var cluster StackHCIClusterDataSourceModel
			if err := metadata.Decode(&cluster); err != nil {
				return err
			}

			id := clusters.NewClusterID(subscriptionId, cluster.ResourceGroupName, cluster.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				cluster.Location = location.Normalize(model.Location)
				cluster.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					cluster.ClientId = pointer.From(props.AadClientId)
					cluster.CloudId = pointer.From(props.CloudId)
					cluster.Identity = flattenSystemAssignedToModel(model.Identity)
					cluster.ResourceProviderObjectId = pointer.From(props.ResourceProviderObjectId)
					cluster.ServiceEndpoint = pointer.From(props.ServiceEndpoint)
					cluster.TenantId = pointer.From(props.AadTenantId)

					assignmentResp, err := hciAssignmentClient.Get(ctx, id.ResourceGroupName, id.ClusterName, "default")
					if err != nil && !utils.ResponseWasNotFound(assignmentResp.Response) {
						return err
					}
					configId := ""
					if !utils.ResponseWasNotFound(assignmentResp.Response) && assignmentResp.Properties != nil && assignmentResp.Properties.ConfigurationProfile != nil {
						automanageConfigId, err := configurationprofiles.ParseConfigurationProfileIDInsensitively(*assignmentResp.Properties.ConfigurationProfile)
						if err != nil {
							return err
						}
						configId = automanageConfigId.ID()
					}
					cluster.AutomanageConfigurationId = configId
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&cluster)
		},
	}
}

func flattenSystemAssignedToModel(input *identity.SystemAndUserAssignedMap) []identity.ModelSystemAssigned {
	if input == nil {
		return []identity.ModelSystemAssigned{}
	}

	if input.Type == identity.TypeNone {
		return []identity.ModelSystemAssigned{}
	}

	return []identity.ModelSystemAssigned{
		{
			Type:        input.Type,
			PrincipalId: input.PrincipalId,
			TenantId:    input.TenantId,
		},
	}
}
