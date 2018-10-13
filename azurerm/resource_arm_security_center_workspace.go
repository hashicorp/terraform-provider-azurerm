package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/2017-08-01-preview/security"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
)

//only valid name is default
// Message="Invalid workspace settings name 'kttest' , only default is allowed "
const resourceArmSecurityCenterWorkspaceName = "default"

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
				ValidateFunc: validation.NoZeroValues,
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
	client := meta.(*ArmClient).securityCenterWorkspaceClient
	ctx := meta.(*ArmClient).StopContext

	name := resourceArmSecurityCenterWorkspaceName

	contact := security.WorkspaceSetting{
		WorkspaceSettingProperties: &security.WorkspaceSettingProperties{
			Scope:       utils.String(d.Get("scope").(string)),
			WorkspaceID: utils.String(d.Get("workspace_id").(string)),
		},
	}

	if d.IsNewResource() {
		_, err := client.Create(ctx, name, contact)
		if err != nil {
			return fmt.Errorf("Error creating Security Center Workspace: %+v", err)
		}
	} else {
		_, err := client.Update(ctx, name, contact)
		if err != nil {
			return fmt.Errorf("Error updating Security Center Workspace: %+v", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Waiting"},
		Target:     []string{"Populated"},
		Timeout:    60 * time.Minute,
		MinTimeout: 30 * time.Second,
		Refresh: func() (interface{}, string, error) {

			resp, err := client.Get(ctx, name)
			if err != nil {
				return resp, "Error", fmt.Errorf("Error reading Security Center Workspace: %+v", err)
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
	client := meta.(*ArmClient).securityCenterWorkspaceClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, resourceArmSecurityCenterWorkspaceName)
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
	client := meta.(*ArmClient).securityCenterWorkspaceClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, resourceArmSecurityCenterWorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Security Center Subscription Workspace was not found: %v", err)
			return nil
		}

		return fmt.Errorf("Error deleting Security Center Workspace: %+v", err)
	}

	return nil
}
