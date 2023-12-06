// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopWorkspaceApplicationGroupAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopWorkspaceApplicationGroupAssociationCreate,
		Read:   resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead,
		Delete: resourceVirtualDesktopWorkspaceApplicationGroupAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceApplicationGroupAssociationID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WorkspaceApplicationGroupAssociationV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
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
		},
	}
}

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace <-> Application Group Association creation.")
	workspaceId, err := workspace.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	applicationGroupId, err := applicationgroup.ParseApplicationGroupID(d.Get("application_group_id").(string))
	if err != nil {
		return err
	}
	associationId := parse.NewWorkspaceApplicationGroupAssociationId(*workspaceId, *applicationGroupId).ID()

	locks.ByName(workspaceId.WorkspaceName, workspaceResourceType)
	defer locks.UnlockByName(workspaceId.WorkspaceName, workspaceResourceType)

	locks.ByName(applicationGroupId.ApplicationGroupName, applicationGroupType)
	defer locks.UnlockByName(applicationGroupId.ApplicationGroupName, applicationGroupType)

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
	model := *existing.Model

	applicationGroupAssociations := []string{}
	if props := model.Properties; props != nil && props.ApplicationGroupReferences != nil {
		applicationGroupAssociations = *props.ApplicationGroupReferences
	}

	applicationGroupIdStr := applicationGroupId.ID()
	if associationExists(model.Properties, applicationGroupIdStr) {
		return tf.ImportAsExistsError("azurerm_virtual_desktop_workspace_application_group_association", associationId)
	}
	applicationGroupAssociations = append(applicationGroupAssociations, applicationGroupIdStr)

	payload := workspace.WorkspacePatch{
		Properties: &workspace.WorkspacePatchProperties{
			ApplicationGroupReferences: &applicationGroupAssociations,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, *workspaceId, payload); err != nil {
		return fmt.Errorf("creating association between %s and %s: %+v", *workspaceId, *applicationGroupId, err)
	}

	d.SetId(associationId)
	return resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead(d, meta)
}

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceApplicationGroupAssociationID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := client.Get(ctx, id.Workspace)
	if err != nil {
		if response.WasNotFound(workspace.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.Workspace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.Workspace, err)
	}
	if model := workspace.Model; model != nil {
		applicationGroupId := id.ApplicationGroup.ID()
		exists := associationExists(model.Properties, applicationGroupId)
		if !exists {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.Workspace, id.ApplicationGroup)
			d.SetId("")
			return nil
		}

		d.Set("workspace_id", id.Workspace.ID())
		d.Set("application_group_id", applicationGroupId)
	}

	return nil
}

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceApplicationGroupAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Workspace.WorkspaceName, workspaceResourceType)
	defer locks.UnlockByName(id.Workspace.WorkspaceName, workspaceResourceType)

	locks.ByName(id.ApplicationGroup.ApplicationGroupName, applicationGroupType)
	defer locks.UnlockByName(id.ApplicationGroup.ApplicationGroupName, applicationGroupType)

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
