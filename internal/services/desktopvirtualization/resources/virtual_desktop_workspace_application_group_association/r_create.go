// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_workspace_application_group_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_workspace_application_group_association", associationId)
		}
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
