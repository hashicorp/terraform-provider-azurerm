package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"regexp"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var serviceBusNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"

func resourceArmServiceBusNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceBusNamespaceCreate,
		Read:   resourceArmServiceBusNamespaceRead,
		Update: resourceArmServiceBusNamespaceCreate,
		Delete: resourceArmServiceBusNamespaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		MigrateState:  resourceAzureRMServiceBusNamespaceMigrateState,
		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{0,100}[a-zA-Z0-9]$"),
					"The namespace can contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number.",
				),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servicebus.Basic),
					string(servicebus.Standard),
					string(servicebus.Premium),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntInSlice([]int{1, 2, 4}),
			},

			"default_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tagsSchema(),
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {

			//If the SKU is not premium the API will always return 0 for capacity
			//so lets only allow it to be set if the SKU is premium
			if _, ok := d.GetOk("capacity"); ok {
				sku := d.Get("sku").(string)
				if strings.ToLower(sku) != strings.ToLower(string(servicebus.Premium)) {
					return fmt.Errorf("`capacity` can only be set for a Premium SKU")
				}
			}

			return nil
		},
	}
}

func resourceArmServiceBusNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM ServiceBus Namespace creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	sku := d.Get("sku").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := servicebus.SBNamespace{
		Location: &location,
		Sku: &servicebus.SBSku{
			Name: servicebus.SkuName(sku),
			Tier: servicebus.SkuTier(sku),
		},
		Tags: expandTags(tags),
	}

	if capacity, ok := d.GetOk("capacity"); ok {
		parameters.Sku.Capacity = utils.Int32(int32(capacity.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ServiceBus Namespace %q (resource group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceBusNamespaceRead(d, meta)
}

func resourceArmServiceBusNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure ServiceBus Namespace %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", strings.ToLower(string(sku.Name)))
		d.Set("capacity", sku.Capacity)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name, serviceBusNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	} else {
		d.Set("default_primary_connection_string", keys.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.SecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmServiceBusNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceBusNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	_, err = client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	return nil
}
