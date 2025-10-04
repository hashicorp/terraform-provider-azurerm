// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource                   = DesktopVirtualizationWorkspaceApplicationGroupAssociationResource{}
	_ sdk.ResourceWithStateMigration = DesktopVirtualizationWorkspaceApplicationGroupAssociationResource{}
)

type DesktopVirtualizationWorkspaceApplicationGroupAssociationResource struct{}

func (DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) ModelObject() interface{} {
	return &DesktopVirtualizationWorkspaceApplicationGroupAssociationResourceModel{}
}

func (DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := parse.WorkspaceApplicationGroupAssociationID(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}

func (DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) ResourceType() string {
	return "azurerm_virtual_desktop_workspace_application_group_association"
}

func (DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceApplicationGroupAssociationV0ToV1{},
		},
	}
}

type DesktopVirtualizationWorkspaceApplicationGroupAssociationResourceModel struct {
	WorkspaceId        string `tfschema:"workspace_id"`
	ApplicationGroupId string `tfschema:"application_group_id"`
}

func (r DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspace.ValidateWorkspaceID,
		},

		"application_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationgroup.ValidateApplicationGroupID,
		},
	}
}

func (r DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient

			var model DesktopVirtualizationWorkspaceApplicationGroupAssociationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace <-> Application Group Association creation.")
			workspaceId, err := workspace.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return err
			}
			applicationGroupId, err := applicationgroup.ParseApplicationGroupID(model.ApplicationGroupId)
			if err != nil {
				return err
			}
			associationId := parse.NewWorkspaceApplicationGroupAssociationId(*workspaceId, *applicationGroupId)

			locks.ByName(workspaceId.WorkspaceName, workspaceResourceType)
			defer locks.UnlockByName(workspaceId.WorkspaceName, workspaceResourceType)

			locks.ByName(applicationGroupId.ApplicationGroupName, DesktopVirtualizationApplicationGroupResource{}.ResourceType())
			defer locks.UnlockByName(applicationGroupId.ApplicationGroupName, DesktopVirtualizationApplicationGroupResource{}.ResourceType())

			existing, err := client.Get(ctx, *workspaceId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", *workspaceId)
				}

				return fmt.Errorf("retrieving %s: %+v", *workspaceId, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *workspaceId)
			}
			workspaceModel := *existing.Model

			applicationGroupAssociations := []string{}
			if props := workspaceModel.Properties; props != nil && props.ApplicationGroupReferences != nil {
				applicationGroupAssociations = *props.ApplicationGroupReferences
			}

			applicationGroupIdStr := applicationGroupId.ID()
			if associationExists(workspaceModel.Properties, applicationGroupIdStr) {
				return metadata.ResourceRequiresImport(r.ResourceType(), associationId)
			}
			applicationGroupAssociations = append(applicationGroupAssociations, applicationGroupIdStr)

			payload := workspace.WorkspacePatch{
				Properties: &workspace.WorkspacePatchProperties{
					ApplicationGroupReferences: &applicationGroupAssociations,
				},
				Tags: workspaceModel.Tags,
			}
			if _, err = client.Update(ctx, *workspaceId, payload); err != nil {
				return fmt.Errorf("creating association between %s and %s: %+v", *workspaceId, *applicationGroupId, err)
			}

			metadata.SetID(associationId)
			return nil
		},
	}
}

func (r DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient

			state := DesktopVirtualizationWorkspaceApplicationGroupAssociationResourceModel{}

			id, err := parse.WorkspaceApplicationGroupAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspace, err := client.Get(ctx, id.Workspace)
			if err != nil {
				if response.WasNotFound(workspace.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", id.Workspace)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id.Workspace, err)
			}
			if model := workspace.Model; model != nil {
				applicationGroupId := id.ApplicationGroup.ID()
				exists := associationExists(model.Properties, applicationGroupId)
				if !exists {
					log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.Workspace, id.ApplicationGroup)
					return metadata.MarkAsGone(id)
				}

				state.WorkspaceId = id.Workspace.ID()
				state.ApplicationGroupId = applicationGroupId
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationWorkspaceApplicationGroupAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.WorkspacesClient

			id, err := parse.WorkspaceApplicationGroupAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.Workspace.WorkspaceName, workspaceResourceType)
			defer locks.UnlockByName(id.Workspace.WorkspaceName, workspaceResourceType)

			locks.ByName(id.ApplicationGroup.ApplicationGroupName, DesktopVirtualizationApplicationGroupResource{}.ResourceType())
			defer locks.UnlockByName(id.ApplicationGroup.ApplicationGroupName, DesktopVirtualizationApplicationGroupResource{}.ResourceType())

			existing, err := client.Get(ctx, id.Workspace)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id.Workspace)
				}

				return fmt.Errorf("retrieving %s: %+v", id.Workspace, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id.Workspace)
			}
			model := *existing.Model

			applicationGroupReferences := []string{}
			applicationGroupId := id.ApplicationGroup.ID()
			if props := model.Properties; props != nil && props.ApplicationGroupReferences != nil {
				for _, referenceId := range *props.ApplicationGroupReferences {
					if strings.EqualFold(referenceId, applicationGroupId) {
						continue
					}

					applicationGroupReferences = append(applicationGroupReferences, referenceId)
				}
			}

			payload := workspace.WorkspacePatch{
				Properties: &workspace.WorkspacePatchProperties{
					ApplicationGroupReferences: &applicationGroupReferences,
				},
				Tags: model.Tags,
			}
			if _, err = client.Update(ctx, id.Workspace, payload); err != nil {
				return fmt.Errorf("removing association between %s and %s: %+v", id.Workspace, id.ApplicationGroup, err)
			}

			return nil
		},
	}
}

func associationExists(props *workspace.WorkspaceProperties, applicationGroupId string) bool {
	if props == nil || props.ApplicationGroupReferences == nil {
		return false
	}

	for _, id := range *props.ApplicationGroupReferences {
		if strings.EqualFold(id, applicationGroupId) {
			return true
		}
	}

	return false
}
