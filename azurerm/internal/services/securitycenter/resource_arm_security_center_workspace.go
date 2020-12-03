package securitycenter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// only valid name is default
// Message="Invalid workspace settings name 'kttest' , only default is allowed "
const securityCenterWorkspaceName = "default"

func resourceArmSecurityCenterWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSecurityCenterWorkspaceCreateUpdate,
		Read:   resourceArmSecurityCenterWorkspaceRead,
		Update: resourceArmSecurityCenterWorkspaceCreateUpdate,
		Delete: resourceArmSecurityCenterWorkspaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmSecurityCenterWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("%v", err)
	}

	contact := security.WorkspaceSetting{
		WorkspaceSettingProperties: &security.WorkspaceSettingProperties{
			Scope:       utils.String(d.Get("scope").(string)),
			WorkspaceID: utils.String(workspaceID.ID("")),
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
	stateConf := &resource.StateChangeConf{
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
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Waiting: %+v", err)
	}

	if d.IsNewResource() {
		d.SetId(*resp.(security.WorkspaceSetting).ID)
	}

	return resourceArmSecurityCenterWorkspaceRead(d, meta)
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

			if resourcePrice.PricingProperties.PricingTier == security.Standard {
				return true, nil
			}
		}
	}

	return false, nil
}

func resourceArmSecurityCenterWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
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
		d.Set("workspace_id", properties.WorkspaceID)
	}

	return nil
}

func resourceArmSecurityCenterWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
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
