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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2025-09-01/playwrightworkspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name playwright_workspace -service-package-name loadtestservice -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type PlaywrightWorkspaceResource struct{}

var _ sdk.ResourceWithIdentity = PlaywrightWorkspaceResource{}

type PlaywrightWorkspaceModel struct {
	Name                    string            `tfschema:"name"`
	ResourceGroupName       string            `tfschema:"resource_group_name"`
	Location                string            `tfschema:"location"`
	LocalAuthEnabled        bool              `tfschema:"local_auth_enabled"`
	RegionalAffinityEnabled bool              `tfschema:"regional_affinity_enabled"`
	DataplaneUri            string            `tfschema:"dataplane_uri"`
	WorkspaceId             string            `tfschema:"workspace_id"`
	Tags                    map[string]string `tfschema:"tags"`
}

func (PlaywrightWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`),
				"length of `name` must be between 3 and 24 characters (inclusive) and contain only numbers, letters, and hyphens (-)"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"regional_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": {
			Type:         schema.TypeMap,
			Optional:     true,
			ValidateFunc: tags.Validate,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (PlaywrightWorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dataplane_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (PlaywrightWorkspaceResource) ModelObject() interface{} {
	return &PlaywrightWorkspaceModel{}
}

func (PlaywrightWorkspaceResource) ResourceType() string {
	return "azurerm_playwright_workspace"
}

func (r PlaywrightWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.PlaywrightWorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config PlaywrightWorkspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := playwrightworkspaces.NewPlaywrightWorkspaceID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &playwrightworkspaces.PlaywrightWorkspaceProperties{}
			properties.LocalAuth = pointer.To(playwrightworkspaces.EnablementStatusDisabled)
			if config.LocalAuthEnabled {
				properties.LocalAuth = pointer.To(playwrightworkspaces.EnablementStatusEnabled)
			}

			properties.RegionalAffinity = pointer.To(playwrightworkspaces.EnablementStatusDisabled)
			if config.RegionalAffinityEnabled {
				properties.RegionalAffinity = pointer.To(playwrightworkspaces.EnablementStatusEnabled)
			}

			param := playwrightworkspaces.PlaywrightWorkspace{
				Location:   location.Normalize(config.Location),
				Properties: properties,
				Tags:       pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return r.Read().Func(ctx, metadata)
		},
	}
}

func (r PlaywrightWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.PlaywrightWorkspacesClient
			id, err := playwrightworkspaces.ParsePlaywrightWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			param := playwrightworkspaces.PlaywrightWorkspaceUpdate{
				Properties: &playwrightworkspaces.PlaywrightWorkspaceUpdateProperties{},
			}

			var config PlaywrightWorkspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("local_auth_enabled") {
				param.Properties.LocalAuth = pointer.To(playwrightworkspaces.EnablementStatusDisabled)
				if config.LocalAuthEnabled {
					param.Properties.LocalAuth = pointer.To(playwrightworkspaces.EnablementStatusEnabled)
				}
			}

			if metadata.ResourceData.HasChange("regional_affinity_enabled") {
				param.Properties.RegionalAffinity = pointer.To(playwrightworkspaces.EnablementStatusDisabled)
				if config.RegionalAffinityEnabled {
					param.Properties.RegionalAffinity = pointer.To(playwrightworkspaces.EnablementStatusEnabled)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(config.Tags)
			}

			if _, err := client.Update(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return r.Read().Func(ctx, metadata)
		},
	}
}

func (PlaywrightWorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.PlaywrightWorkspacesClient
			id, err := playwrightworkspaces.ParsePlaywrightWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := PlaywrightWorkspaceModel{
				Name:              id.PlaywrightWorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(&model.Location)

				if properties := model.Properties; properties != nil {
					if localAuth := properties.LocalAuth; localAuth != nil {
						state.LocalAuthEnabled = pointer.From(localAuth) == playwrightworkspaces.EnablementStatusEnabled
					}

					if regionalAffinity := properties.RegionalAffinity; regionalAffinity != nil {
						state.RegionalAffinityEnabled = pointer.From(regionalAffinity) == playwrightworkspaces.EnablementStatusEnabled
					}

					if dataplaneUri := properties.DataplaneUri; dataplaneUri != nil {
						state.DataplaneUri = pointer.From(dataplaneUri)
					}

					if workspaceId := properties.WorkspaceId; workspaceId != nil {
						state.WorkspaceId = pointer.From(workspaceId)
					}
				}

				state.Tags = pointer.From(model.Tags)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (PlaywrightWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.PlaywrightWorkspacesClient

			id, err := playwrightworkspaces.ParsePlaywrightWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (PlaywrightWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return playwrightworkspaces.ValidatePlaywrightWorkspaceID
}

func (PlaywrightWorkspaceResource) Identity() resourceids.ResourceId {
	return &playwrightworkspaces.PlaywrightWorkspaceId{}
}
