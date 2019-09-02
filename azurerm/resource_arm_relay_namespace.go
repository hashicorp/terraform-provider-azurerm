package azurerm

import (
	"context"
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"time"

	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRelayNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRelayNamespaceCreateUpdate,
		Read:   resourceArmRelayNamespaceRead,
		Update: resourceArmRelayNamespaceCreateUpdate,
		Delete: resourceArmRelayNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(6, 50),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Deprecated:    "This property has been deprecated in favour of the 'sku_name' property and will be removed in version 2.0 of the provider",
				ConflictsWith: []string{"sku_name"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(relay.Standard),
							}, true),
						},
					},
				},
			},

			"sku_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"sku"},
				ValidateFunc: validation.StringInSlice([]string{
					string(relay.Standard),
				}, false),
			},

			"metric_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmRelayNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).relay.NamespacesClient
	ctx := meta.(*ArmClient).StopContext

	// Remove in 2.0
	var sku relay.Sku

	if inputs := d.Get("sku").([]interface{}); len(inputs) != 0 {
		input := inputs[0].(map[string]interface{})
		v := input["name"].(string)

		sku = relay.Sku{
			Name: utils.String(v),
			Tier: relay.SkuTier(v),
		}
	} else {
		// Keep in 2.0
		sku = relay.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
			Tier: relay.SkuTier(d.Get("sku_name").(string)),
		}
	}

	if *sku.Name == "" {
		return fmt.Errorf("either 'sku_name' or 'sku' must be defined in the configuration file")
	}

	log.Printf("[INFO] preparing arguments for Relay Namespace creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Relay Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_relay_namespace", *existing.ID)
		}
	}

	parameters := relay.Namespace{
		Location:            utils.String(location),
		Sku:                 &sku,
		NamespaceProperties: &relay.NamespaceProperties{},
		Tags:                expandedTags,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Relay Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on future for Relay Namespace %q (Resource Group %q) creation: %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Relay Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Relay Namespace %q (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmRelayNamespaceRead(d, meta)
}

func resourceArmRelayNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).relay.NamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

		return fmt.Errorf("Error making Read request on Relay Namespace %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		// Remove in 2.0
		if err := d.Set("sku", flattenRelayNamespaceSku(sku)); err != nil {
			return fmt.Errorf("Error setting 'sku': %+v", err)
		}

		if err := d.Set("sku_name", sku.Name); err != nil {
			return fmt.Errorf("Error setting 'sku_name': %+v", err)
		}
	} else {
		return fmt.Errorf("Error making Read request on Relay Namespace %q (Resource Group %q): Unable to retrieve 'sku' value", name, resourceGroup)
	}

	if props := resp.NamespaceProperties; props != nil {
		d.Set("metric_id", props.MetricID)
	}

	keysResp, err := client.ListKeys(ctx, resourceGroup, name, "RootManageSharedAccessKey")
	if err != nil {
		return fmt.Errorf("Error making ListKeys request on Relay Namespace %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmRelayNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).relay.NamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	// we can't make use of the Future here due to a bug where 404 isn't tracked as Successful
	log.Printf("[DEBUG] Waiting for Relay Namespace %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Deleted"},
		Refresh:    relayNamespaceDeleteRefreshFunc(ctx, client, resourceGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Relay Namespace %q (Resource Group %q) to be deleted: %s", name, resourceGroup, err)
	}

	return nil
}

func relayNamespaceDeleteRefreshFunc(ctx context.Context, client *relay.NamespacesClient, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}

			return nil, "Error", fmt.Errorf("Error issuing read request in relayNamespaceDeleteRefreshFunc to Relay Namespace %q (Resource Group %q): %s", name, resourceGroupName, err)
		}

		return res, "Pending", nil
	}
}

// Remove in 2.0
func flattenRelayNamespaceSku(input *relay.Sku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})
	if name := input.Name; name != nil {
		output["name"] = *name
	}
	return []interface{}{output}
}
