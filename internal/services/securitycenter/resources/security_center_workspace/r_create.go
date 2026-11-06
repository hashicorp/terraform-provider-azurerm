// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_workspace

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	// TODO: split this create/update

	client := meta.(*clients.Client).SecurityCenter.WorkspaceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewWorkspaceID(subscriptionId, securityCenterWorkspaceName)
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.WorkspaceSettingName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_security_center_workspace", id.ID())
			}
		}
	}

	logAnalyticsWorkspaceId, err := workspaces.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}

	contact := security.WorkspaceSetting{
		WorkspaceSettingProperties: &security.WorkspaceSettingProperties{
			Scope:       pointer.To(d.Get("scope").(string)),
			WorkspaceID: pointer.To(logAnalyticsWorkspaceId.ID()),
		},
	}

	if d.IsNewResource() {
		if _, err = client.Create(ctx, id.WorkspaceSettingName, contact); err != nil {
			return fmt.Errorf("creating Security Center Workspace: %+v", err)
		}
	} else if _, err = client.Update(ctx, id.WorkspaceSettingName, contact); err != nil {
		return fmt.Errorf("updating Security Center Workspace: %+v", err)
	}

	// api returns "" for workspace id after an create/update and eventually the new value
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Waiting"},
		Target:     []string{"Populated"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (any, string, error) {
			resp, err2 := client.Get(ctx, id.WorkspaceSettingName)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("reading Security Center Workspace: %+v", err2)
			}

			if properties := resp.WorkspaceSettingProperties; properties != nil {
				if properties.WorkspaceID != nil && *properties.WorkspaceID != "" {
					return resp, "Populated", nil
				}
			}

			return resp, "Waiting", nil
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSecurityCenterWorkspaceRead(d, meta)
}
