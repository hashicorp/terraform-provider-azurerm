package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

//only valid name is default
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

		Schema: map[string]*schema.Schema{
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
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
	priceClient := meta.(*ArmClient).securityCenter.PricingClient
	client := meta.(*ArmClient).securityCenter.WorkspaceClient
	ctx := meta.(*ArmClient).StopContext

	name := securityCenterWorkspaceName

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Security Center Workspace: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_workspace", *existing.ID)
		}
	}

	//get pricing tier, workspace can only be configured when tier is not Free.
	//API does not error, it just doesn't set the workspace scope
	price, err := priceClient.GetSubscriptionPricing(ctx, securityCenterSubscriptionPricingName)
	if err != nil {
		return fmt.Errorf("Error reading Security Center Subscription pricing: %+v", err)
	}

	if price.PricingProperties == nil {
		return fmt.Errorf("Security Center Subscription pricing propertier is nil")
	}
	if price.PricingProperties.PricingTier == security.Free {
		return fmt.Errorf("Security Center Subscription workspace cannot be set when pricing tier is `Free`")
	}

	contact := security.WorkspaceSetting{
		WorkspaceSettingProperties: &security.WorkspaceSettingProperties{
			Scope:       utils.String(d.Get("scope").(string)),
			WorkspaceID: utils.String(d.Get("workspace_id").(string)),
		},
	}

	if d.IsNewResource() {
		if _, err = client.Create(ctx, name, contact); err != nil {
			return fmt.Errorf("Error creating Security Center Workspace: %+v", err)
		}
	} else {
		if _, err = client.Update(ctx, name, contact); err != nil {
			return fmt.Errorf("Error updating Security Center Workspace: %+v", err)
		}
	}

	//api returns "" for workspace id after an create/update and eventually the new value
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Waiting"},
		Target:     []string{"Populated"},
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {

			resp, err2 := client.Get(ctx, name)
			if err2 != nil {
				return resp, "Error", fmt.Errorf("Error reading Security Center Workspace: %+v", err2)
			}

			if properties := resp.WorkspaceSettingProperties; properties != nil {
				if properties.WorkspaceID != nil && *properties.WorkspaceID != "" {
					return resp, "Populated", nil
				}
			}

			return resp, "Waiting", nil
		},
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting: %+v", err)
	}

	if d.IsNewResource() {
		d.SetId(*resp.(security.WorkspaceSetting).ID)
	}

	return resourceArmSecurityCenterWorkspaceRead(d, meta)
}

func resourceArmSecurityCenterWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenter.WorkspaceClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, securityCenterWorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Security Center Subscription Workspace was not found: %v", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Security Center Workspace: %+v", err)
	}

	if properties := resp.WorkspaceSettingProperties; properties != nil {
		d.Set("scope", properties.Scope)
		d.Set("workspace_id", properties.WorkspaceID)
	}

	return nil
}

func resourceArmSecurityCenterWorkspaceDelete(_ *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).securityCenter.WorkspaceClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, securityCenterWorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center Subscription Workspace was not found: %v", err)
			return nil
		}

		return fmt.Errorf("Error deleting Security Center Workspace: %+v", err)
	}

	return nil
}
