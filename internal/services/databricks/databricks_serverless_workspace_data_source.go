// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DatabricksServerlessWorkspaceDataSource struct{}

var _ sdk.DataSource = DatabricksServerlessWorkspaceDataSource{}

type DatabricksServerlessWorkspaceDataSourceModel struct {
	Name                       string                            `tfschema:"name"`
	ResourceGroupName          string                            `tfschema:"resource_group_name"`
	Location                   string                            `tfschema:"location"`
	EnhancedSecurityCompliance []EnhancedSecurityComplianceModel `tfschema:"enhanced_security_compliance"`
	WorkspaceId                string                            `tfschema:"workspace_id"`
	WorkspaceUrl               string                            `tfschema:"workspace_url"`
	Tags                       map[string]string                 `tfschema:"tags"`
}

func (DatabricksServerlessWorkspaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (DatabricksServerlessWorkspaceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"enhanced_security_compliance": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"automatic_cluster_update_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"compliance_security_profile_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"compliance_security_profile_standards": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"enhanced_security_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (DatabricksServerlessWorkspaceDataSource) ModelObject() interface{} {
	return nil
}

func (DatabricksServerlessWorkspaceDataSource) ResourceType() string {
	return "azurerm_databricks_serverless_workspace"
}

func (DatabricksServerlessWorkspaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			ctx, cancel := context.WithTimeout(ctx, metadata.ResourceData.Timeout(schema.TimeoutRead))
			defer cancel()

			client := metadata.Client.DataBricks.WorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state DatabricksServerlessWorkspaceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := workspaces.NewWorkspaceID(subscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.EnhancedSecurityCompliance = DatabricksServerlessWorkspaceResource{}.flattenDatabricksServerlessWorkspaceEnhancedSecurityComplianceDefinition(model.Properties.EnhancedSecurityCompliance)

				if model.Properties.WorkspaceId != nil {
					state.WorkspaceId = *model.Properties.WorkspaceId
				}

				if model.Properties.WorkspaceURL != nil {
					state.WorkspaceUrl = *model.Properties.WorkspaceURL
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
