package azurerm

import (
	"fmt"
	"log"
	"strings"

	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/eventhub"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"

func resourceArmEventHubNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubNamespaceCreate,
		Read:   resourceArmEventHubNamespaceRead,
		Update: resourceArmEventHubNamespaceCreate,
		Delete: resourceArmEventHubNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateEventHubNamespaceSku,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateEventHubNamespaceCapacity,
			},

			"auto_inflate_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"maximum_throughput_units": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"default_primary_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmEventHubNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	namespaceClient := client.eventHubNamespacesClient
	log.Printf("[INFO] preparing arguments for Azure ARM EventHub Namespace creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	tags := d.Get("tags").(map[string]interface{})

	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)

	parameters := eventhub.EHNamespace{
		Location: &location,
		Sku: &eventhub.Sku{
			Name:     eventhub.SkuName(sku),
			Tier:     eventhub.SkuTier(sku),
			Capacity: &capacity,
		},
		EHNamespaceProperties: &eventhub.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		maximumThroughputUnits := v.(int)
		parameters.EHNamespaceProperties.MaximumThroughputUnits = utils.Int32(int32(maximumThroughputUnits))
	}

	_, error := namespaceClient.CreateOrUpdate(resGroup, name, parameters, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := namespaceClient.Get(resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub Namespace %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventHubNamespaceRead(d, meta)
}

func resourceArmEventHubNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	namespaceClient := meta.(*ArmClient).eventHubNamespacesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	resp, err := namespaceClient.Get(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure EventHub Namespace %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("resource_group_name", resGroup)
	d.Set("sku", string(resp.Sku.Name))
	d.Set("capacity", resp.Sku.Capacity)

	keys, err := namespaceClient.ListKeys(resGroup, name, eventHubNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[ERROR] Unable to List default keys for Namespace %s: %+v", name, err)
	} else {
		d.Set("default_primary_connection_string", keys.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.SecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	if props := resp.EHNamespaceProperties; props != nil {
		d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
		d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmEventHubNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	namespaceClient := meta.(*ArmClient).eventHubNamespacesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	deleteResp, error := namespaceClient.Delete(resGroup, name, make(chan struct{}))
	resp := <-deleteResp
	err = <-error

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error issuing Azure ARM delete request of EventHub Namespace '%s': %+v", name, err)
	}

	return nil
}

func validateEventHubNamespaceSku(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	skus := map[string]bool{
		"basic":    true,
		"standard": true,
	}

	if !skus[value] {
		errors = append(errors, fmt.Errorf("EventHub Namespace SKU can only be Basic or Standard"))
	}
	return
}

func validateEventHubNamespaceCapacity(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	maxCapacity := 20

	if value > maxCapacity || value < 1 {
		errors = append(errors, fmt.Errorf("EventHub Namespace Capacity must be 20 or fewer Throughput Units for Basic or Standard SKU"))
	}
	return
}
