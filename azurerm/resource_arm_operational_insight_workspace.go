package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/arm/operationalinsights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmOperationalInsightWorkspaceService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmOperationalInsightWorkspaceCreateUpdate,
		Read:   resourceArmOperationalInsightWorkspaceRead,
		Update: resourceArmOperationalInsightWorkspaceCreateUpdate,
		Delete: resourceArmOperationalInsightWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRmOperationalInsightWorkspaceName,
			},
			"location":            locationSchema(),
			"resource_group_name": resourceGroupNameSchema(),
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},
			"retention_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"primary_shared_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_shared_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceArmOperationalInsightWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspacesClient
	log.Printf("[INFO] preparing arguments for AzureRM Operational Insight workspace creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
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

	_, error := client.CreateOrUpdate(resGroup, name, parameters, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Operational Inight Workspace '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmOperationalInsightWorkspaceRead(d, meta)

}

func resourceArmOperationalInsightWorkspaceRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).workspacesClient
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Get(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Operational Insight workspaces '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", resp.Location)
	d.Set("resource_group_name", resGroup)
	d.Set("workspace_id", resp.CustomerID)
	d.Set("portal_url", resp.PortalURL)
	d.Set("sku", resp.Sku.Name)
	d.Set("retention_in_days", resp.RetentionInDays)

	sharedKeys, err := client.GetSharedKeys(resGroup, name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Operatinal Insight workspaces %s: %+v", name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}

func resourceArmOperationalInsightWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).workspacesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Delete(resGroup, name)

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Operational Insight Workspaces '%s': %+v", name, err)
	}

	return nil
}

func validateAzureRmOperationalInsightWorkspaceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	r, _ := regexp.Compile("^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$")
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("Workspace Name can only contain alphabet, number, and '-' charactor. You can not use '-' as the start and end of the name."))
	}
	return
}
