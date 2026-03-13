// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadtestservice

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PlaywrightWorkspaceDataSource struct{}

var _ sdk.DataSource = PlaywrightWorkspaceDataSource{}

type PlaywrightWorkspaceDataSourceModel struct {
	Name                    string            `tfschema:"name"`
	ResourceGroupName       string            `tfschema:"resource_group_name"`
	Location                string            `tfschema:"location"`
	LocalAuthEnabled        bool              `tfschema:"local_auth_enabled"`
	RegionalAffinityEnabled bool              `tfschema:"regional_affinity_enabled"`
	DataplaneUri            string            `tfschema:"dataplane_uri"`
	WorkspaceId             string            `tfschema:"workspace_id"`
	Tags                    map[string]string `tfschema:"tags"`
}

func (PlaywrightWorkspaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`[a-zA-Z0-9-]{3,24}`),
				"length of `name` must be between 3 and 24 characters (inclusive) and contain only numbers, letters, and hyphens (-)"),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (PlaywrightWorkspaceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"dataplane_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"regional_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (PlaywrightWorkspaceDataSource) ModelObject() interface{} {
	return nil
}

func (PlaywrightWorkspaceDataSource) ResourceType() string {
	return "azurerm_playwright_workspace"
}

func (PlaywrightWorkspaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.PlaywrightWorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state PlaywrightWorkspaceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := playwrightworkspaces.NewPlaywrightWorkspaceID(subscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(&model.Location)

				if properties := model.Properties; properties != nil {
					if dataplaneUri := properties.DataplaneUri; dataplaneUri != nil {
						state.DataplaneUri = pointer.From(dataplaneUri)
					}

					if localAuth := properties.LocalAuth; localAuth != nil {
						state.LocalAuthEnabled = pointer.From(localAuth) == playwrightworkspaces.EnablementStatusEnabled
					}

					if regionalAffinity := properties.RegionalAffinity; regionalAffinity != nil {
						state.RegionalAffinityEnabled = pointer.From(regionalAffinity) == playwrightworkspaces.EnablementStatusEnabled
					}

					if workspaceId := properties.WorkspaceId; workspaceId != nil {
						state.WorkspaceId = pointer.From(workspaceId)
					}
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
