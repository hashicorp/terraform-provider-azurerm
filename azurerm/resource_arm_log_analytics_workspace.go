package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsWorkspaceCreateUpdate,
		Read:   resourceArmLogAnalyticsWorkspaceRead,
		Update: resourceArmLogAnalyticsWorkspaceCreateUpdate,
		Delete: resourceArmLogAnalyticsWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRmLogAnalyticsWorkspaceName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(operationalinsights.Free),
					string(operationalinsights.PerNode),
					string(operationalinsights.Premium),
					string(operationalinsights.Standalone),
					string(operationalinsights.Standard),
					string(operationalinsights.Unlimited),
					string(operationalinsights.PerGB2018),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"retention_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(30, 730),
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmLogAnalyticsWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspacesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics workspace creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)

	skuName := d.Get("sku").(string)
	sku := &operationalinsights.Sku{
		Name: operationalinsights.SkuNameEnum(skuName),
	}

	retentionInDays := int32(d.Get("retention_in_days").(int))

	tags := d.Get("tags").(map[string]interface{})

	parameters := operationalinsights.Workspace{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		WorkspaceProperties: &operationalinsights.WorkspaceProperties{
			Sku:             sku,
			RetentionInDays: &retentionInDays,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Log Analytics Workspace '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLogAnalyticsWorkspaceRead(d, meta)

}

func resourceArmLogAnalyticsWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspacesClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics workspaces '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("workspace_id", resp.CustomerID)
	d.Set("portal_url", resp.PortalURL)
	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}
	d.Set("retention_in_days", resp.RetentionInDays)

	sharedKeys, err := client.GetSharedKeys(ctx, resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}

func resourceArmLogAnalyticsWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspacesClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Delete(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Log Analytics Workspaces '%s': %+v", name, err)
	}

	return nil
}

func validateAzureRmLogAnalyticsWorkspaceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	r, _ := regexp.Compile("^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$")
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("Workspace Name can only contain alphabet, number, and '-' character. You can not use '-' as the start and end of the name"))
	}

	length := len(value)
	if length > 63 || 4 > length {
		errors = append(errors, fmt.Errorf("Workspace Name can only be between 4 and 63 letters"))
	}

	return
}
