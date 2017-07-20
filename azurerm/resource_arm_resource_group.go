package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmResourceGroupCreate,
		Read:   resourceArmResourceGroupRead,
		Update: resourceArmResourceGroupUpdate,
		Exists: resourceArmResourceGroupExists,
		Delete: resourceArmResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
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

func validateArmResourceGroupName(v interface{}, k string) (ws []string, es []error) {
	value := v.(string)

	if len(value) > 80 {
		es = append(es, fmt.Errorf("%q may not exceed 80 characters in length", k))
	}

	if strings.HasSuffix(value, ".") {
		es = append(es, fmt.Errorf("%q may not end with a period", k))
	}

	if matched := regexp.MustCompile(`[\(\)\.a-zA-Z0-9_-]`).Match([]byte(value)); !matched {
		es = append(es, fmt.Errorf("%q may only contain alphanumeric characters, dash, underscores, parentheses and periods", k))
	}

	return
}

func resourceArmResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	resourceGroupClient := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Group Update, error parsing ID: %s", err)
	}

	resourceGroupName := id.ResourceGroup

	if !d.HasChange("tags") {
		return nil
	}

	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := resources.Group{
		Location: &location,
		Tags:     expandTags(tags),
	}

	result, err := resourceGroupClient.CreateOrUpdate(resourceGroupName, parameters)

	if err != nil {
		return fmt.Errorf("Error updating resource group: %s", err)
	}
	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error updating resource group: %v", result.Response.StatusCode)
	}

	return resourceArmResourceGroupRead(d, meta)
}

func resourceArmResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	resourceGroupClient := meta.(*ArmClient).resourceGroupClient

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := resources.Group{
		Location: &location,
		Tags:     expandTags(tags),
	}

	result, err := resourceGroupClient.CreateOrUpdate(name, parameters)

	if err != nil {
		return fmt.Errorf("Error creating resource group: %s", err)
	}
	if result.Response.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error creating resource group: %v", result.Response.StatusCode)
	}

	d.SetId(*result.ID)

	return resourceArmResourceGroupRead(d, meta)
}

func resourceArmResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	resourceGroupClient := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Group Read, error parsing ID: %s", err)
	}

	resourceGroupName := id.ResourceGroup

	result, err := resourceGroupClient.Get(resourceGroupName)

	if err != nil {
		return fmt.Errorf("Resource group %s, error reading: %s", resourceGroupName, err)
	}

	// covers the case where resource has been deleted outside TF but is still in state
	if result.Response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if result.Response.StatusCode != http.StatusOK {
		log.Printf("[INFO] Error reading resource group %q - removing from state", d.Id())
		d.SetId("")
		return fmt.Errorf("Error reading resource group: %v", result.Response.StatusCode)
	}

	d.Set("name", *result.Name)
	d.Set("location", azureRMNormalizeLocation(*result.Location))
	flattenAndSetTags(d, result.Tags)

	return nil
}

func resourceArmResourceGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	resourceGroupClient := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return false, fmt.Errorf("Resource Group Exists, error parsing ID: %s", err)
	}

	resourceGroupName := id.ResourceGroup

	result, err := resourceGroupClient.CheckExistence(resourceGroupName)

	if err != nil {
		return false, fmt.Errorf("Error checking existence of resource group: %s", err)
	}
	if result.Response.StatusCode == http.StatusNoContent {
		return true, nil
	}

	return false, nil
}

func resourceArmResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	resourceGroupClient := meta.(*ArmClient).resourceGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Group Delete, error parsing ID: %s", err)
	}

	resourceGroupName := id.ResourceGroup

	delResult, delErr := resourceGroupClient.Delete(resourceGroupName, make(chan struct{}))
	err = <-delErr
	result := <-delResult

	if err != nil {
		return fmt.Errorf("Error deleting resource group: %s", err)
	}
	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting resource group: %v", result.Response.StatusCode)
	}

	return nil

}
