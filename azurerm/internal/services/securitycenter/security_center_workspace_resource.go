package securitycenter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// only valid name is default
// Message="Invalid workspace settings name 'kttest' , only default is allowed "
const securityCenterWorkspaceName = "default"

func resourceSecurityCenterWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterWorkspaceCreateUpdate,
		Read:   resourceSecurityCenterWorkspaceRead,
		Update: resourceSecurityCenterWorkspaceCreateUpdate,
		Delete: resourceSecurityCenterWorkspaceDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceSecurityCenterWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	priceClient := meta.(*clients.Client).SecurityCenter.PricingClient
	client := meta.(*clients.Client).SecurityCenter.WorkspaceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := securityCenterWorkspaceName

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Checking for presence of existing Security Center Workspace: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_workspace", *existing.ID)
		}
	}

	// get pricing tier, workspace can only be configured when tier is not Free.
	// API does not error, it just doesn't set the workspace scope
	isPricingStandard, err := isPricingStandard(ctx, priceClient)
	if err != nil {
		return fmt.Errorf("Checking Security Center Subscription pricing tier %v", err)
	}

	if !isPricingStandard {
		return fmt.Errorf("Security Center Subscription workspace cannot be set when pricing tier is `Free`")
	}

	workspaceID, err := parse.LogAnalyticsWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}

	contact := security.WorkspaceSetting{
		WorkspaceSettingProperties: &security.WorkspaceSettingProperties{
			Scope:       utils.String(d.Get("scope").(string)),
			WorkspaceID: utils.String(workspaceID.ID()),
		},
	}

	if d.IsNewResource() {
		if _, err = client.Create(ctx, name, contact); err != nil {
			return fmt.Errorf("Creating Security Center Workspace: %+v", err)
		}
	} else if _, err = client.Update(ctx, name, contact); err != nil {
		return fmt.Errorf("Updating Security Center Workspace: %+v", err)
	}

	// api returns "" for workspace id after an create/update and eventually the new value
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Waiting"},
		Target:     []string{"Populated"},
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {
			resp, err2 := client.Get(ctx, name)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("Reading Security Center Workspace: %+v", err2)
			}

			if properties := resp.WorkspaceSettingProperties; properties != nil {
				if properties.WorkspaceID != nil && *properties.WorkspaceID != "" {
					return resp, "Populated", nil
				}
			}

			return resp, "Waiting", nil
		},
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("Waiting: %+v", err)
	}

	if d.IsNewResource() {
		d.SetId(*resp.(security.WorkspaceSetting).ID)
	}

	return resourceSecurityCenterWorkspaceRead(d, meta)
}

func isPricingStandard(ctx context.Context, priceClient *security.PricingsClient) (bool, error) {
	prices, err := priceClient.List(ctx)
	if err != nil {
		return false, fmt.Errorf("Listing Security Center Subscription pricing: %+v", err)
	}

	if prices.Value != nil {
		for _, resourcePrice := range *prices.Value {
			if resourcePrice.PricingProperties == nil {
				return false, fmt.Errorf("%v Security Center Subscription pricing properties is nil", *resourcePrice.Type)
			}

			if resourcePrice.PricingProperties.PricingTier == security.PricingTierStandard {
				return true, nil
			}
		}
	}

	return false, nil
}

func resourceSecurityCenterWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.WorkspaceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, securityCenterWorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription Workspace was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Reading Security Center Workspace: %+v", err)
	}

	if properties := resp.WorkspaceSettingProperties; properties != nil {
		d.Set("scope", properties.Scope)
		workspaceId := ""
		if properties.WorkspaceID != nil {
			id, err := parse.LogAnalyticsWorkspaceID(*properties.WorkspaceID)
			if err != nil {
				return fmt.Errorf("Reading Security Center Log Analytics Workspace ID: %+v", err)
			}
			workspaceId = id.ID()
		}
		d.Set("workspace_id", utils.String(workspaceId))
	}

	return nil
}

func resourceSecurityCenterWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.WorkspaceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Delete(ctx, securityCenterWorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center Subscription Workspace was not found: %v", err)
			return nil
		}

		return fmt.Errorf("Deleting Security Center Workspace: %+v", err)
	}

	return nil
}
