package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Default Authorization Rule/Policy created by Azure, used to populate the
// default connection strings and keys
var eventHubNamespaceDefaultAuthorizationRule = "RootManageSharedAccessKey"

func resourceArmEventHubNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubNamespaceCreateUpdate,
		Read:   resourceArmEventHubNamespaceRead,
		Update: resourceArmEventHubNamespaceCreateUpdate,
		Delete: resourceArmEventHubNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(eventhub.Basic),
					string(eventhub.Standard),
				}, true),
			},

			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"auto_inflate_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"kafka_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"maximum_throughput_units": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 20),
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
	}
}

func resourceArmEventHubNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventhub.NamespacesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Namespace %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	capacity := int32(d.Get("capacity").(int))
	tags := d.Get("tags").(map[string]interface{})
	autoInflateEnabled := d.Get("auto_inflate_enabled").(bool)
	kafkaEnabled := d.Get("kafka_enabled").(bool)

	parameters := eventhub.EHNamespace{
		Location: &location,
		Sku: &eventhub.Sku{
			Name:     eventhub.SkuName(sku),
			Tier:     eventhub.SkuTier(sku),
			Capacity: &capacity,
		},
		EHNamespaceProperties: &eventhub.EHNamespaceProperties{
			IsAutoInflateEnabled: utils.Bool(autoInflateEnabled),
			KafkaEnabled:         utils.Bool(kafkaEnabled),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("maximum_throughput_units"); ok {
		parameters.EHNamespaceProperties.MaximumThroughputUnits = utils.Int32(int32(v.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating eventhub namespace: %+v", err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub Namespace %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventHubNamespaceRead(d, meta)
}

func resourceArmEventHubNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventhub.NamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on EventHub Namespace %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("sku", string(resp.Sku.Name))
	d.Set("capacity", resp.Sku.Capacity)

	keys, err := client.ListKeys(ctx, resGroup, name, eventHubNamespaceDefaultAuthorizationRule)
	if err != nil {
		log.Printf("[WARN] Unable to List default keys for EventHub Namespace %q: %+v", name, err)
	} else {
		d.Set("default_primary_connection_string", keys.PrimaryConnectionString)
		d.Set("default_secondary_connection_string", keys.SecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	if props := resp.EHNamespaceProperties; props != nil {
		d.Set("auto_inflate_enabled", props.IsAutoInflateEnabled)
		d.Set("kafka_enabled", props.KafkaEnabled)
		d.Set("maximum_throughput_units", int(*props.MaximumThroughputUnits))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmEventHubNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventhub.NamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request of EventHub Namespace %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return waitForEventHubNamespaceToBeDeleted(ctx, client, resGroup, name)
}

func waitForEventHubNamespaceToBeDeleted(ctx context.Context, client *eventhub.NamespacesClient, resourceGroup, name string) error {
	// we can't use the Waiter here since the API returns a 200 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for EventHub Namespace (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: eventHubNamespaceStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: 40 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for EventHub NameSpace (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func eventHubNamespaceStateStatusCodeRefreshFunc(ctx context.Context, client *eventhub.NamespacesClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving EventHub Namespace %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the EventHub Namespace %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
