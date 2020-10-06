package desktopvirtualization

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualDesktopWorkspaceApplicationGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationCreate,
		Read:   resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationRead,
		Delete: resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DesktopVirtualizationWorkspaceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"application_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace <-> Application Group Association creation.")

	workspaceID := d.Get("workspace_id").(string)
	applicationGroupReferenceID := d.Get("application_group_id").(string)

	wsID, err := azure.ParseAzureResourceID(workspaceID)
	if err != nil {
		return err
	}

	agID, err := azure.ParseAzureResourceID(applicationGroupReferenceID)
	if err != nil {
		return err
	}

	workspaceName := wsID.Path["workspaces"]
	resourceGroup := wsID.ResourceGroup
	agName := agID.Path["applicationgroups"]

	resourceID := fmt.Sprintf("%s|%s", workspaceID, applicationGroupReferenceID)

	locks.ByName(workspaceName, "Microsoft.DesktopVirtualization/workspaces")
	defer locks.UnlockByName(workspaceName, "Microsoft.DesktopVirtualization/workspaces")

	locks.ByName(agName, "Microsoft.DesktopVirtualization/applicationgroups")
	defer locks.UnlockByName(agName, "Microsoft.DesktopVirtualization/applicationgroups")

	read, err := client.Get(ctx, resourceGroup, workspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Retrieving Virtual Desktop Workspace %q (Resource Group %q) was not found", workspaceName, resourceGroup)
		}

		return fmt.Errorf("Retrieving Virtual Desktop Workspace for Association %q (Resource Group %q): %+v", workspaceName, resourceGroup, err)
	}

	refs := read.ApplicationGroupReferences

	output := make([]string, 0)
	output = append(output, *refs...)
	if utils.SliceContainsValue(output, applicationGroupReferenceID) {
		return tf.ImportAsExistsError("azurerm_virtual_desktop_workspace_application_group_association", resourceID)
	}
	output = append(output, applicationGroupReferenceID)

	read.ApplicationGroupReferences = &output

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, workspaceName, read); err != nil {
		return fmt.Errorf("Updating Virtual Desktop Workspace Association for Application Group %q (Resource Group %q): %+v", workspaceName, resourceGroup, err)
	}

	d.SetId(resourceID)

	return resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationRead(d, meta)
}

func resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitID := strings.Split(d.Id(), "|")
	if len(splitID) != 2 {
		return fmt.Errorf("Expected ID to be in the format {workspaceID}/{networkSecurityGroup} but got %q", d.Id())
	}

	wsID, err := parse.DesktopVirtualizationWorkspaceID(splitID[0])
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, wsID.ResourceGroup, wsID.Name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Virtual Desktop Workspace %q was not found in Resource Group %q - removing from state!", wsID.Name, wsID.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Making Read request on Virtual Desktop Desktop Workspace %q (Resource Group %q): %+v", wsID.Name, wsID.ResourceGroup, err)
	}

	output := make([]string, 0)
	output = append(output, *read.ApplicationGroupReferences...)

	if !utils.SliceContainsValue(output, splitID[1]) {
		log.Printf("[DEBUG] Association between Virtual Desktop Workspace %q (Resource Group %q) and Virtual Desktop Application Group %q was not found - removing from state!", wsID.Name, wsID.ResourceGroup, splitID[1])
		d.SetId("")
		return nil
	}

	d.Set("application_group_id", splitID[1])

	return nil
}

func resourceArmVirtualDesktopWorkspaceApplicationGroupAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitID := strings.Split(d.Id(), "|")
	if len(splitID) != 2 {
		return fmt.Errorf("Expected ID to be in the format {workspaceID}/{networkSecurityGroup} but got %q", d.Id())
	}

	applicationGroupReferenceID := d.Get("application_group_id").(string)

	wsID, err := parse.DesktopVirtualizationWorkspaceID(splitID[0])
	if err != nil {
		return err
	}

	agID, err := parse.VirtualDesktopApplicationGroupID(splitID[1])
	if err != nil {
		return err
	}

	locks.ByName(wsID.Name, "Microsoft.DesktopVirtualization/workspaces")
	defer locks.UnlockByName(wsID.Name, "Microsoft.DesktopVirtualization/workspaces")

	locks.ByName(agID.Name, "Microsoft.DesktopVirtualization/applicationgroups")
	defer locks.UnlockByName(agID.Name, "Microsoft.DesktopVirtualization/applicationgroups")

	read, err := client.Get(ctx, wsID.ResourceGroup, wsID.Name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("Virtual Desktop Workspace %q (Resource Group %q) was not found", wsID.Name, wsID.ResourceGroup)
		}

		return fmt.Errorf("Retrieving Virtual Desktop Workspace %q (Resource Group %q): %+v", wsID.Name, wsID.ResourceGroup, err)
	}

	refs := read.ApplicationGroupReferences
	if refs == nil {
		return fmt.Errorf("ApplicationGroupReferences was nil for Virtual Desktop Workspace %q (Resource Group %q)", wsID.Name, wsID.ResourceGroup)
	}

	output := make([]string, 0)
	output = append(output, *refs...)
	output = utils.RemoveFromStringArray(output, applicationGroupReferenceID)

	read.ApplicationGroupReferences = &output

	if _, err = client.CreateOrUpdate(ctx, wsID.ResourceGroup, wsID.Name, read); err != nil {
		return fmt.Errorf("Updating Virtual Desktop Workspace Association for Application Group %q (Resource Group %q): %+v", wsID.Name, wsID.ResourceGroup, err)
	}

	return nil
}
