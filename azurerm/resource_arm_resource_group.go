package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmResourceGroupCreateUpdate,
		Read:   resourceArmResourceGroupRead,
		Update: resourceArmResourceGroupCreateUpdate,
		Exists: resourceArmResourceGroupExists,
		Delete: resourceArmResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmResourceGroupName,
			},

			"location": locationSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmResourceGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resourceGroupClient

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})
	parameters := resources.Group{
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}
	_, err := client.CreateOrUpdate(name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating resource group: %+v", err)
	}

	resp, err := client.Get(name)
	if err != nil {
		return fmt.Errorf("Error retrieving resource group: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmResourceGroupRead(d, meta)
}

func resourceArmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.ResourceGroup

	resp, err := client.Get(name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading resource group: %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmResourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.ResourceGroup

	resp, err := client.Get(name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			return false, nil
		}

		return false, fmt.Errorf("Error reading resource group: %+v", err)
	}

	return true, nil
}

func resourceArmResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.ResourceGroup

	deleteResp, deleteErr := client.Delete(name, make(chan struct{}))
	resp := <-deleteResp
	err = <-deleteErr
	if err != nil {
		if responseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting resource group: %+v", err)
	}

	return nil
}

func validateArmResourceGroupName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if len(value) > 80 {
		es = append(es, fmt.Errorf("%q may not exceed 80 characters in length", k))
	}

	if strings.HasSuffix(value, ".") {
		es = append(es, fmt.Errorf("%q may not end with a period", k))
	}

	// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
	if matched := regexp.MustCompile(`^[-\w\._\(\)]+$`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters, dash, underscores, parentheses and periods", k))
	}

	return
}
