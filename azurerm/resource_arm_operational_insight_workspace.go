package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/arm/operationalinsights"
	"github.com/hashicorp/terraform/helper/schema"
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
			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azureRMNormalizeLocation,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"workspace_id": { // a.k.a. customer_id
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

	skuName := d.Get("sku")
	sku, err := getSku(skuName)
	if err != nil {
		return err
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

	cancel := make(chan struct{})
	workspaceChannel, error := client.CreateOrUpdate(resGroup, name, parameters, cancel)
	workspace := <-workspaceChannel
	err = <-error
	if err != nil {
		return err
	}
	// The cosmos DB read rest api again for getting id. Try remove it.
	// read, err := client.Get(resGroup, name)
	// if err != nil {
	//	return err
	//}

	//if read.ID == nil {
	//	return fmt.Errorf("Cannot read Operational Inight Workspace '%s' (resource group %s) ID", name, resGroup)
	//}
	d.SetId(*workspace.ID)

	return resourceArmOperationalInsightWorkspaceRead(d, meta)

}

func resourceArmOperationalInsightWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	// I don't understand why we can get the meta data. How the framework set it.
	client := meta.(*ArmClient).workspacesClient
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workspaces"]

	resp, err := client.Get(resGroup, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Operational Insight workspaces '%s': %s", name, err)
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
		log.Printf("[ERROR] Unable to List Shared keys for Operatinal Insight workspaces %s: %s", name, err)
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

func getSku(skuName interface{}) (*operationalinsights.Sku, error) {
	if skuName == nil {
		return nil, nil
	}
	skuEnum, err := getSkuNameEnum(skuName.(string))
	if err != nil {
		return nil, err
	}
	return &operationalinsights.Sku{
		Name: skuEnum,
	}, nil
}

func getSkuNameEnum(skuName string) (operationalinsights.SkuNameEnum, error) {
	switch skuName {
	case "Free":
		return operationalinsights.Free, nil
	case "PerNode":
		return operationalinsights.PerNode, nil
	case "Premium":
		return operationalinsights.Premium, nil
	case "Standalone":
		return operationalinsights.Standalone, nil
	case "Standard":
		return operationalinsights.Standard, nil
	case "Unlimited":
		return operationalinsights.Unlimited, nil
	default:
		return operationalinsights.Free, fmt.Errorf("Sku name not found")
	}
}

func validateAzureRmOperationalInsightWorkspaceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	r, _ := regexp.Compile("^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$")
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("Workspace Name can only contain alphabet, number, and '-' charactor. You can not use '-' as the start and end of the name."))
	}
	return
}
