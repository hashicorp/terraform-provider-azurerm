package azurerm

import (
	"context"
	"fmt"
	"log"

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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(relay.Standard),
							}, true),
						},
					},
				},
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmRelayNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).relayNamespacesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Relay Namespace creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	sku := expandRelayNamespaceSku(d)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := relay.Namespace{
		Location:            utils.String(location),
		Sku:                 sku,
		NamespaceProperties: &relay.NamespaceProperties{},
		Tags:                expandedTags,
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
		return fmt.Errorf("Cannot read Relay Namespace %q (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmRelayNamespaceRead(d, meta)
}

func resourceArmRelayNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).relayNamespacesClient
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

		return fmt.Errorf("Error making Read request on Relay Namespace %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		flattenedSku := flattenRelayNamespaceSku(sku)
		if err := d.Set("sku", flattenedSku); err != nil {
			return fmt.Errorf("Error setting `sku`: %+v", err)
		}
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

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRelayNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).relayNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

func relayNamespaceDeleteRefreshFunc(ctx context.Context, client relay.NamespacesClient, resourceGroupName string, name string) resource.StateRefreshFunc {
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

func expandRelayNamespaceSku(d *schema.ResourceData) *relay.Sku {
	vs := d.Get("sku").([]interface{})
	v := vs[0].(map[string]interface{})

	name := v["name"].(string)

	return &relay.Sku{
		Name: utils.String(name),
		Tier: relay.SkuTier(name),
	}
}

func flattenRelayNamespaceSku(input *relay.Sku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{}, 0)
	if name := input.Name; name != nil {
		output["name"] = *name
	}
	return []interface{}{output}
}
